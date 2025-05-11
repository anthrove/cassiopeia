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

// @Summary	Creates a new Enforcer
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string									true	"Tenant ID"
// @Param		"Enforcer"	body		object.CreateEnforcer					true	"Create Enforcer Data"
// @Success	200			{object}	HttpResponse{data=object.Enforcer{}}	"Enforcer"
// @Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer [post]
func (ir IdentityRoutes) createEnforcer(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateEnforcer
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	enforcer, err := ir.service.CreateEnforcer(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: enforcer,
	})
}

// @Summary	Update an existing Enforcer
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string					true	"Tenant ID"
// @Param		enforcer_id	path	string					true	"Enforcer ID"
// @Param		"Enforcer"	body	object.UpdateEnforcer	true	"Create Enforcer Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer/{enforcer_id} [put]
func (ir IdentityRoutes) updateEnforcer(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	enforcerID := c.Param("enforcer_id")

	var body object.UpdateEnforcer
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateEnforcer(c, tenantID, enforcerID, body)

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

// @Summary	Kill an existing Enforcer
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		enforcer_id	path	string	true	"Enforcer ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer/{enforcer_id} [delete]
func (ir IdentityRoutes) killEnforcer(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	enforcerID := c.Param("enforcer_id")

	err := ir.service.KillEnforcer(c, tenantID, enforcerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing Enforcer
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string									true	"Tenant ID"
// @Param		enforcer_id	path		string									true	"Enforcer ID"
// @Success	200			{object}	HttpResponse{data=object.Enforcer{}}	"Enforcer"
// @Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer/{enforcer_id} [get]
func (ir IdentityRoutes) findEnforcer(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	enforcerID := c.Param("enforcer_id")

	enforcer, err := ir.service.FindEnforcer(c, tenantID, enforcerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: enforcer,
	})
}

// @Summary	Get existing Enforcers
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		page		query		string									false	"Page"
// @Param		page_limit	query		string									false	"Page Limit"
// @Param		tenant_id	path		string									true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.Enforcer{}}	"Enforcer"
// @Failure	400			{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer [get]
func (ir IdentityRoutes) findEnforcers(c *gin.Context) {
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

	enforcers, err := ir.service.FindEnforcers(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: enforcers,
	})
}

// @Summary	Check if the request has Permissions
// @Tags		Enforcer API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		enforcer_id	path	string	true	"Enforcer ID"
// @Param		"Enforcer"	body	[]any	true	"Create Enforcer Data"
// @Success	200
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/enforcer/{enforcer_id}/enforce [post]
func (ir IdentityRoutes) enforce(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	enforcerID := c.Param("enforcer_id")

	var body []any
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	success, err := ir.service.Enforce(c, tenantID, enforcerID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: success,
	})
}
