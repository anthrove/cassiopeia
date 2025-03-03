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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/go-playground/validator/v10"
)

type smtpConfiguration struct {
	Host     string `json:"host" validate:"required,max=100"`
	Port     int    `json:"port" validate:"required"`
	SSL      bool   `json:"ssl"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	From     string `json:"from" validate:"required,max=100,email"`
	FromName string `json:"from_name" validate:"required,max=100"`
}

type smtpProvider struct {
	provider object.Provider
}

func newSMTPProvider(provider object.Provider) (Provider, error) {
	return smtpProvider{provider: provider}, nil
}

func (s smtpProvider) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{
		{
			FieldKey:  "host",
			FieldType: "text",
		},
		{
			FieldKey:  "port",
			FieldType: "int",
		},
		{
			FieldKey:  "ssl",
			FieldType: "bool",
		},
		{
			FieldKey:  "username",
			FieldType: "text",
		},
		{
			FieldKey:  "password",
			FieldType: "secret",
		},
		{
			FieldKey:  "from",
			FieldType: "text",
		},
		{
			FieldKey:  "from_name",
			FieldType: "text",
		},
	}
}

func (s smtpProvider) ValidateConfigurationFields() error {
	smtpConfig := smtpConfiguration{}

	err := json.Unmarshal(s.provider.Parameter, &smtpConfig)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(smtpConfig)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return nil
}

func (s smtpProvider) SendMail(fromAddress string, fromName string, toAddress, subject string, body string) error {
	//TODO implement me
	panic("implement me")
}
