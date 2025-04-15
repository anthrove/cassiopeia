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
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary	Get an file
// @Tags		CDN API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		file_path	path	string	true	"FilePath"
// @Success	200
// @Router		/api/v1/cdn/{tenant_id}/{file_path} [get]
func (ir IdentityRoutes) cdnGetFile(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	filePath := c.Param("file_path")

	localFilePath, err := ir.service.ServeResource(c, tenantID, filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.File(localFilePath)
}
