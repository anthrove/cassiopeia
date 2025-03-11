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
	"os"
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

func GetStorageProvider(provider object.Provider) (Provider, error) {
	switch provider.ProviderType {
	case "local":
		return newLocalProvider(provider)
	case "s3":
		return newS3Provider(provider)
	}
	return nil, errors.New("unknown storage provider: " + provider.ProviderType)

}
