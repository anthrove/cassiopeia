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

// @Summary	Creates a new Model
// @Tags		Model API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path		string								true	"Tenant ID"
//
// @Param		"Model"		body		object.CreateModel					true	"Create Model Data"
// @Success	200			{object}	HttpResponse{data=object.Model{}}	"Model"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/model [post]
func (ir IdentityRoutes) createModel(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateModel
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	model, err := ir.service.CreateModel(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: model,
	})
}

// @Summary	Update an existing Model
// @Tags		Model API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path	string				true	"Tenant ID"
// @Param		model_id	path	string				true	"Model ID"
//
// @Param		"Model"		body	object.UpdateModel	true	"Create Model Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/model/{model_id} [put]
func (ir IdentityRoutes) updateModel(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	modelID := c.Param("model_id")

	var body object.UpdateModel
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateModel(c, tenantID, modelID, body)

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

// @Summary	Kill an existing Model
// @Tags		Model API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		model_id	path	string	true	"Model ID"
//
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/model/{model_id} [delete]
func (ir IdentityRoutes) killModel(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	modelID := c.Param("model_id")

	err := ir.service.KillModel(c, tenantID, modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing Model
// @Tags		Model API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Param		model_id	path		string								true	"Model ID"
// @Success	200			{object}	HttpResponse{data=object.Model{}}	"Model"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/model/{model_id} [get]
func (ir IdentityRoutes) findModel(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	modelID := c.Param("model_id")

	model, err := ir.service.FindModel(c, tenantID, modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: model,
	})
}

// @Summary	Get existing Models
// @Tags		Model API
// @Accept		json
// @Produce	json
//
// @Param		page		query		string								false	"Page"
// @Param		page_limit	query		string								false	"Page Limit"
//
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.Model{}}	"Model"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/model [get]
func (ir IdentityRoutes) findModels(c *gin.Context) {
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

	models, err := ir.service.FindModels(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: models,
	})
}
