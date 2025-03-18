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
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/pquerna/otp"
)

type Provider interface {
	GetConfigurationFields() []object.ProviderConfigurationField
	ValidateConfigurationFields() error
	Create(username string) (object.MFAProviderData, error)
	Validate(secret string, data map[string]any) (bool, error)
	ValidateMFAMethode(secret string, data map[string]any) (bool, error)
}

func GetMFAProvider(provider object.Provider) (Provider, error) {
	switch provider.ProviderType {
	case "totp":
		// AlgorithmSHA1 should be used for compatibility with Google Authenticator, See https://github.com/pquerna/otp/issues/55 for additional details.
		return newTOTPProvider(provider, 30, otp.DigitsSix, otp.AlgorithmSHA1)
	}
	return nil, errors.New("unknown mfa provider: " + provider.ProviderType)
}
