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

package mfa

import (
	"encoding/json"
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"maps"
	"slices"
)

type Provider interface {
	GetConfigurationFields() []object.ProviderConfigurationField
	ValidateConfigurationFields() error
	GenerateUserConfig(username string) (object.MFAProviderData, error)
	InitDataFlow(mfaConfig json.RawMessage) (map[string]any, error)
	ValidateDatFlow(mfaConfig json.RawMessage, data map[string]any) (bool, error)
}

var providerMap = map[string]func(provider object.Provider) (Provider, error){
	"totp": newTOTPProvider,
}

func GetMFAProvider(provider object.Provider) (Provider, error) {
	newFunc, exists := providerMap[provider.ProviderType]

	if !exists {
		return nil, errors.New("unknown mfa provider: " + provider.ProviderType)
	}

	return newFunc(provider)
}

func ConfigurationFields(providerType string) []object.ProviderConfigurationField {
	switch providerType {
	case "totp":
		return totpProvider{}.GetConfigurationFields()
	}

	return nil
}

func GetMfaTypes() []string {
	return slices.Collect(maps.Keys(providerMap))
}
