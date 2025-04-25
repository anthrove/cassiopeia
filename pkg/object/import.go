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

package object

import "encoding/json"

type ImportTenant struct {
	ID                   string  `json:"id"`
	DisplayName          string  `json:"display_name"`
	PasswordType         string  `json:"password_type"`
	SigningCertificateID *string `json:"signing_certificate_id"`
}

type ImportProvider struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"display_name"`
	Category    string          `json:"category"`
	Type        string          `json:"type"`
	Parameter   json.RawMessage `json:"parameter"`
}

type ImportGroup struct {
	ID            string  `json:"id"`
	DisplayName   string  `json:"display_name"`
	ParentGroupID *string `json:"parent_group_id"`
}

type ImportApplication struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`

	ClientSecret string `json:"client_secret"`

	Logo string `json:"logo"`

	SignInURL string `json:"sign_in_url"`
	SignUpURL string `json:"sign_up_url"`
	ForgetURL string `json:"forget_url"`
	TermsURL  string `json:"terms_url"`

	RedirectURLs []string `json:"redirect_urls"`
}

type ImportUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
