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

// @Summary	Creates a MFA from a profile
// @Tags		Profile API
// @Accept		json
// @Produce	json
// @Failure	200	{object}	HttpResponse{data=object.MFA{}}	"Profile"
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa [post]
func (ir IdentityRoutes) profileCreateMFA(c *gin.Context) {
	user, err := sessionConvert(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var body object.CreateMFA
	err = c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	mfa, err := ir.service.CreateMFA(c, user.TenantID, user.ID, body)
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

// @Summary	Verify MFA from a profile
// @Tags		Profile API
// @Accept		json
// @Produce	json
// @Failure	200	{object}	HttpResponse{data=nil}	"Profile"
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa/{mfa_id}/verify [post]
func (ir IdentityRoutes) profileVerifyMFA(c *gin.Context) {
	mfaID := c.Param("mfa_id")
	user, err := sessionConvert(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var body map[string]any
	err = c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.VerifyMFA(c, user.TenantID, user.ID, mfaID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: nil,
	})
}

// @Summary	Update MFA from a profile
// @Tags		Profile API
// @Accept		json
// @Produce	json
// @Failure	200	{object}	HttpResponse{data=nil}	"Profile"
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa/{mfa_id} [post]
func (ir IdentityRoutes) profileUpdateMFA(c *gin.Context) {
	mfaID := c.Param("mfa_id")

	user, err := sessionConvert(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var body object.UpdateMFA
	err = c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateMFA(c, user.TenantID, user.ID, mfaID, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: nil,
	})
}

// @Summary	List all MFA from a profile
// @Tags		Profile API
// @Accept		json
// @Produce	json
//
//	@Param		page		query		string									false	"Page"
//	@Param		page_limit	query		string									false	"Page Limit"
//
// @Failure	200	{object}	HttpResponse{data=[]object.MFA{}}	"Get all Profiles"
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa [get]
func (ir IdentityRoutes) profileGetMFAs(c *gin.Context) {
	user, err := sessionConvert(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

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

	mfas, err := ir.service.FindMFAs(c, user.TenantID, user.ID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: mfas,
	})
}

// @Summary	List all MFA from a profile
// @Tags		Profile API
// @Accept		json
// @Produce	json
//
//	@Param		page		query		string									false	"Page"
//	@Param		page_limit	query		string									false	"Page Limit"
//
// @Failure	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa [get]
func (ir IdentityRoutes) profileKillMFAs(c *gin.Context) {
	mfaID := c.Param("mfa_id")
	user, err := sessionConvert(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = ir.service.KillMFA(c, user.TenantID, user.ID, mfaID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
