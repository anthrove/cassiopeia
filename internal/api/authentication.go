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
	"net/http"

	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
)

// @Summary	Login
// @Tags		Authentication API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path	string	true	"Tenant ID"
// @Param		application_id	path	string	true	"Application ID"
// @Param		username		query	string	true	"Username"
// @Param		type			query	string	true	"Type"
// @Success	200
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application/{application_id}/login/begin [get]
func (ir IdentityRoutes) signInBegin(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	applicationID := c.Param("application_id")

	username := c.Query("username")
	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: "username is a required parameter",
		})
		return
	}

	credentialType := c.Query("type")
	if len(credentialType) == 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: "type is a required parameter",
		})
		return
	}

	data, err := ir.service.SignInStart(c, tenantID, applicationID, object.SignInRequest{
		Username: username,
		Type:     credentialType,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: data,
	})
}

// @Summary	Login
// @Tags		Authentication API
// @Accept		json
// @Produce	json
// @Param		tenant_id		path	string					true	"Tenant ID"
// @Param		application_id	path	string					true	"Application ID"
// @Param		"Sign In"		body	object.SignInRequest	true	"SignIn Data"
// @Success	200
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/application/{application_id}/login [post]
func (ir IdentityRoutes) signInSubmit(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	applicationID := c.Param("application_id")

	var body object.SignInRequest
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	session, user, err := ir.service.SignInSubmit(c, tenantID, applicationID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.SetCookie("identity_session_id", session, 60*60*24*30, "", "", false, true)
	c.JSON(http.StatusOK, HttpResponse{
		Data: user,
	})
}
