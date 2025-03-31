/*
 * Copyright (C) 2025 Anthrove
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"github.com/anthrove/identity/pkg/logic"
	"github.com/anthrove/identity/pkg/oidc"
	"github.com/gin-gonic/gin"
	"github.com/zitadel/oidc/v3/pkg/op"
	"net/http"
	"sync"
)

var (
	providers = make(map[string]*op.Provider)
	lock      sync.Mutex
)

func GetProvider(service logic.IdentityService, tenantID string) (*op.Provider, error) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := providers[tenantID]; !ok {
		storage := oidc.NewStorage(service, tenantID)
		provider, err := oidc.NewProvider(storage, tenantID)

		if err != nil {
			return nil, err
		}

		providers[tenantID] = provider
	}

	return providers[tenantID], nil
}

func (ir IdentityRoutes) OIDCEndpoints(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	provider, err := GetProvider(ir.service, tenantID)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	request := c.Request
	request.URL.Path = c.Request.URL.Path[26:]

	provider.ServeHTTP(c.Writer, c.Request)
}
