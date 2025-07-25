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

//	@Summary	Creates a new MFA for a User
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string							true	"Tenant ID"
//	@Param		user_id		path		string							true	"User ID"
//	@Param		MFA			body		object.CreateMFA				true	"Create MFA Data"
//	@Success	200			{object}	HttpResponse{data=object.MFA{}}	"MFA"
//	@Failure	400			{object}	HttpResponse{data=nil}			"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa [post]
func (ir IdentityRoutes) createMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	var body object.CreateMFA
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	mfa, err := ir.service.CreateMFA(c, tenantID, userID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: mfa,
	})
}

//	@Summary	Update an existing MFA for a User
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string				true	"Tenant ID"
//	@Param		user_id		path	string				true	"User ID"
//	@Param		mfa_id		path	string				true	"MFA ID"
//	@Param		MFA			body	object.UpdateMFA	true	"Update MFA Data"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [put]
func (ir IdentityRoutes) updateMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")
	mfaID := c.Param("mfa_id")

	var body object.UpdateMFA
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateMFA(c, tenantID, userID, mfaID, body)

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

//	@Summary	Delete an existing MFA for User
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string	true	"Tenant ID"
//	@Param		user_id		path	string	true	"User ID"
//	@Param		mfa_id		path	string	true	"MFA ID"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [delete]
func (ir IdentityRoutes) killMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")
	mfaID := c.Param("mfa_id")

	err := ir.service.KillMFA(c, tenantID, userID, mfaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get an existing MFA for User
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path		string							true	"Tenant ID"
//	@Param		user_id		path		string							true	"User ID"
//	@Param		mfa_id		path		string							true	"MFA ID"
//	@Success	200			{object}	HttpResponse{data=object.MFA{}}	"MFA"
//	@Failure	400			{object}	HttpResponse{data=nil}			"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [get]
func (ir IdentityRoutes) findMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")
	mfaID := c.Param("mfa_id")

	mfa, err := ir.service.FindMFA(c, tenantID, userID, mfaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: mfa,
	})
}

//	@Summary	Get existing MFAs
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		page		query		string								false	"Page"
//	@Param		page_limit	query		string								false	"Page Limit"
//	@Param		tenant_id	path		string								true	"Tenant ID"
//	@Param		user_id		path		string								true	"User ID"
//	@Success	200			{object}	HttpResponse{data=[]object.MFA{}}	"MFA"
//	@Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa [get]
func (ir IdentityRoutes) findMFAs(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

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

	mfas, err := ir.service.FindMFAs(c, tenantID, userID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: mfas,
	})
}

//	@Summary	Validates a given MFA
//	@Tags		MFA API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id	path	string			true	"Tenant ID"
//	@Param		user_id		path	string			true	"User ID"
//	@Param		mfa_id		path	string			true	"MFA ID"
//	@Param		MFA			body	map[string]any	true	"Verify MFA Body"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id}/verify [post]
func (ir IdentityRoutes) verifyMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")
	mfaID := c.Param("mfa_id")

	var body map[string]any
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.VerifyMFA(c, tenantID, userID, mfaID, body)

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
