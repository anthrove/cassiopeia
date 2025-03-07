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
	"encoding/json"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

// @Summary	Get an file
// @Tags		CDN API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Router		/api/v1/cdn/{tenant_id}/{file_path} [get]
func (ir IdentityRoutes) cdnGetFile(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	filePath := c.Param("file_path")

	providers, err := ir.service.FindProviders(c, tenantID, object.Pagination{})
	if err != nil {
		return
	}
	sanitizedFilePath := strings.TrimPrefix(filepath.Clean(filePath), "/")

	for _, provider := range providers {
		if provider.Category == "storage" && provider.ProviderType == "local" {
			var parameters map[string]string
			if err := json.Unmarshal(provider.Parameter, &parameters); err != nil {
				return
			}

			sanitizedBasePath := strings.TrimPrefix(filepath.Clean(parameters["base_path"]), "/")
			localFilePath := fmt.Sprintf("local_storage_provider/%s/%s/%s", tenantID, sanitizedBasePath, sanitizedFilePath)

			c.File(localFilePath)
			return
		}
	}
	return
}
