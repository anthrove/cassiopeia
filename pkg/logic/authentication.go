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
	"database/sql"
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/auth"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

func (is IdentityService) SignInStart(ctx context.Context, tenantID string, applicationID string, signInData object.SignInRequest) (map[string]any, error) {
	tenant, err := is.FindTenant(ctx, tenantID)

	if err != nil {
		return nil, err
	}

	// get for making sure application exists and also later its planed you can enable and disable signin and stuff
	application, err := is.FindApplication(ctx, tenantID, applicationID)

	if err != nil {
		return nil, err
	}

	user, err := is.FindUserByUsername(ctx, tenantID, signInData.Username)

	if err != nil {
		return nil, err
	}

	userCredentials, err := is.FindCredentialsByUser(ctx, tenantID, user.ID)

	if err != nil {
		return nil, err
	}

	var selectedCredential object.Credentials
	for _, userCred := range userCredentials {
		if signInData.Type == userCred.Type {
			selectedCredential = userCred
			break
		}
	}

	if len(selectedCredential.ID) == 0 {
		return nil, errors.New("no configured credential found")
	}

	var authProviderObj object.Provider
	for _, provider := range application.AuthProvider {
		if signInData.Type == provider.ProviderType {
			authProviderObj = provider
			break
		}
	}

	if len(authProviderObj.ID) == 0 {
		return nil, errors.New("no provider was configured with given type")
	}

	authProvider, err := auth.GetAuthProvider(authProviderObj)

	if err != nil {
		return nil, err
	}

	return authProvider.Begin(ctx, auth.ProviderContext{
		Tenant:     tenant,
		User:       user,
		Credential: selectedCredential,
		SendMail: func(ctx context.Context, data object.SendMailData) {
			// TODO do nothing  right now.. we need to check to fix it
		},
	})
}

func (is IdentityService) SignInSubmit(ctx context.Context, tenantID string, applicationID string, signInData object.SignInRequest) (string, object.User, error) {
	tenant, err := is.FindTenant(ctx, tenantID)

	if err != nil {
		return "", object.User{}, err
	}

	// get for making sure application exists and also later its planed you can enable and disable signin and stuff
	application, err := is.FindApplication(ctx, tenantID, applicationID)

	if err != nil {
		return "", object.User{}, err
	}

	user, err := is.FindUserByUsername(ctx, tenantID, signInData.Username)

	if err != nil {
		return "", object.User{}, err
	}

	userCredentials, err := is.FindCredentialsByUser(ctx, tenantID, user.ID)

	if err != nil {
		return "", object.User{}, err
	}

	var selectedCredential object.Credentials
	for _, userCred := range userCredentials {
		if signInData.Type == userCred.Type {
			selectedCredential = userCred
			break
		}
	}

	if len(selectedCredential.ID) == 0 {
		return "", object.User{}, errors.New("no configured credential found")
	}

	var authProviderObj object.Provider
	for _, provider := range application.AuthProvider {
		if signInData.Type == provider.ProviderType {
			authProviderObj = provider
			break
		}
	}

	if len(authProviderObj.ID) == 0 {
		return "", object.User{}, errors.New("no provider was configured with given type")
	}

	authProvider, err := auth.GetAuthProvider(authProviderObj)

	if err != nil {
		return "", object.User{}, err
	}

	success, err := authProvider.Submit(ctx, auth.ProviderContext{
		Tenant:     tenant,
		User:       user,
		Credential: selectedCredential,
		SendMail: func(ctx context.Context, data object.SendMailData) {
			// TODO do nothing  right now.. we need to check to fix it
		},
	}, signInData.Metadata)

	if !success {
		return "", object.User{}, errors.New("credential were incorrect")
	}

	sessionID, err := gonanoid.New(50)

	if err != nil {
		return "", object.User{}, err
	}

	session := is.FindSession(ctx, sessionID)

	if session == nil {
		session = map[string]any{}
	}

	session["tenant_id"] = tenantID
	session["application_id"] = applicationID
	session["user"] = user
	session["logged_in"] = true

	is.UpdateSession(ctx, sessionID, session)

	if signInData.RequestID != "" {
		err := is.UpdateAuthRequest(ctx, tenantID, signInData.RequestID, object.UpdateAuthRequest{
			UserID: sql.NullString{
				String: user.ID,
				Valid:  true,
			},
			Authenticated:   true,
			AuthenticatedAt: time.Now(),
		})

		if err != nil {
			return "", object.User{}, err
		}
	}

	return sessionID, user, nil
}
