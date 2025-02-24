/*
Tencent is pleased to support the open source community by making Basic Service Configuration Platform available.
Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at
http://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
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
	pbds "bscp.io/pkg/protocol/data-service"
	"bscp.io/pkg/types"
)

// Publish publish a strategy
func (s *Service) Publish(ctx context.Context, req *pbcs.PublishReq) (
	*pbcs.PublishResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.PublishResp)

	res := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.Strategy, Action: meta.Publish,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, res)
	if err != nil {
		return resp, nil
	}

	r := &pbds.PublishReq{
		BizId:     req.BizId,
		AppId:     req.AppId,
		ReleaseId: req.ReleaseId,
		All:       req.All,
		Groups:    req.Groups,
	}
	rp, err := s.client.DS.Publish(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("publish strategy failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	resp.Data = &pbcs.PublishResp_RespData{
		Id: rp.PublishedStrategyHistoryId,
	}
	return resp, nil
}

// FinishPublish finish the published strategy
func (s *Service) FinishPublish(ctx context.Context, req *pbcs.FinishPublishReq) (
	*pbcs.FinishPublishResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.FinishPublishResp)

	res := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.Strategy, Action: meta.FinishPublish,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, res)
	if err != nil {
		return resp, nil
	}

	r := &pbds.FinishPublishReq{
		BizId:      req.BizId,
		AppId:      req.AppId,
		StrategyId: req.Id,
	}
	_, err = s.client.DS.FinishPublish(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("finish publish strategy failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	return resp, nil
}

// ListPublishedStrategyHistories list the published strategy histories.
func (s *Service) ListPublishedStrategyHistories(ctx context.Context, req *pbcs.ListPubStrategyHistoriesReq) (
	*pbcs.ListPubStrategyHistoriesResp, error) {

	grpcKit := kit.FromGrpcContext(ctx)
	resp := new(pbcs.ListPubStrategyHistoriesResp)

	res := &meta.ResourceAttribute{Basic: &meta.Basic{Type: meta.PSH, Action: meta.Find,
		ResourceID: req.AppId}, BizID: req.BizId}
	err := s.authorizer.AuthorizeWithResp(grpcKit, resp, res)
	if err != nil {
		return resp, nil
	}

	if req.Page == nil {
		errf.Error(errf.New(errf.InvalidParameter, "page is null")).AssignResp(grpcKit, resp)
		return resp, nil
	}

	if err = req.Page.BasePage().Validate(types.DefaultPageOption); err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		return resp, nil
	}

	r := &pbds.ListPubStrategyHistoriesReq{
		BizId:  req.BizId,
		AppId:  req.AppId,
		Filter: req.Filter,
		Page:   req.Page,
	}
	rp, err := s.client.DS.ListPublishedStrategyHistories(grpcKit.RpcCtx(), r)
	if err != nil {
		errf.Error(err).AssignResp(grpcKit, resp)
		logs.Errorf("list published strategy histories failed, err: %v, rid: %s", err, grpcKit.Rid)
		return resp, nil
	}

	resp.Code = errf.OK
	resp.Data = &pbcs.ListPubStrategyHistoriesResp_RespData{
		Count:   rp.Count,
		Details: rp.Details,
	}
	return resp, nil
}
