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
	"database/sql"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"github.com/zitadel/oidc/v3/pkg/op"
	"net/http"
	"time"
)

//	@Summary	Login
//	@Tags		Authentication API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id		path		string							true	"Tenant ID"
//	@Param		application_id	path		string							true	"Application ID"
//	@Param		"Sign In"		body		object.SignInRequest			true	"SignIn Data"
//	@Success	200				{object}	HttpResponse{data=object.User}	"Success"
//	@Failure	400				{object}	HttpResponse{data=nil}			"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/application/{application_id}/login [put]
func (ir IdentityRoutes) signIn(c *gin.Context) {
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

	session, user, err := ir.service.SignIn(c, tenantID, applicationID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	var callbackURL string
	if body.RequestID != "" {
		authRequest, err := ir.service.FindAuthRequest(c, tenantID, body.RequestID)
		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}

		if authRequest.ApplicationID != applicationID {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: "Application ID mismatch",
			})
			return
		}

		err = ir.service.UpdateAuthRequest(c, tenantID, body.RequestID, object.UpdateAuthRequest{
			UserID: sql.NullString{
				String: user.ID,
				Valid:  true,
			},
			Authenticated:   true,
			AuthenticatedAt: time.Now(),
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}

		provider, err := GetProvider(c, ir.service, tenantID)

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}

		callbackURLFn := op.AuthCallbackURL(provider)
		callbackURL = "/" + tenantID + callbackURLFn(c, body.RequestID)
	}

	c.SetCookie("identity_session_id", session, 60*60*24*30, "", "", false, true)

	type resp struct {
		User        object.User `json:"user"`
		RedirectURI string      `json:"redirect_uri,omitempty"`
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: resp{
			user,
			callbackURL,
		},
	})
}

func (ir IdentityRoutes) profile(c *gin.Context) {
	session, exists := c.Get("session")

	if !exists {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: "internal session information are missing",
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: session,
	})
}
