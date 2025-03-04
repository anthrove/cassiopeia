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

// @Summary	Creates a new resource
// @Tags		Resource API
// @Accept		json
// @Param		provider_id	query	string	true	"Provider ID"
// @Param		tag			query	string	true	"Tag"
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource [post]
func (ir IdentityRoutes) createResource(c *gin.Context) {
	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}

// @Summary	Deleates a resource
// @Tags		Resource API
// @Accept		json
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource/{resource_id} [delete]
func (ir IdentityRoutes) killResource(c *gin.Context) {
	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}

// @Summary	Lists all resource
// @Tags		Resource API
// @Accept		json
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource/{resource_id} [get]
func (ir IdentityRoutes) findResource(c *gin.Context) {
	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}

// @Summary	Lists all resource
// @Tags		Resource API
// @Accept		json
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource [get]
func (ir IdentityRoutes) findResources(c *gin.Context) {
	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}

// @Summary	Retrieves only the URL of the resource
// @Tags		Resource API
// @Accept		json
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource/{resource_id}/url [get]
func (ir IdentityRoutes) findResourceURL(c *gin.Context) {
	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}
