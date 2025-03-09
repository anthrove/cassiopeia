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
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary	Creates a new resource
// @Tags		Resource API
// @Accept		json
// @Param		tenant_id	path		string									true	"Tenant ID"
// @Param		provider_id	query		string									true	"Provider ID"
// @Param		tag			query		string									true	"Tag"
// @Success	200			{object}	HttpResponse{data=object.Resource{}}	"Resource"
// @Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource [post]
func (ir IdentityRoutes) createResource(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	providerID := c.Query("provider_id")
	tag := c.Query("tag")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	createResource := object.CreateResource{
		ProviderID: providerID,
		Tag:        tag,
		FileSize:   file.Size,
		FileName:   file.Filename,
		MimeType:   file.Header.Get("Content-Type"),
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}
	defer fileContent.Close()

	resource, err := ir.service.CreateResource(c, tenantID, createResource, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: resource,
	})
}

// @Summary	Deleates a resource
// @Tags		Resource API
// @Accept		json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		resource_id	query	string	true	"Resource ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource/{resource_id} [delete]
func (ir IdentityRoutes) killResource(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	resourceID := c.Param("resource_id")

	err := ir.service.KillResource(c, tenantID, resourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary	Lists all resource
// @Tags		Resource API
// @Accept		json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		resource_id	query	string	true	"Resource ID"
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource/{resource_id} [get]
func (ir IdentityRoutes) findResource(c *gin.Context) {
	// tenantID := c.Param("tenant_id")
	// resourceID := c.Query("resource_id")

	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}

// @Summary	Lists all resource
// @Tags		Resource API
// @Accept		json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Produce	json
// @Router		/api/v1/tenant/{tenant_id}/resource [get]
func (ir IdentityRoutes) findResources(c *gin.Context) {
	// tenantID := c.Param("tenant_id")
	// resourceID := c.Query("resource_id")

	// TODO: Implement Me
	c.JSON(http.StatusForbidden, "Not Implemented yet")
}
