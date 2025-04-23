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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/crypto"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-playground/validator/v10"
	"strconv"
	"unicode"
)

type passwordConfiguration struct {
	MinPasswordLength  int `json:"min_password_length"`
	MaxPasswordLength  int `json:"max_password_length"`
	MinLowercaseLength int `json:"min_lowercase_length"`
	MinUppercaseLetter int `json:"min_uppercase_letter"`
	MinDigitLetter     int `json:"min_digit_letter"`
	MinSpecialLetter   int `json:"min_special_letter"`
}

type passwordAuth struct {
	provider object.Provider
}

func newPasswordAuth(provider object.Provider) Provider {
	return &passwordAuth{
		provider: provider,
	}
}

func (p passwordAuth) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{
		{
			FieldKey:  "min_password_length",
			FieldType: "int",
		},
		{
			FieldKey:  "max_password_length",
			FieldType: "int",
		},
		{
			FieldKey:  "min_lowercase_letter",
			FieldType: "int",
		},
		{
			FieldKey:  "min_uppercase_letter",
			FieldType: "int",
		},
		{
			FieldKey:  "min_digit_letter",
			FieldType: "int",
		},
		{
			FieldKey:  "min_special_letter",
			FieldType: "int",
		},
	}
}

func (p passwordAuth) ValidateConfigurationFields() error {
	passwordConfig := passwordConfiguration{}

	err := json.Unmarshal(p.provider.Parameter, &passwordConfig)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(passwordConfig)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating password configuration data"), validateErrs)
		}
	}

	return nil
}

func (p passwordAuth) Configure(_ context.Context, providerContext ProviderContext, data map[string]any) (map[string]any, error) {
	passwordSalt, err := util.RandomSaltString(25)
	if err != nil {
		return nil, err
	}

	passwordHasher, err := crypto.GetPasswordHasher(providerContext.Tenant.PasswordType)
	if err != nil {
		return nil, err
	}

	password, exist := data["password"]

	if !exist {
		return nil, errors.New("password is a required data field")
	}

	passwordStr, ok := password.(string)

	if !ok {
		return nil, errors.New("password is not a string")
	}

	passwordConfig := passwordConfiguration{}

	err = json.Unmarshal(p.provider.Parameter, &passwordConfig)
	if err != nil {
		return nil, err
	}

	if len(passwordStr) >= passwordConfig.MinPasswordLength {
		return nil, errors.New("password is too short")
	}

	if len(passwordStr) >= passwordConfig.MaxPasswordLength {
		return nil, errors.New("password is too long")
	}

	lowercase := 0
	uppercase := 0
	digit := 0
	special := 0
	for _, char := range passwordStr {
		if unicode.IsUpper(char) {
			uppercase += 1
		} else if unicode.IsLower(char) {
			lowercase += 1
		} else if unicode.IsDigit(char) {
			digit += 1
		} else if unicode.IsSymbol(char) {
			special += 1
		}
	}

	if lowercase > passwordConfig.MinLowercaseLength {
		return nil, errors.New("password requires multiple (" + strconv.Itoa(passwordConfig.MinLowercaseLength) + ") lowercase letters")
	}

	if uppercase > passwordConfig.MinUppercaseLetter {
		return nil, errors.New("password requires multiple (" + strconv.Itoa(passwordConfig.MinUppercaseLetter) + ") uppercase letters")
	}

	if digit > passwordConfig.MinDigitLetter {
		return nil, errors.New("password requires multiple (" + strconv.Itoa(passwordConfig.MinDigitLetter) + ") digit letters")
	}

	if special > passwordConfig.MinSpecialLetter {
		return nil, errors.New("password requires multiple (" + strconv.Itoa(passwordConfig.MinSpecialLetter) + ") special letters")
	}

	passwordHash, err := passwordHasher.HashPassword(passwordStr, passwordSalt)
	if err != nil {
		return nil, err
	}

	metadata := map[string]any{
		"hash": passwordHash,
		"salt": passwordSalt,
		"type": providerContext.Tenant.PasswordType,
	}

	return metadata, nil
}

func (p passwordAuth) Validate(ctx context.Context, providerContext ProviderContext, data map[string]any) (bool, map[string]any, error) {
	success, err := p.Submit(ctx, providerContext, data)

	if err != nil {
		return false, nil, err
	}

	return success, providerContext.Credential.Metadata, nil
}

func (p passwordAuth) Begin(_ context.Context, _ ProviderContext) (map[string]any, error) {
	return make(map[string]any), nil
}

func (p passwordAuth) Submit(_ context.Context, providerContext ProviderContext, data map[string]any) (bool, error) {
	hash := providerContext.Credential.Metadata["hash"].(string)
	salt := providerContext.Credential.Metadata["salt"].(string)
	hashType := providerContext.Credential.Metadata["type"].(string)

	hasher, err := crypto.GetPasswordHasher(hashType)

	if err != nil {
		return false, err
	}

	password, exist := data["password"]

	if !exist {
		return false, errors.New("password is a required data field")
	}

	passwordStr, ok := password.(string)

	if !ok {
		return false, errors.New("password is not a string")
	}

	success, err := hasher.ComparePassword(passwordStr, hash, salt)
	if err != nil {
		return false, err
	}

	return success, nil
}
