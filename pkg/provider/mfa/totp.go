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
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type totpProperties struct {
	URI    string `json:"uri"`
	Secret string `json:"secret"`
}

type totpBodyData struct {
	OTP string `json:"otp" validate:"required"`
}

type totpProvider struct {
	provider      object.Provider
	period        uint
	digits        otp.Digits
	hashAlgorithm otp.Algorithm
}

func newTOTPProvider(provider object.Provider, period uint, digits otp.Digits, hashAlgorithm otp.Algorithm) (Provider, error) {
	return totpProvider{
		provider:      provider,
		hashAlgorithm: hashAlgorithm,
		period:        period,
		digits:        digits,
	}, nil
}

func (t totpProvider) Create(username string) (object.MFAProviderData, error) {
	if len(username) == 0 {
		return object.MFAProviderData{}, errors.New("username is required")
	}
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.provider.TenantID, //TODO: maybe refactor to use the name instead of the ID
		AccountName: username,
		Period:      t.period,
		Digits:      t.digits,
		Algorithm:   t.hashAlgorithm,
	})

	if err != nil {
		return object.MFAProviderData{}, err
	}

	propertiesJson, err := json.Marshal(totpProperties{
		URI:    secret.URL(),
		Secret: secret.Secret(),
	})

	if err != nil {
		return object.MFAProviderData{}, err
	}

	return object.MFAProviderData{
		Properties: propertiesJson,
	}, nil

}

func (t totpProvider) Validate(mfaConfig json.RawMessage, data map[string]any) (bool, error) {
	var parameters totpBodyData
	var totpProperties totpProperties

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(jsonData, &parameters)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(mfaConfig, &totpProperties)
	if err != nil {
		return false, err
	}

	if len(parameters.OTP) == 0 {
		return false, errors.New("missing otp")
	}

	// Validate OTP
	return totp.Validate(parameters.OTP, totpProperties.Secret), nil
}

func (t totpProvider) GetConfigurationFields() []object.ProviderConfigurationField {
	// There is no configuration that can be returned
	return []object.ProviderConfigurationField{}
}

func (t totpProvider) ValidateConfigurationFields() error {
	// There is no configuration that needs to be validated
	return nil
}
