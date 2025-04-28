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

package email

import (
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"maps"
	"slices"
)

type Provider interface {
	GetConfigurationFields() []object.ProviderConfigurationField
	ValidateConfigurationFields() error
	SendMail(toAddress, subject string, body string) error
}

var providerMap = map[string]func(provider object.Provider) (Provider, error){
	"smtp": newSMTPProvider,
}

func GetEMailProvider(provider object.Provider) (Provider, error) {
	newFunc, exists := providerMap[provider.ProviderType]

	if !exists {
		return nil, errors.New("unknown email provider: " + provider.ProviderType)
	}

	return newFunc(provider)
}

func ConfigurationFields(providerType string) []object.ProviderConfigurationField {
	switch providerType {
	case "smtp":
		return smtpProvider{}.GetConfigurationFields()
	}

	return nil
}

func GetEMailTypes() []string {
	return slices.Collect(maps.Keys(providerMap))
}
