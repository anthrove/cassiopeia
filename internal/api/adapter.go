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

//	@Summary	Creates a new Adapter
//	@Tags		Adapter API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id	path		string								true	"Tenant ID"
//
//	@Param		"Adapter"	body		object.CreateAdapter				true	"Create Adapter Data"
//	@Success	200			{object}	HttpResponse{data=object.Adapter{}}	"Adapter"
//	@Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/adapter [post]
func (ir IdentityRoutes) createAdapter(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateAdapter
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	adapter, err := ir.service.CreateAdapter(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: adapter,
	})
}

//	@Summary	Update an existing Adapter
//	@Tags		Adapter API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id	path	string					true	"Tenant ID"
//	@Param		adapter_id	path	string					true	"Adapter ID"
//
//	@Param		"Adapter"	body	object.UpdateAdapter	true	"Create Adapter Data"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/adapter/{adapter_id} [put]
func (ir IdentityRoutes) updateAdapter(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	adapterID := c.Param("adapter_id")

	var body object.UpdateAdapter
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateAdapter(c, tenantID, adapterID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: gin.H{},
	})
}

//	@Summary	Kill an existing Adapter
//	@Tags		Adapter API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id	path	string	true	"Tenant ID"
//	@Param		adapter_id	path	string	true	"Adapter ID"
//
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/adapter/{adapter_id} [delete]
func (ir IdentityRoutes) killAdapter(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	adapterID := c.Param("adapter_id")

	err := ir.service.KillAdapter(c, tenantID, adapterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get an existing Adapter
//	@Tags		Adapter API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string								true	"Tenant ID"
//	@Param		adapter_id	path		string								true	"Adapter ID"
//	@Success	200			{object}	HttpResponse{data=object.Adapter{}}	"Adapter"
//	@Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/adapter/{adapter_id} [get]
func (ir IdentityRoutes) findAdapter(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	adapterID := c.Param("adapter_id")

	adapter, err := ir.service.FindAdapter(c, tenantID, adapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: adapter,
	})
}

//	@Summary	Get existing Adapters
//	@Tags		Adapter API
//	@Accept		json
//	@Produce	json
//
//	@Param		page		query		string									false	"Page"
//	@Param		page_limit	query		string									false	"Page Limit"
//
//	@Param		tenant_id	path		string									true	"Tenant ID"
//	@Success	200			{object}	HttpResponse{data=[]object.Adapter{}}	"Adapter"
//	@Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/adapter [get]
func (ir IdentityRoutes) findAdapters(c *gin.Context) {
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

	adapters, err := ir.service.FindAdapters(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: adapters,
	})
}
