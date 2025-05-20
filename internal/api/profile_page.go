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

type ProfileFieldValue struct {
	Field object.ProfileField `json:"field"`
	Value any                 `json:"value"`
}

// @Summary	Get Profile fields
// @Tags		Profile API
// @Accept		json
// @Produce	json
// @Failure	200	{object}	HttpResponse{data=[]ProfileFieldValue{}}	"Profile"
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile [get]
func (ir IdentityRoutes) getProfileFields(c *gin.Context) {
	sessionData, exists := c.Get("session")

	if !exists {
		c.JSON(http.StatusInternalServerError, errors.New("this should never happen. Contact an Administrator"))
		return
	}

	sessionObj, ok := sessionData.(map[string]any)

	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("session should be of type session.Session! Contact an Administrator"))
		return
	}

	userData, exists := sessionObj["user"]

	if !exists {
		c.JSON(http.StatusInternalServerError, errors.New("don't get userData. Contact an Administrator"))
		return
	}

	user, ok := userData.(object.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("don't get user. Contact an Administrator"))
		return
	}

	profilePage, err := ir.service.FindProfilePage(c, user.TenantID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	tenant, err := ir.service.FindTenant(c, user.TenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	values := make([]ProfileFieldValue, 0, len(tenant.ProfileFields))
profilePage:
	for _, tenantField := range tenant.ProfileFields {
		for _, profile := range profilePage.Fields {
			if tenantField.Identifier == profile.Identifier {
				values = append(values, ProfileFieldValue{
					Field: tenantField,
					Value: profile.Value,
				})
				continue profilePage
			}
		}

		values = append(values, ProfileFieldValue{
			Field: tenantField,
			Value: nil,
		})
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: values,
	})
}

// @Summary	Upsert Profile fields
// @Tags		Profile API
// @Accept		json
// @Produce	json
// @Param		"Profile Pages"	body	object.UpdateProfilePage	true	"Create Provider Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/profile [put]
func (ir IdentityRoutes) upsertProfileFields(c *gin.Context) {
	sessionData, exists := c.Get("session")

	if !exists {
		c.JSON(http.StatusInternalServerError, errors.New("this should never happen. Contact an Administrator"))
		return
	}

	sessionObj, ok := sessionData.(map[string]any)

	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("session should be of type session.Session! Contact an Administrator"))
		return
	}

	userData, exists := sessionObj["user"]

	if !exists {
		c.JSON(http.StatusInternalServerError, errors.New("don't get userData. Contact an Administrator"))
		return
	}

	user, ok := userData.(object.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("don't get user. Contact an Administrator"))
		return
	}

	_, err := ir.service.FindProfilePage(c, user.TenantID, user.ID)

	if err == nil {
		// Update
		var body object.UpdateProfilePage
		err := c.ShouldBind(&body)

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}

		err = ir.service.UpdateProfilePage(c, user.TenantID, user.ID, body)

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}
	} else {
		// Create
		var body object.CreateProfilePage
		err := c.ShouldBind(&body)

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}

		_, err = ir.service.CreateProfilePage(c, user.TenantID, user.ID, body)

		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{
				Error: err.Error(),
			})
			return
		}
	}

	c.Status(http.StatusNoContent)
}
