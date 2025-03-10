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

package logic

import (
	"context"
	"errors"
	"github.com/anthrove/identity/pkg/crypto"
	"github.com/anthrove/identity/pkg/object"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (is IdentityService) SignIn(ctx context.Context, tenantID string, applicationID string, signInData object.SignInRequest) (string, object.User, error) {
	// get for making sure application exists and also later its planed you can enable and disable signin and stuff
	_, err := is.FindApplication(ctx, tenantID, applicationID)

	if err != nil {
		return "", object.User{}, err
	}

	user, err := is.FindUserByUsername(ctx, tenantID, signInData.Username)

	if err != nil {
		return "", object.User{}, err
	}

	if !user.EmailVerified {
		return "", object.User{}, errors.New("user is not verified")
	}

	hasher, err := crypto.GetPasswordHasher(user.PasswordType)

	if err != nil {
		return "", object.User{}, err
	}

	password, err := hasher.ComparePassword(signInData.Password, user.PasswordHash, user.PasswordSalt)
	if err != nil {
		return "", object.User{}, err
	}

	if !password {
		return "", object.User{}, errors.New("password is incorrect")
	}

	sessionID, err := gonanoid.New(50)

	if err != nil {
		return "", object.User{}, err
	}

	session := is.FindSession(ctx, sessionID)
	session["tenant_id"] = tenantID
	session["application_id"] = applicationID
	session["user"] = user
	session["logged_in"] = true

	is.UpdateSession(ctx, sessionID, session)

	return sessionID, user, nil
}
