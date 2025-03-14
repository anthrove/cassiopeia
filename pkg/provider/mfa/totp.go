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
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type totpConfiguration struct{}

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

	return object.MFAProviderData{
		Secret: secret.Secret(),
		URI:    secret.URL(),
	}, nil

}

func (t totpProvider) Validate(secret string, data map[string]any) (bool, error) {
	var parameters totpBodyData

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(jsonData, &parameters)
	if err != nil {
		return false, err
	}

	if len(parameters.OTP) == 0 {
		return false, errors.New("missing otp")
	}

	// Validate OTP
	return totp.Validate(parameters.OTP, secret), nil

}

func (t totpProvider) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{}
}

func (t totpProvider) ValidateConfigurationFields() error {
	localConfig := totpConfiguration{}

	err := json.Unmarshal(t.provider.Parameter, &localConfig)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(localConfig)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create totp data"), validateErrs)
		}
	}

	return nil
}
