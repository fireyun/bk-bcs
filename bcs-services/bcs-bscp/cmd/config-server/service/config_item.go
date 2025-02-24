/*
Tencent is pleased to support the open source community by making Basic Service Configuration Platform available.
Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at
http://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "as IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"

	"bscp.io/pkg/criteria/errf"
	"bscp.io/pkg/iam/meta"
	"bscp.io/pkg/kit"
	"bscp.io/pkg/logs"
	pbcs "bscp.io/pkg/protocol/config-server"
	pbcommit "bscp.io/pkg/protocol/core/commit"
	pbci "bscp.io/pkg/protocol/core/config-item"
	pbcontent "bscp.io/pkg/protocol/core/content"
	pbds "bscp.io/pkg/protocol/data-service"
	"bscp.io/pkg/types"
)

// CreateConfigItem create config item with option
func (s *Service) CreateConfigItem(ctx context.Context, req *pbcs.CreateConfigItemReq) (
	*pbcs.CreateConfigItemResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.CreateConfigItemResp)

	authRes := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.ConfigItem, Action: meta.Create,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, authRes)
	if err != nil {
		return resp, nil
	}
	// 1. create config_item
	cciReq := &pbds.CreateConfigItemReq{
		Attachment: &pbci.ConfigItemAttachment{
			BizId: req.BizId,
			AppId: req.AppId,
		},
		Spec: &pbci.ConfigItemSpec{
			Name:     req.Name,
			Path:     req.Path,
			FileType: req.FileType,
			FileMode: req.FileMode,
			Memo:     req.Memo,
			Permission: &pbci.FilePermission{
				User:      req.User,
				UserGroup: req.UserGroup,
				Privilege: req.Privilege,
			},
		},
	}
	cciResp, err := s.client.DS.CreateConfigItem(grpcKit.RpcCtx(), cciReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("create config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 2. create content
	ccReq := &pbds.CreateContentReq{
		Attachment: &pbcontent.ContentAttachment{
			ConfigItemId: cciResp.Id,
			BizId:        req.BizId,
			AppId:        req.AppId,
		},
		Spec: &pbcontent.ContentSpec{
			Signature: req.Sign,
			ByteSize:  req.ByteSize,
		},
	}
	ccResp, err := s.client.DS.CreateContent(grpcKit.RpcCtx(), ccReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("create config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 3. create commit
	ccmReq := &pbds.CreateCommitReq{
		Attachment: &pbcommit.CommitAttachment{
			BizId:        req.BizId,
			AppId:        req.AppId,
			ConfigItemId: cciResp.Id,
		},
		ContentId: ccResp.Id,
		Memo:      req.Memo,
	}
	_, err = s.client.DS.CreateCommit(grpcKit.RpcCtx(), ccmReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("create config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	resp.Data = &pbcs.CreateConfigItemResp_RespData{
		Id: cciResp.Id,
	}

	return resp, nil
}

// UpdateConfigItem update config item with option
func (s *Service) UpdateConfigItem(ctx context.Context, req *pbcs.UpdateConfigItemReq) (
	*pbcs.UpdateConfigItemResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.UpdateConfigItemResp)

	authRes := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.ConfigItem, Action: meta.Update,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, authRes)
	if err != nil {
		return resp, nil
	}

	// 1. update config_item
	r := &pbds.UpdateConfigItemReq{
		Id: req.Id,
		Attachment: &pbci.ConfigItemAttachment{
			BizId: req.BizId,
			AppId: req.AppId,
		},
		Spec: &pbci.ConfigItemSpec{
			Name:     req.Name,
			Path:     req.Path,
			FileType: req.FileType,
			FileMode: req.FileMode,
			Memo:     req.Memo,
			Permission: &pbci.FilePermission{
				User:      req.User,
				UserGroup: req.UserGroup,
				Privilege: req.Privilege,
			},
		},
	}
	_, err = s.client.DS.UpdateConfigItem(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("update config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 2. check if content sign changed,if changed,create content and commit
	// 2.1. get latest commit and compare content sign
	glcReq := &pbds.GetLatestCommitReq{
		BizId:        req.BizId,
		AppId:        req.AppId,
		ConfigItemId: req.Id,
	}
	glcResp, err := s.client.DS.GetLatestCommit(grpcKit.RpcCtx(), glcReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("get config item latest commit failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 2.2 if latest content sign equals request content sign,no need to commit
	if glcResp.Data.Spec.Content.Signature == req.Sign {
		resp.Code = errf.OK
		return resp, nil
	}

	// 2.3 if latest content sign not equals request content sign,create content and commit
	ccReq := &pbds.CreateContentReq{
		Attachment: &pbcontent.ContentAttachment{
			ConfigItemId: req.Id,
			BizId:        req.BizId,
			AppId:        req.AppId,
		},
		Spec: &pbcontent.ContentSpec{
			Signature: req.Sign,
			ByteSize:  req.ByteSize,
		},
	}
	ccResp, err := s.client.DS.CreateContent(grpcKit.RpcCtx(), ccReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("create config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 4. create commit
	ccmReq := &pbds.CreateCommitReq{
		Attachment: &pbcommit.CommitAttachment{
			BizId:        req.BizId,
			AppId:        req.AppId,
			ConfigItemId: req.Id,
		},
		ContentId: ccResp.Id,
		Memo:      req.Memo,
	}
	_, err = s.client.DS.CreateCommit(grpcKit.RpcCtx(), ccmReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("create config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}
	resp.Code = errf.OK
	return resp, nil
}

// DeleteConfigItem delete config item with option
func (s *Service) DeleteConfigItem(ctx context.Context, req *pbcs.DeleteConfigItemReq) (
	*pbcs.DeleteConfigItemResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.DeleteConfigItemResp)

	authRes := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.ConfigItem, Action: meta.Delete,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, authRes)
	if err != nil {
		return resp, nil
	}

	r := &pbds.DeleteConfigItemReq{
		Id: req.Id,
		Attachment: &pbci.ConfigItemAttachment{
			BizId: req.BizId,
			AppId: req.AppId,
		},
	}
	_, err = s.client.DS.DeleteConfigItem(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("delete config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	return resp, nil
}

// GetConfigItem get config item with content
func (s *Service) GetConfigItem(ctx context.Context, req *pbcs.GetConfigItemReq) (
	*pbcs.GetConfigItemResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.GetConfigItemResp)

	authRes := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.ConfigItem, Action: meta.Find}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, authRes)
	if err != nil {
		return resp, nil
	}

	// 1. get config item
	gciReq := &pbds.GetConfigItemReq{
		Id:    req.Id,
		BizId: req.BizId,
		AppId: req.AppId,
	}
	gciResp, err := s.client.DS.GetConfigItem(grpcKit.RpcCtx(), gciReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("get config item failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 2. get latest commit
	glcReq := &pbds.GetLatestCommitReq{
		BizId:        req.BizId,
		AppId:        req.AppId,
		ConfigItemId: req.Id,
	}
	glcResp, err := s.client.DS.GetLatestCommit(grpcKit.RpcCtx(), glcReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("get config item latest commit failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	// 3. get content
	gcReq := &pbds.GetContentReq{
		Id:    glcResp.Data.Spec.ContentId,
		BizId: req.BizId,
		AppId: req.AppId,
	}
	gcResp, err := s.client.DS.GetContent(grpcKit.RpcCtx(), gcReq)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("get config item content failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	resp.Data = &pbcs.GetConfigItemResp_RespData{
		ConfigItem: gciResp.Data,
		Content:    gcResp.Data,
	}
	return resp, nil
}

// ListConfigItems list config item with filter
func (s *Service) ListConfigItems(ctx context.Context, req *pbcs.ListConfigItemsReq) (
	*pbcs.ListConfigItemsResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.ListConfigItemsResp)

	authRes := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.ConfigItem, Action: meta.Find}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, authRes)
	if err != nil {
		return resp, nil
	}

	if req.Page == nil {
		errf.Error(errf.New(errf.InvalidParameter, "page is null")).AssignResp(grpcKit, resp)
		return resp, nil
	}

	// TODO: list latest release and compare each config item exists and latest commit id to get changing status

	if err := req.Page.BasePage().Validate(types.DefaultPageOption); err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		return resp, nil
	}

	r := &pbds.ListConfigItemsReq{
		BizId:  req.BizId,
		AppId:  req.AppId,
		Filter: req.Filter,
		Page:   req.Page,
	}
	rp, err := s.client.DS.ListConfigItems(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("list config items failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	resp.Data = &pbcs.ListConfigItemsResp_RespData{
		Count:   rp.Count,
		Details: rp.Details,
	}
	return resp, nil
}
