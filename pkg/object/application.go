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

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"gorm.io/gorm"
	"time"
)

type Application struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	ClientSecret string `json:"client_secret"`

	DisplayName string `json:"display_name" gorm:"type:varchar(100)" example:"Frontend Application"`
	Logo        string `json:"logo" gorm:"type:varchar(255)" example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url" gorm:"type:varchar(255)"`
	SignUpURL string `json:"sign_up_url" gorm:"type:varchar(255)"`
	ForgetURL string `json:"forget_url" gorm:"type:varchar(255)"`
	TermsURL  string `json:"terms_url" gorm:"type:varchar(255)"`

	RedirectURLs []string `json:"redirect_urls" gorm:"serializer:json"`

	Tokens       []Token    `json:"-" swaggerignore:"true"`
	AuthProvider []Provider `json:"auth_provider" gorm:"many2many:auth_application_provider;"`
}

// BeforeCreate is a GORM hook that is called before a new group record is inserted into the database.
// It generates a unique ID for the group if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *Application) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	if base.ClientSecret == "" {
		clientSecret, err := gonanoid.New(50)
		if err != nil {
			return err
		}

		base.ClientSecret = clientSecret
	}

	return nil
}

// Documentation to this function are from here: https://github.com/zitadel/oidc/blob/main/example/server/storage/client.go
// ============================================

// GetID must return the client_id
func (base *Application) GetID() string {
	return base.ID
}

// RedirectURIs must return the registered redirect_uris for Code and Implicit Flow
func (base *Application) RedirectURIs() []string {
	return base.RedirectURLs
}

// PostLogoutRedirectURIs must return the registered post_logout_redirect_uris for sign-outs
func (base *Application) PostLogoutRedirectURIs() []string {
	return []string{base.SignInURL}
}

// ApplicationType must return the type of the client (app, native, user agent)
func (base *Application) ApplicationType() op.ApplicationType {
	return op.ApplicationTypeWeb // TODO check if this is the only variant we want?
}

// AuthMethod must return the authentication method (client_secret_basic, client_secret_post, none, private_key_jwt)
func (base *Application) AuthMethod() oidc.AuthMethod {
	// https://connect2id.com/products/server/docs/guides/oauth-client-authentication#shared-secret-based
	return oidc.AuthMethodPost
}

// ResponseTypes must return all allowed response types (code, id_token token, id_token)
// these must match with the allowed grant types
func (base *Application) ResponseTypes() []oidc.ResponseType {
	return []oidc.ResponseType{oidc.ResponseTypeCode, oidc.ResponseTypeIDTokenOnly, oidc.ResponseTypeIDToken}
}

// GrantTypes must return all allowed grant types (authorization_code, refresh_token, urn:ietf:params:oauth:grant-type:jwt-bearer)
func (base *Application) GrantTypes() []oidc.GrantType {
	return oidc.AllGrantTypes //TODO make it configurable
}

// LoginURL will be called to redirect the user (agent) to the login UI
// you could implement some logic here to redirect the users to different login UIs depending on the client
func (base *Application) LoginURL(requestID string) string {
	return base.SignInURL + "?request_id=" + requestID
}

// AccessTokenType must return the type of access token the client uses (Bearer (opaque) or JWT)
func (base *Application) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeJWT
}

// IDTokenLifetime must return the lifetime of the client's id_tokens
func (base *Application) IDTokenLifetime() time.Duration {
	return time.Hour * 1
}

// DevMode enables the use of non-compliant configs such as redirect_uris (e.g. http schema for user agent client)
func (base *Application) DevMode() bool {
	return false
}

// RestrictAdditionalIdTokenScopes allows specifying which custom scopes shall be asserted into the id_token
func (base *Application) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// RestrictAdditionalAccessTokenScopes allows specifying which custom scopes shall be asserted into the JWT access_token
func (base *Application) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// IsScopeAllowed enables Client specific custom scopes validation
// in this example we allow the CustomScope for all clients
func (base *Application) IsScopeAllowed(scope string) bool {
	return true // TODO Scope validation?
}

// IDTokenUserinfoClaimsAssertion allows specifying if claims of scope profile, email, phone and address are asserted into the id_token
// even if an access token if issued which violates the OIDC Core spec
// (5.4. Requesting Claims using Scope Values: https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims)
// some clients though require that e.g. email is always in the id_token when requested even if an access_token is issued
func (base *Application) IDTokenUserinfoClaimsAssertion() bool {
	return true
}

// ClockSkew enables clients to instruct the OP to apply a clock skew on the various times and expirations
// (subtract from issued_at, add to expiration, ...)
func (base *Application) ClockSkew() time.Duration {
	return 0
}

// ============================================

type CreateApplication struct {
	DisplayName string `json:"display_name" validate:"required,max=100" example:"Frontend Application"`
	Logo        string `json:"logo"  example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url"`
	SignUpURL string `json:"sign_up_url"`
	ForgetURL string `json:"forget_url"`
	TermsURL  string `json:"terms_url"`

	RedirectURLs []string `json:"redirect_urls"`
}

type UpdateApplication struct {
	DisplayName string `json:"display_name" validate:"required,max=100" example:"Frontend Application"`
	Logo        string `json:"logo"  example:"https://domain.tld/files/logo.png"`

	SignInURL string `json:"sign_in_url"`
	SignUpURL string `json:"sign_up_url"`
	ForgetURL string `json:"forget_url"`
	TermsURL  string `json:"terms_url"`

	RedirectURLs []string `json:"redirect_urls"`
}
