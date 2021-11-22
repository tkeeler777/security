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
package filters

import (
	"github.com/tkeel-io/security/pkg/apiserver/response"
	"github.com/tkeel-io/security/pkg/errcode"
	"github.com/tkeel-io/security/pkg/models/oauth"

	"github.com/emicklei/go-restful"
)

func Auth() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		operator := oauth.GetOauthOperator()
		if operator == nil {
			_log.Errorf("nil oauth operator")
		}
		token, err := operator.ValidationBearerToken(req.Request)
		if err != nil {
			_log.Error(err)
			response.SrvErrWithRest(resp, errcode.ErrInvalidAccessRequest, nil)
			return
		}

		req.SetAttribute("userID", token.GetUserID())
		chain.ProcessFilter(req, resp)
	}
}