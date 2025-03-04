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

// createMessageTemplate handles the HTTP POST request to create a new tenant.
// It binds the request body to a CreateMessageTemplate object, validates it, and calls the service to create the tenant.
// If successful, it returns the created tenant; otherwise, it returns an error response.
//
// Parameters:
//
//   - c: a gin.Context instance representing the context of the HTTP request.
//
//     @Summary	Creates a new MessageTemplate
//     @Tags		MessageTemplate API
//     @Accept		json
//     @Produce	json
//
//     @Param		tenant_id			path		string										true	"Tenant ID"
//     @Param		"MessageTemplate"	body		object.CreateMessageTemplate				true	"Create MessageTemplate Data"
//     @Success	200					{object}	HttpResponse{data=object.MessageTemplate{}}	"MessageTemplate"
//     @Failure	400					{object}	HttpResponse{data=nil}						"Bad Request"
//     @Router		/api/v1/tenant/{tenant_id}/template [post]
func (ir IdentityRoutes) createMessageTemplate(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateMessageTemplate
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	tenant, err := ir.service.CreateMessageTemplate(c, tenantID, body)

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

// updateMessageTemplate handles the HTTP PUT request to update an existing tenant.
// It binds the request body to an UpdateMessageTemplate object, validates it, and calls the service to update the tenant.
// If successful, it returns an empty response; otherwise, it returns an error response.
//
// Parameters:
//
//   - c: a gin.Context instance representing the context of the HTTP request.
//
//     @Summary	Update an existing MessageTemplate
//     @Tags		MessageTemplate API
//     @Accept		json
//     @Produce	json
//     @Param		tenant_id			path	string							true	"Tenant ID"
//     @Param		template_id			path	string							true	"MessageTemplate ID"
//     @Param		"MessageTemplate"	body	object.UpdateMessageTemplate	true	"Update MessageTemplate Data"
//     @Success	200
//     @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//     @Router		/api/v1/tenant/{tenant_id}/template/{template_id} [put]
func (ir IdentityRoutes) updateMessageTemplate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	templateID := c.Param("template_id")

	var body object.UpdateMessageTemplate
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateMessageTemplate(c, tenantID, templateID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary	Delete a MessageTemplate
// @Tags		MessageTemplate API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		template_id	path	string	true	"MessageTemplate ID"
//
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/template/{template_id} [delete]
func (ir IdentityRoutes) killMessageTemplate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	templateID := c.Param("template_id")

	err := ir.service.KillMessageTemplate(c, tenantID, templateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get a MessageTemplate by ID
// @Tags		MessageTemplate API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path		string										true	"Tenant ID"
// @Param		template_id	path		string										true	"MessageTemplate ID"
//
// @Success	200			{object}	HttpResponse{data=object.MessageTemplate{}}	"MessageTemplate"
// @Failure	400			{object}	HttpResponse{data=nil}						"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/template/{template_id} [get]
func (ir IdentityRoutes) findMessageTemplate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	templateID := c.Param("template_id")

	tenant, err := ir.service.FindMessageTemplate(c, tenantID, templateID)
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

//	@Summary	Get all MessageTemplates
//	@Tags		MessageTemplate API
//	@Accept		json
//	@Produce	json
//	@Param		page		query	string	false	"Page"
//	@Param		page_limit	query	string	false	"Page Limit"
//	@Param		tenant_id	path	string	true	"Tenant ID"
//	@Param		template_id	path	string	true	"MessageTemplate ID"

// @Success	200	{object}	HttpResponse{data=[]object.MessageTemplate{}}	"MessageTemplate"
// @Failure	400	{object}	HttpResponse{data=nil}							"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/template [get]
func (ir IdentityRoutes) findMessageTemplates(c *gin.Context) {
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

	tenants, err := ir.service.FindMessageTemplates(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: tenants,
	})
}
