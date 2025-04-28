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

package provider

import (
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/provider/auth"
	"github.com/anthrove/identity/pkg/provider/email"
	"github.com/anthrove/identity/pkg/provider/storage"
)

func ConfigurationFields(category string, providerType string) []object.ProviderConfigurationField {
	switch category {
	case "auth":
		return auth.ConfigurationFields(providerType)
	case "email":
		return email.ConfigurationFields(providerType)
	case "storage":
		return storage.ConfigurationFields(providerType)
	}

	return nil
}

func Types(category string) []string {
	switch category {
	case "auth":
		return auth.GetAuthTypes()
	case "email":
		return email.GetEMailTypes()
	case "storage":
		return storage.GetStorageTypes()
	}

	return nil
}
