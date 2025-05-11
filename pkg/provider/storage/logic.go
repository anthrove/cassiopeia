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
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/qor/oss"
	"io"
	"maps"
	"os"
	"slices"
)

type Provider interface {
	GetConfigurationFields() []object.ProviderConfigurationField
	ValidateConfigurationFields() error
	Get(path string) (*os.File, error)
	GetStream(path string) (io.ReadCloser, error)
	Put(path string, reader io.Reader) (*oss.Object, error)
	Delete(path string) error
	List(path string) ([]*oss.Object, error)
	GetEndpoint() string
	GetURL(path string) (string, error)
}

var providerMap = map[string]func(provider object.Provider) (Provider, error){
	"local": newLocalProvider,
	"s3":    newS3Provider,
}

func GetStorageProvider(provider object.Provider) (Provider, error) {
	newFunc, exists := providerMap[provider.ProviderType]

	if !exists {
		return nil, errors.New("unknown storage provider: " + provider.ProviderType)
	}

	return newFunc(provider)
}

func ConfigurationFields(providerType string) []object.ProviderConfigurationField {
	switch providerType {
	case "local":
		return localProvider{}.GetConfigurationFields()
	case "s3":
		return s3Provider{}.GetConfigurationFields()

	}

	return nil
}

func GetStorageTypes() []string {
	return slices.Collect(maps.Keys(providerMap))
}
