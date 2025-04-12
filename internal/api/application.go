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

// @Summary	Creates a new Application
// @Tags		Application API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path		string									true	"Tenant ID"
// @Param		"Application"	body		object.CreateApplication				true	"Create Application Data"
// @Success	200				{object}	HttpResponse{data=object.Application{}}	"Application"
// @Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application [post]
func (ir IdentityRoutes) createApplication(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateApplication
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	application, err := ir.service.CreateApplication(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: application,
	})
}

// @Summary	Update an existing Application
// @Tags		Application API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path	string						true	"Tenant ID"
// @Param		application_id	path	string						true	"Application ID"
// @Param		"Application"	body	object.UpdateApplication	true	"Create Application Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application/{application_id} [put]
func (ir IdentityRoutes) updateApplication(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	applicationID := c.Param("application_id")

	var body object.UpdateApplication
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateApplication(c, tenantID, applicationID, body)

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

// @Summary	Kill an existing Application
// @Tags		Application API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id		path	string	true	"Tenant ID"
// @Param		application_id	path	string	true	"Application ID"
//
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application/{application_id} [delete]
func (ir IdentityRoutes) killApplication(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	applicationID := c.Param("application_id")

	err := ir.service.KillApplication(c, tenantID, applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing Application
// @Tags		Application API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path		string									true	"Tenant ID"
// @Param		application_id	path		string									true	"Application ID"
// @Success	200				{object}	HttpResponse{data=object.Application{}}	"Application"
// @Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application/{application_id} [get]
func (ir IdentityRoutes) findApplication(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	applicationID := c.Param("application_id")

	application, err := ir.service.FindApplication(c, tenantID, applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: application,
	})
}

// @Summary	Get existing Applications
// @Tags		Application API
// @Accept		json
// @Produce	json
//
// @Param		page		query		string										false	"Page"
// @Param		page_limit	query		string										false	"Page Limit"
//
// @Param		tenant_id	path		string										true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.Application{}}	"Application"
// @Failure	400			{object}	HttpResponse{data=nil}						"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application [get]
func (ir IdentityRoutes) findApplications(c *gin.Context) {
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

	applications, err := ir.service.FindApplications(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: applications,
	})
}
