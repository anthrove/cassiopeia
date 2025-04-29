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

//	@Summary	Creates a new Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string									true	"Tenant ID"
//	@Param		"Provider"	body		object.CreateProvider					true	"Create Provider Data"
//	@Success	200			{object}	HttpResponse{data=object.Provider{}}	"Provider"
//	@Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider [post]
func (ir IdentityRoutes) createProvider(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateProvider
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	provider, err := ir.service.CreateProvider(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: provider,
	})
}

//	@Summary	Update an existing Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string					true	"Tenant ID"
//	@Param		provider_id	path	string					true	"Provider ID"
//	@Param		"Provider"	body	object.UpdateProvider	true	"Create Provider Data"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/{provider_id} [put]
func (ir IdentityRoutes) updateProvider(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	providerID := c.Param("provider_id")

	var body object.UpdateProvider
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateProvider(c, tenantID, providerID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Kill an existing Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string	true	"Tenant ID"
//	@Param		provider_id	path	string	true	"Provider ID"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/{provider_id} [delete]
func (ir IdentityRoutes) killProvider(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	providerID := c.Param("provider_id")

	err := ir.service.KillProvider(c, tenantID, providerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get an existing Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string									true	"Tenant ID"
//	@Param		provider_id	path		string									true	"Provider ID"
//	@Success	200			{object}	HttpResponse{data=object.Provider{}}	"Provider"
//	@Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/{provider_id} [get]
func (ir IdentityRoutes) findProvider(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	providerID := c.Param("provider_id")

	provider, err := ir.service.FindProvider(c, tenantID, providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: provider,
	})
}

//	@Summary	Get all Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		page		query		string									false	"Page"
//	@Param		page_limit	query		string									false	"Page Limit"
//	@Param		tenant_id	path		string									true	"Tenant ID"
//	@Success	200			{object}	HttpResponse{data=[]object.Provider{}}	"Provider"
//	@Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/ [get]
func (ir IdentityRoutes) findProviders(c *gin.Context) {
	tenantID := c.Param("tenant_id")
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
	provider, err := ir.service.FindProviders(c, tenantID, paginationObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: provider,
	})
}

//	@Summary	Send Mail from Provider
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string	true	"Tenant ID"
//	@Param		provider_id	path	string	true	"Provider ID"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/{provider_id}/mail [post]
func (ir IdentityRoutes) SendMail(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	providerID := c.Param("provider_id")

	var body object.SendMailData
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.SendMail(c, tenantID, providerID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get all Provider Category
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string						true	"Tenant ID"
//	@Success	200			{object}	HttpResponse{data=[]string}	"Provider Categories"
//	@Failure	400			{object}	HttpResponse{data=nil}		"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/category [get]
func (ir IdentityRoutes) findProviderCategories(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	categories, err := ir.service.FindProviderCategories(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: categories,
	})
}

//	@Summary	Get all Provider Category Types
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string						true	"Tenant ID"
//	@Param		category	path		string						true	"Category"
//	@Success	200			{object}	HttpResponse{data=[]string}	"Provider Types"
//	@Failure	400			{object}	HttpResponse{data=nil}		"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/category/{category} [get]
func (ir IdentityRoutes) findProviderCategoryTypes(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	category := c.Param("category")

	providerTypes := ir.service.FindProviderTypes(c, tenantID, category)

	c.JSON(http.StatusOK, HttpResponse{
		Data: providerTypes,
	})
}

//	@Summary	Get all Provider Category Type Configuration fields
//	@Tags		Provider API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string														true	"Tenant ID"
//	@Param		category	path		string														true	"Category"
//	@Param		type		path		string														true	"Provider Type"
//	@Success	200			{object}	HttpResponse{data=[]object.ProviderConfigurationField{}}	"Provider Types"
//	@Failure	400			{object}	HttpResponse{data=nil}										"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/provider/category/{category}/{type} [get]
func (ir IdentityRoutes) findProviderCategoryTypeConfiguration(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	category := c.Param("category")
	providerType := c.Param("type")

	providerTypes := ir.service.FindProviderConfiguration(c, tenantID, category, providerType)

	c.JSON(http.StatusOK, HttpResponse{
		Data: providerTypes,
	})
}
