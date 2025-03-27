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
	"database/sql"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"gorm.io/gorm"
	"time"
)

type AuthRequest struct {
	ID            string         `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID      string         `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	ApplicationID string         `json:"application_id" maxLength:"25" minLength:"25" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	UserID        sql.NullString `json:"user_id" maxLength:"25" gorm:"type:char(25)" `

	CreatedAt time.Time `json:"created_at" format:"date-time"`

	CallbackURI   string             `json:"callback_uri"`
	TransferState string             `json:"transfer_state"`
	Prompt        []string           `json:"prompt" gorm:"type:text[]; serializer:json"`
	LoginHint     string             `json:"login_hint"`
	MaxAuthAge    *time.Duration     `json:"max_auth_age"`
	Scopes        []string           `json:"scopes" gorm:"type:text[]; serializer:json"`
	ResponseType  oidc.ResponseType  `json:"response_type"`
	ResponseMode  oidc.ResponseMode  `json:"response_mode"`
	Nonce         string             `json:"nonce"`
	CodeChallenge *OIDCCodeChallenge `json:"code_challenge" gorm:"type:text; serializer:json"`

	Authenticated   bool      `json:"authenticated" format:"date-time"`
	AuthenticatedAt time.Time `json:"authenticated_at" format:"date-time"`
}

func (a AuthRequest) GetID() string {
	return a.ID
}

func (a AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (a AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if a.Authenticated {
		return []string{"pwd"}
	}
	return nil
}

func (a AuthRequest) GetAudience() []string {
	return []string{a.ApplicationID} // this example will always just use the client_id as audience
}

func (a AuthRequest) GetAuthTime() time.Time {
	return a.AuthenticatedAt
}

func (a AuthRequest) GetClientID() string {
	return a.ApplicationID
}

func (a AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.CodeChallenge)
}

func (a AuthRequest) GetNonce() string {
	return a.Nonce
}

func (a AuthRequest) GetRedirectURI() string {
	return a.CallbackURI
}

func (a AuthRequest) GetResponseType() oidc.ResponseType {
	return a.ResponseType
}

func (a AuthRequest) GetResponseMode() oidc.ResponseMode {
	return a.ResponseMode
}

func (a AuthRequest) GetScopes() []string {
	return a.Scopes
}

func (a AuthRequest) GetState() string {
	return a.TransferState
}

func (a AuthRequest) GetSubject() string {
	return a.UserID.String
}

func (a AuthRequest) Done() bool {
	return a.Authenticated
}

// BeforeCreate is a GORM hook that is called before a new group record is inserted into the database.
// It generates a unique ID for the group if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (a *AuthRequest) BeforeCreate(db *gorm.DB) error {
	if a.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		a.ID = id
	}

	return nil
}

type CreateAuthRequest struct {
	ApplicationID string         `json:"application_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	UserID        sql.NullString `json:"user_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CallbackURI   string             `json:"callback_uri"`
	TransferState string             `json:"transfer_state"`
	Prompt        []string           `json:"prompt" gorm:"type:text[]; serializer:json"`
	LoginHint     string             `json:"login_hint"`
	MaxAuthAge    *time.Duration     `json:"max_auth_age"`
	Scopes        []string           `json:"scopes" gorm:"type:text[]; serializer:json"`
	ResponseType  oidc.ResponseType  `json:"response_type"`
	ResponseMode  oidc.ResponseMode  `json:"response_mode"`
	Nonce         string             `json:"nonce"`
	CodeChallenge *OIDCCodeChallenge `json:"code_challenge" gorm:"type:text; serializer:json"`
}

type UpdateAuthRequest struct {
	UserID sql.NullString `json:"user_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	Authenticated   bool      `json:"authenticated" format:"date-time"`
	AuthenticatedAt time.Time `json:"authenticated_at" format:"date-time"`
}

// Copyright https://github.com/zitadel/oidc/blob/eb2f912c5e5a783e6fb682d5eeea3a13b1ad12c7/example/server/storage/oidc.go#L145
// =============================================

type AccessToken struct {
	ID             string
	ApplicationID  string
	Subject        string
	RefreshTokenID string
	Audience       []string
	Expiration     time.Time
	Scopes         []string
}

type RefreshToken struct {
	ID            string
	Token         string
	AuthTime      time.Time
	AMR           []string
	Audience      []string
	UserID        string
	ApplicationID string
	Expiration    time.Time
	Scopes        []string
	AccessToken   string // Token.ID
}

type OIDCCodeChallenge struct {
	Challenge string
	Method    string
}

func CodeChallengeToOIDC(challenge *OIDCCodeChallenge) *oidc.CodeChallenge {
	if challenge == nil {
		return nil
	}
	challengeMethod := oidc.CodeChallengeMethodPlain
	if challenge.Method == "S256" {
		challengeMethod = oidc.CodeChallengeMethodS256
	}
	return &oidc.CodeChallenge{
		Challenge: challenge.Challenge,
		Method:    challengeMethod,
	}
}

// RefreshTokenRequestFromBusiness will simply wrap the storage RefreshToken to implement the op.RefreshTokenRequest interface
func RefreshTokenRequestFromBusiness(token *RefreshToken) op.RefreshTokenRequest {
	return &RefreshTokenRequest{token}
}

type RefreshTokenRequest struct {
	*RefreshToken
}

func (r *RefreshTokenRequest) GetAMR() []string {
	return r.AMR
}

func (r *RefreshTokenRequest) GetAudience() []string {
	return r.Audience
}

func (r *RefreshTokenRequest) GetAuthTime() time.Time {
	return r.AuthTime
}

func (r *RefreshTokenRequest) GetClientID() string {
	return r.ApplicationID
}

func (r *RefreshTokenRequest) GetScopes() []string {
	return r.Scopes
}

func (r *RefreshTokenRequest) GetSubject() string {
	return r.UserID
}

func (r *RefreshTokenRequest) SetCurrentScopes(scopes []string) {
	r.Scopes = scopes
}
