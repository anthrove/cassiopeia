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
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"net/http"
)

// createTenant handles the HTTP POST request to create a new tenant.
// It binds the request body to a CreateTenant object, validates it, and calls the service to create the tenant.
// If successful, it returns the created tenant; otherwise, it returns an error response.
//
// Parameters:
//
//   - c: a gin.Context instance representing the context of the HTTP request.
//
//	@Summary	Creates a new Tenant
//	@Tags		Tenant API
//	@Accept		json
//	@Produce	json
//	@Param		"Tenant"	body		object.CreateTenant					true	"Create Tenant Data"
//	@Success	200			{object}	HttpResponse{data=object.Tenant{}}	"Tenant"
//	@Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
//	@Router		/api/v1/tenant [post]
func (ir IdentityRoutes) createTenant(c *gin.Context) {
	var body object.CreateTenant
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	tenant, err := ir.service.CreateTenant(c, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: tenant,
	})
}

// updateTenant handles the HTTP PUT request to update an existing tenant.
// It binds the request body to an UpdateTenant object, validates it, and calls the service to update the tenant.
// If successful, it returns an empty response; otherwise, it returns an error response.
//
// Parameters:
//
//   - c: a gin.Context instance representing the context of the HTTP request.
//
//	@Summary	Update an existing Tenant
//	@Tags		Tenant API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string				true	"Tenant ID"
//	@Param		"Tenant"	body	object.UpdateTenant	true	"Update Tenant Data"
//	@Success	200
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id} [put]
func (ir IdentityRoutes) updateTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.UpdateTenant
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateTenant(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

//	@Summary	Delete a Tenant
//	@Tags		Tenant API
//	@Accept		json
//	@Produce	json
//	@Param		id	body	string	true	"Tenant ID"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id} [delete]
func (ir IdentityRoutes) killTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	err := ir.service.KillTenant(c, tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get a Tenant by ID
//	@Tags		Tenant API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string								true	"Tenant ID"
//	@Success	200			{object}	HttpResponse{data=object.Tenant{}}	"Tenant"
//	@Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id} [get]
func (ir IdentityRoutes) findTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	tenant, err := ir.service.FindTenant(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: tenant,
	})
}

//	@Summary	Get all Tenants
//	@Tags		Tenant API
//	@Accept		json
//	@Produce	json
//	@Param		page		query		string									false	"Page"
//	@Param		page_limit	query		string									false	"Page Limit"
//	@Success	200			{object}	HttpResponse{data=[]object.Tenant{}}	"Tenant"
//	@Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant [get]
func (ir IdentityRoutes) findTenants(c *gin.Context) {
	pagination, ok := c.Get("pagination")
	if !ok {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: "pagination parameter is missing",
		})
		return
	}

	paginationObj, ok := pagination.(object.Pagination)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("pagination parameter cant be converted to object.Pagination"))
		return
	}

	tenants, err := ir.service.FindTenants(c, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: tenants,
	})
}
