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

// Package service NOTES
package service

import (
	"errors"
	"fmt"
	"net/http"

	"bscp.io/pkg/cc"
	"bscp.io/pkg/criteria/errf"
	"bscp.io/pkg/iam/auth"
	"bscp.io/pkg/logs"
	"bscp.io/pkg/metrics"
	pbds "bscp.io/pkg/protocol/data-service"
	"bscp.io/pkg/serviced"
	esbcli "bscp.io/pkg/thirdparty/esb/client"
	"bscp.io/pkg/thirdparty/repo"
	"bscp.io/pkg/tools"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Service do all the data service's work
type Service struct {
	client  *ClientSet
	gateway *gateway
	// authorizer auth related operations.
	authorizer auth.Authorizer
}

// NewService create a service instance.
func NewService(sd serviced.Discover) (*Service, error) {
	client, err := newClientSet(sd, cc.ConfigServer().Network.TLS)
	if err != nil {
		return nil, fmt.Errorf("new client set failed, err: %v", err)
	}

	state, ok := sd.(serviced.State)
	if !ok {
		return nil, errors.New("discover convert state failed")
	}
	gateway, err := newGateway(state)
	if err != nil {
		return nil, fmt.Errorf("new gateway failed, err: %v", err)
	}

	authorizer, err := auth.NewAuthorizer(sd, cc.ConfigServer().Network.TLS)
	if err != nil {
		return nil, fmt.Errorf("new authorizer failed, err: %v", err)
	}

	return &Service{
		client:     client,
		gateway:    gateway,
		authorizer: authorizer,
	}, nil
}

// Handler return service's handler.
func (s *Service) Handler() (http.Handler, error) {
	if s.gateway == nil {
		return nil, errors.New("gateway is nil")
	}

	return s.gateway.handler(), nil
}

func newClientSet(sd serviced.Discover, tls cc.TLSConfig) (*ClientSet, error) {
	logs.Infof("start initialize the client set.")

	opts := make([]grpc.DialOption, 0)

	// add dial load balancer.
	opts = append(opts, sd.LBRoundRobin())

	if !tls.Enable() {
		// dial without ssl
		opts = append(opts, grpc.WithInsecure())
	} else {
		// dial with ssl.
		tlsC, err := tools.ClientTLSConfVerify(tls.InsecureSkipVerify, tls.CAFile, tls.CertFile, tls.KeyFile,
			tls.Password)
		if err != nil {
			return nil, fmt.Errorf("init client set tls config failed, err: %v", err)
		}

		cred := credentials.NewTLS(tlsC)
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	// connect data service.
	dsConn, err := grpc.Dial(serviced.GrpcServiceDiscoveryName(cc.DataServiceName), opts...)
	if err != nil {
		logs.Errorf("dial data service failed, err: %v", err)
		return nil, errf.New(errf.Unknown, fmt.Sprintf("dial data service failed, err: %v", err))
	}

	repoSetting := cc.ConfigServer().Repo
	repoCli, err := repo.NewClient(&repoSetting, metrics.Register())
	if err != nil {
		return nil, err
	}

	esbSetting := cc.ConfigServer().Esb
	esbCli, err := esbcli.NewClient(&esbSetting, metrics.Register())
	if err != nil {
		return nil, err
	}

	cs := &ClientSet{
		DS:   pbds.NewDataClient(dsConn),
		Repo: repoCli,
		Esb:  esbCli,
	}

	logs.Infof("initialize the client set success.")
	return cs, nil
}

// ClientSet defines configure server's all the depends api client.
type ClientSet struct {
	// DS data service's client api
	DS pbds.DataClient
	// Repo repoCli repo client api.
	Repo *repo.Client
	// Esb Esb client api
	Esb esbcli.Client
}
