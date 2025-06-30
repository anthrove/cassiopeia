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
	"github.com/anthrove/identity/pkg/object"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"testing"
	"time"
)

func TestToTpProviderFlow(t *testing.T) {
	provider, err := newTOTPProvider(object.Provider{
		ID:           "test",
		TenantID:     "test",
		DisplayName:  "Test",
		Category:     "mfa",
		ProviderType: "totp",
		Parameter:    nil,
	}, 30, otp.DigitsSix, otp.AlgorithmSHA1)

	if err != nil {
		t.Fatal(err)
	}

	data, err := provider.GenerateUserConfig("testuser")

	if err != nil {
		t.Fatal(err)
	}

	var totpProperties totpProperties

	err = json.Unmarshal(data.Properties, &totpProperties)

	if err != nil {
		t.Fatal(err)
	}

	code, err := totp.GenerateCode(totpProperties.Secret, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	success, err := provider.ValidateDatFlow(data.Properties, map[string]any{
		"otp": code,
	})

	if err != nil {
		t.Fatal(err)
	}

	if !success {
		t.Fatal("invalid otp")
	}
}
