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

package auth

import (
	"context"
	"errors"
	"github.com/anthrove/identity/pkg/object"
)

type ProviderContext struct {
	Tenant     object.Tenant
	User       object.User
	Credential object.Credentials

	SendMail func(ctx context.Context, data object.SendMailData)
}

type Provider interface {
	GetConfigurationFields() []object.ProviderConfigurationField
	ValidateConfigurationFields() error
	// Configure is used for set up a new credential object. It returns the finished metadata which can be saved.
	Configure(ctx context.Context, providerContext ProviderContext, data map[string]any) (map[string]any, error)
	Validate(ctx context.Context, providerContext ProviderContext, data map[string]any) (bool, map[string]any, error)
	Begin(ctx context.Context, providerContext ProviderContext) (map[string]any, error)
	Submit(ctx context.Context, providerContext ProviderContext, data map[string]any) (bool, error)
}

func GetAuthProvider(provider object.Provider) (Provider, error) {
	switch provider.ProviderType {
	case "password":
		return newPasswordAuth(provider), nil
	}
	return nil, errors.New("unknown auth provider: " + provider.ProviderType)
}
