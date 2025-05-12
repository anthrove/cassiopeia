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
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		provider_id	path	string	true	"Provider ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile/mfa [post]
func (ir IdentityRoutes) ProfileCreateMFA(c *gin.Context) {
	sessionData, exists := c.Get("session")

	if !exists {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("this should never happen. Contact an Administrator"))
		return
	}

	sessionObj, ok := sessionData.(map[string]any)

	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("session should be of type session.Session! Contact an Administrator"))
		return
	}

	userData, exists := sessionObj["user"]

	if !exists {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("don't get userData. Contact an Administrator"))
		return
	}

	user, ok := userData.(object.User)

	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("don't get user. Contact an Administrator"))
		return
	}

	var body object.CreateMFA
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	mfa, err := ir.service.CreateMFA(c, user.TenantID, user.TenantID, body)
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
