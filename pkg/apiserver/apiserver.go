/*
Copyright 2021 The tKeel Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package apiserver

import (
	"context"
	"net/http"

	entityrouter "github.com/tkeel-io/security/pkg/apirouter/entity/v1"
	oauthrouter "github.com/tkeel-io/security/pkg/apirouter/oauth"
	openapirouter "github.com/tkeel-io/security/pkg/apirouter/openapi/v1"
	rbacrouter "github.com/tkeel-io/security/pkg/apirouter/rbac/v1"
	tenantrouter "github.com/tkeel-io/security/pkg/apirouter/tenant/v1"
	"github.com/tkeel-io/security/pkg/apiserver/config"
	"github.com/tkeel-io/security/pkg/apiserver/filters"
	"github.com/tkeel-io/security/pkg/logger"
	"github.com/tkeel-io/security/pkg/models/dao"

	"github.com/emicklei/go-restful"
)

var (
	_log = logger.NewLogger("auth.apiserver")
)

type APIServer struct {
	Config        *config.Config
	Server        *http.Server
	restContainer *restful.Container
}

func (s *APIServer) Run(ctx context.Context) (err error) {
	dao.SetUp()
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
		_log.Warn("shutdown server...")
	}()

	_log.Infof("start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.restContainer = restful.NewContainer()
	s.restContainer.Router(restful.CurlyRouter{})

	s.installApis()

	for _, webservice := range s.restContainer.RegisteredWebServices() {
		_log.Infof("%s", webservice.RootPath())
	}
	s.Server.Handler = s.restContainer

	return nil
}

func (s *APIServer) installApis() {
	s.restContainer.Filter(filters.GlobalLog())
	must(oauthrouter.AddToRestContainer(s.restContainer))
	must(openapirouter.AddToRestContainer(s.restContainer))
	must(rbacrouter.AddToRestContainer(s.restContainer))
	must(tenantrouter.AddToRestContainer(s.restContainer))
	must(entityrouter.AddToRestContainer(s.restContainer))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
