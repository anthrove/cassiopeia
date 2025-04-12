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

// @Summary	Creates a new Permission
// @Tags		Permission API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path		string									true	"Tenant ID"
// @Param		"Permission"	body		object.CreatePermission					true	"Create Permission Data"
// @Success	200				{object}	HttpResponse{data=object.Permission{}}	"Permission"
// @Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/permission [post]
func (ir IdentityRoutes) createPermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreatePermission
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	permission, err := ir.service.CreatePermission(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: permission,
	})
}

// @Summary	Update an existing Permission
// @Tags		Permission API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path	string					true	"Tenant ID"
// @Param		permission_id	path	string					true	"Permission ID"
// @Param		"Permission"	body	object.UpdatePermission	true	"Create Permission Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/permission/{permission_id} [put]
func (ir IdentityRoutes) updatePermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	permissionID := c.Param("permission_id")

	var body object.UpdatePermission
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdatePermission(c, tenantID, permissionID, body)

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

// @Summary	Kill an existing Permission
// @Tags		Permission API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path	string	true	"Tenant ID"
// @Param		permission_id	path	string	true	"Permission ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/permission/{permission_id} [delete]
func (ir IdentityRoutes) killPermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	permissionID := c.Param("permission_id")

	err := ir.service.KillPermission(c, tenantID, permissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing Permission
// @Tags		Permission API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path		string									true	"Tenant ID"
// @Param		permission_id	path		string									true	"Permission ID"
// @Success	200				{object}	HttpResponse{data=object.Permission{}}	"Permission"
// @Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/permission/{permission_id} [get]
func (ir IdentityRoutes) findPermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	permissionID := c.Param("permission_id")

	permission, err := ir.service.FindPermission(c, tenantID, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: permission,
	})
}

// @Summary	Get existing Permissions
// @Tags		Permission API
// @Accept		json
// @Produce	json
// @Param		page		query		string										false	"Page"
// @Param		page_limit	query		string										false	"Page Limit"
// @Param		tenant_id	path		string										true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.Permission{}}	"Permission"
// @Failure	400			{object}	HttpResponse{data=nil}						"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/permission [get]
func (ir IdentityRoutes) findPermissions(c *gin.Context) {
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

	permissions, err := ir.service.FindPermissions(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: permissions,
	})
}
