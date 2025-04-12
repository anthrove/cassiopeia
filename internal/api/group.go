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

// @Summary	Creates a new Group
// @Tags		Group API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Param		"Group"		body		object.CreateGroup					true	"Create Group Data"
// @Success	200			{object}	HttpResponse{data=object.Group{}}	"Group"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/group [post]
func (ir IdentityRoutes) createGroup(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateGroup
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	group, err := ir.service.CreateGroup(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: group,
	})
}

// @Summary	Update an existing Group
// @Tags		Group API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string				true	"Tenant ID"
// @Param		group_id	path	string				true	"Group ID"
// @Param		"Group"		body	object.UpdateGroup	true	"Create Group Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/group/{group_id} [put]
func (ir IdentityRoutes) updateGroup(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	groupID := c.Param("group_id")

	var body object.UpdateGroup
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateGroup(c, tenantID, groupID, body)

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

// @Summary	Kill an existing Group
// @Tags		Group API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		group_id	path	string	true	"Group ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/group/{group_id} [delete]
func (ir IdentityRoutes) killGroup(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	groupID := c.Param("group_id")

	err := ir.service.KillGroup(c, tenantID, groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing Group
// @Tags		Group API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Param		group_id	path		string								true	"Group ID"
// @Success	200			{object}	HttpResponse{data=object.Group{}}	"Group"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/group/{group_id} [get]
func (ir IdentityRoutes) findGroup(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	groupID := c.Param("group_id")

	group, err := ir.service.FindGroup(c, tenantID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: group,
	})
}

// @Summary	Get existing Groups
// @Tags		Group API
// @Accept		json
// @Produce	json
// @Param		page		query		string								false	"Page"
// @Param		page_limit	query		string								false	"Page Limit"
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.Group{}}	"Group"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/group [get]
func (ir IdentityRoutes) findGroups(c *gin.Context) {
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

	groups, err := ir.service.FindGroups(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: groups,
	})
}
