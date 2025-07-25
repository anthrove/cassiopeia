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
	"strconv"
	"strings"
)

// Pagination is a middleware function for handling pagination in HTTP requests.
// It extracts the page and page_limit query parameters, validates them, and sets a Pagination object in the context.
//
// Returns:
//   - A gin.HandlerFunc that can be used as middleware in a Gin router.
func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		if page == "" {
			err := errors.New("page parameter is missing")
			c.AbortWithError(http.StatusBadRequest, err)
		}

		pageLimit := c.DefaultQuery("page_limit", "50")
		if pageLimit == "" {
			err := errors.New("page limit parameter is missing")
			c.AbortWithError(http.StatusBadRequest, err)
		}

		pageLimitInt, err := strconv.Atoi(pageLimit)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		pagination := object.Pagination{
			Limit: pageLimitInt,
			Page:  pageInt,
		}

		c.Set("pagination", pagination)
		c.Next()
	}
}

func (ir IdentityRoutes) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("identity_session_id")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if sessionID == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session := ir.service.FindSession(c, sessionID)
		if session == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("session", session)

		tenantID := c.Param("tenant_id")

		if tenantID == "" {
			// Endpoint without tenant id!!!
			c.Next()
			return
		}

		user := session["user"].(object.User)

		access, err := ir.service.Enforce(c, tenantID, tenantID, []any{
			strings.Trim(user.TenantID, " "),
			user.ID,
			c.Request.URL.Path,
			strings.ToLower(c.Request.Method),
		})

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !access {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
