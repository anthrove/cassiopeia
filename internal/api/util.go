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
)

func sessionConvert(c *gin.Context) (object.User, error) {
	sessionData, exists := c.Get("session")

	if !exists {
		return object.User{}, errors.New("this should never happen. Contact an Administrator")
	}

	sessionObj, ok := sessionData.(map[string]any)

	if !ok {
		return object.User{}, errors.New("session should be of type session.Session! Contact an Administrator")
	}

	userData, exists := sessionObj["user"]

	if !exists {
		return object.User{}, errors.New("don't get userData. Contact an Administrator")
	}

	user, ok := userData.(object.User)

	if !ok {
		return object.User{}, errors.New("don't get userData. Contact an Administrator")
	}

	return user, nil
}
