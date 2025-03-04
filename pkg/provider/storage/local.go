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

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"github.com/go-playground/validator/v10"
)

type localConfiguration struct {
	basePath string
}

type localProvider struct {
	provider object.Provider
}

func newLocalProvider(provider object.Provider) (Provider, error) {
	return localProvider{provider: provider}, nil
}

func (l localProvider) GetConfigurationFields() []object.ProviderConfigurationField {
	return []object.ProviderConfigurationField{
		{
			FieldKey:  "basePath",
			FieldType: "text",
		},
	}
}

func (l localProvider) ValidateConfigurationFields() error {
	localConfig := localConfiguration{}

	err := json.Unmarshal(l.provider.Parameter, &localConfig)
	if err != nil {
		return err
	}

	// use a single instance of Validate, it caches struct info
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(localConfig)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return errors.Join(fmt.Errorf("problem while validating create tenant data"), validateErrs)
		}
	}

	return nil
}

func (l localProvider) GetFile() error {
	//TODO implement me
	panic("implement me")
}

func (l localProvider) SaveFile() error {
	//TODO implement me
	panic("implement me")
}
