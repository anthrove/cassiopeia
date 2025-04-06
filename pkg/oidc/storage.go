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

package oidc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/anthrove/identity/pkg/logic"
	"github.com/anthrove/identity/pkg/object"
	"github.com/anthrove/identity/pkg/util"
	"github.com/go-jose/go-jose/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"gorm.io/gorm"
	"math"
	"slices"
	"strings"
	"sync"
	"time"
)

type fullStorage interface {
	op.Storage
	op.TokenExchangeStorage
}

type storage struct {
	service logic.IdentityService
	tenant  object.Tenant

	//TODO put this data into db
	lock         sync.Mutex
	requestCodes map[string]string
}

func NewStorage(ctx context.Context, is logic.IdentityService, tenantID string) (fullStorage, error) {
	tenant, err := is.FindTenant(ctx, tenantID)

	if err != nil {
		return nil, err
	}

	return &storage{
		service:      is,
		tenant:       tenant,
		requestCodes: make(map[string]string),
	}, nil
}

// Documentation to this function are from here: https://github.com/zitadel/oidc/blob/main/example/server/storage/storage.go
// ============================================

// CreateAuthRequest implements the op.Storage interface
// it will be called after parsing and validation of the authentication request
func (s *storage) CreateAuthRequest(ctx context.Context, request *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(request.Prompt) == 1 && request.Prompt[0] == "none" {
		// With prompt=none, there is no way for the user to log in
		// so return error right away.
		return nil, oidc.ErrLoginRequired()
	}

	id, err := gonanoid.New(25)
	if err != nil {
		return nil, err
	}

	req := authRequestToInternal(id, request, userID)

	req, err = s.service.CreateAuthRequest(ctx, s.tenant.ID, object.CreateAuthRequest{
		ApplicationID: request.ClientID,
		CallbackURI:   request.RedirectURI,
		TransferState: request.State,
		Prompt:        PromptToInternal(request.Prompt),
		LoginHint:     request.LoginHint,
		MaxAuthAge:    MaxAgeToInternal(request.MaxAge),
		UserID:        sql.NullString{String: userID, Valid: true},
		Scopes:        request.Scopes,
		ResponseType:  request.ResponseType,
		ResponseMode:  request.ResponseMode,
		Nonce:         request.Nonce,
		CodeChallenge: &object.OIDCCodeChallenge{
			Challenge: request.CodeChallenge,
			Method:    string(request.CodeChallengeMethod),
		},
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}

// AuthRequestByID implements the op.Storage interface
// it will be called after the Login UI redirects back to the OIDC endpoint
func (s *storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	return s.service.FindAuthRequest(ctx, s.tenant.ID, id)
}

// AuthRequestByCode implements the op.Storage interface
// it will be called after parsing and validation of the token request (in an authorization code flow)
func (s *storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	requestID, ok := s.requestCodes[code]

	if !ok {
		return nil, errors.New("auth request not found")
	}

	return s.service.FindAuthRequest(ctx, s.tenant.ID, requestID)
}

// SaveAuthCode implements the op.Storage interface
// it will be called after the authentication has been successful and before redirecting the user agent to the redirect_uri
// (in an authorization code flow)
func (s *storage) SaveAuthCode(_ context.Context, id string, code string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.requestCodes[code] = id

	return nil
}

// DeleteAuthRequest implements the op.Storage interface
// it will be called after creating the token response (id and access tokens) for a valid
// - authentication request (in an implicit flow)
// - token request (in an authorization code flow)
func (s *storage) DeleteAuthRequest(ctx context.Context, id string) error {
	return s.service.KillAuthRequest(ctx, s.tenant.ID, id)
}

// CreateAccessToken implements the op.Storage interface
// it will be called for all requests able to return an access token (Authorization Code Flow, Implicit Flow, JWT Profile, ...)
func (s *storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	token, err := s.service.CreateToken(ctx, s.tenant.ID, object.CreateToken{
		ApplicationID: "",
		UserID:        request.GetSubject(),
		Scope:         strings.Join(request.GetScopes(), " "),
		Audience:      strings.Join(request.GetAudience(), " "),
		ExpiredAt:     time.Now().Add(24 * time.Hour),
	})

	if err != nil {
		return "", time.Time{}, err
	}

	return token.ID, token.ExpiredAt, nil
}

// CreateAccessAndRefreshTokens implements the op.Storage interface
// it will be called for all requests able to return an access and refresh token (Authorization Code Flow, Refresh Token Request)
func (s *storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {
	//TODO implement me
	panic("implement me")
}

// TokenRequestByRefreshToken implements the op.Storage interface
// it will be called after parsing and validation of the refresh token request
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	//TODO implement me
	panic("implement me")
}

// TerminateSession implements the op.Storage interface
// it will be called after the user signed out, therefore the access and refresh token of the user of this client must be removed
func (s *storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	userTokens, err := s.service.FindUserTokens(ctx, s.tenant.ID, clientID, userID)

	if err != nil {
		return err
	}

	tokenIDs := make([]string, len(userTokens))
	for i, token := range userTokens {
		tokenIDs[i] = token.ID
	}

	err = s.service.KillTokens(ctx, s.tenant.ID, tokenIDs)

	return err
}

// RevokeToken implements the op.Storage interface
// it will be called after parsing and validation of the token revocation request
func (s *storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	//TODO What if token is a refresh-token instead of token id
	token, err := s.service.FindToken(ctx, s.tenant.ID, tokenOrTokenID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Ignore if not exists!
			return nil
		}

		return oidc.ErrServerError().WithParent(err)
	}

	if token.ApplicationID != clientID {
		return oidc.ErrInvalidClient().WithDescription("wrong client for this token")
	}

	err = s.service.KillToken(ctx, s.tenant.ID, tokenOrTokenID)

	if err != nil {
		return oidc.ErrServerError().WithParent(err)
	}

	return nil
}

// GetRefreshTokenInfo looks up a refresh token and returns the token id and user id.
// If given something that is not a refresh token, it must return error.
func (s *storage) GetRefreshTokenInfo(ctx context.Context, clientID string, refreshToken string) (userID string, tokenID string, err error) {
	token, err := s.service.FindTokenByRefresh(ctx, s.tenant.ID, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", op.ErrInvalidRefreshToken
		}
	}

	if token.ApplicationID != clientID {
		return "", "", oidc.ErrInvalidClient().WithDescription("wrong client for this refresh token")
	}

	return token.UserID.String, token.ID, nil
}

// SigningKey implements the op.Storage interface
// it will be called when creating the OpenID Provider
func (s *storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	certificate, err := s.service.FindCertificate(ctx, s.tenant.ID, *s.tenant.SigningCertificateID)

	if err != nil {
		return nil, err
	}

	return certificate.ToSigningCert(), nil
}

// SignatureAlgorithms implements the op.Storage interface
// it will be called to get the sign
func (s *storage) SignatureAlgorithms(_ context.Context) ([]jose.SignatureAlgorithm, error) {
	return []jose.SignatureAlgorithm{jose.RS256, jose.RS384, jose.RS512, jose.ES256, jose.ES384, jose.ES512}, nil
}

// KeySet implements the op.Storage interface
// it will be called to get the current (public) keys, among others for the keys_endpoint or for validating access_tokens on the userinfo_endpoint, ...
func (s *storage) KeySet(ctx context.Context) ([]op.Key, error) {
	certs, err := s.service.FindCertificates(ctx, s.tenant.ID, object.Pagination{
		Limit: math.MaxInt,
		Page:  0,
	})

	if err != nil {
		return nil, err
	}

	keys := make([]op.Key, 0, len(certs))
	for _, cert := range certs {
		keys = append(keys, &cert)
	}

	return keys, err
}

// GetClientByClientID implements the op.Storage interface
// it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	application, err := s.service.FindApplication(ctx, s.tenant.ID, clientID)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

// AuthorizeClientIDSecret implements the op.Storage interface
// it will be called for validating the client_id, client_secret on token or introspection requests
func (s *storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	application, err := s.service.FindApplication(ctx, s.tenant.ID, clientID)

	if err != nil {
		return err
	}

	if application.ClientSecret != clientSecret {
		return errors.New("authorization client secret does not match")
	}

	return nil
}

// SetUserinfoFromScopes implements the op.Storage interface.
// Provide an empty implementation and use SetUserinfoFromRequest instead.
func (s *storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	return nil
}

// SetUserinfoFromToken implements the op.Storage interface
// it will be called for the userinfo endpoint, so we read the token and pass the information from that to the private function
func (s *storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
	token, err := s.service.FindToken(ctx, s.tenant.ID, tokenID)

	if err != nil {
		return err
	}

	if token.ExpiredAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return s.setUserinfo(ctx, userinfo, token.UserID.String, token.ApplicationID, strings.Split(token.Scope, " "))
}

// SetIntrospectionFromToken implements the op.Storage interface
// it will be called for the introspection endpoint, so we read the token and pass the information from that to the private function
func (s *storage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	token, err := s.service.FindToken(ctx, s.tenant.ID, tokenID)

	if err != nil {
		return err
	}

	if token.ExpiredAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	if token.ApplicationID != clientID {
		return errors.New("token application does not match")
	}

	userInfo := new(oidc.UserInfo)
	err = s.setUserinfo(ctx, userInfo, subject, clientID, strings.Split(token.Scope, " "))
	if err != nil {
		return err
	}

	introspection.SetUserInfo(userInfo)
	//...and also the requested scopes...
	introspection.Scope = strings.Split(token.Scope, " ")
	//...and the client the token was issued to
	introspection.ClientID = token.ApplicationID
	return nil
}

// setUserinfo sets the info based on the user, scopes and if necessary the clientID
func (s *storage) setUserinfo(ctx context.Context, userInfo *oidc.UserInfo, userID, clientID string, scopes []string) (err error) {
	user, err := s.service.FindUser(ctx, s.tenant.ID, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.ID
		case oidc.ScopeEmail:
			userInfo.Email = user.Email
			userInfo.EmailVerified = oidc.Bool(user.EmailVerified)
		case oidc.ScopeProfile:
			userInfo.PreferredUsername = user.Username
			userInfo.Name = user.DisplayName
		case oidc.ScopePhone:
			//TODO setup phone scope
		}
	}
	return nil
}

// GetPrivateClaimsFromScopes implements the op.Storage interface
// it will be called for the creation of a JWT access token to assert claims for custom scopes
func (s *storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (map[string]any, error) {
	user, err := s.service.FindUser(ctx, s.tenant.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	claims := make(map[string]any)
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			claims["sub"] = user.ID
		case oidc.ScopeEmail:
			claims["email"] = user.Email
			claims["email_verified"] = user.EmailVerified
		case oidc.ScopeProfile:
			claims["preferred_username"] = user.Username
			claims["name"] = user.DisplayName
		case oidc.ScopePhone:
			//TODO setup phone scope
		}
	}
	return claims, nil
}

// GetKeyByIDAndClientID implements the op.Storage interface
// it will be called to validate the signatures of a JWT (JWT Profile Grant and Authentication)
func (s *storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	token, err := s.service.FindToken(ctx, s.tenant.ID, keyID)

	if err != nil {
		return nil, err
	}

	if token.ApplicationID != clientID {
		return nil, errors.New("clientID not found")
	}

	certificate, err := s.service.FindCertificate(ctx, s.tenant.ID, *s.tenant.SigningCertificateID)

	if err != nil {
		return nil, err
	}

	publicKey, err := util.BytesToPublicKey([]byte(certificate.Certificate))

	if err != nil {
		return nil, err
	}

	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}

var allowedScopes = []string{
	oidc.ScopeOpenID,
	oidc.ScopeEmail,
	oidc.ScopeProfile,
	oidc.ScopeOfflineAccess,
}

// ValidateJWTProfileScopes implements the op.Storage interface
// it will be called to validate the scopes of a JWT Profile Authorization Grant request
func (s *storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	allowed := make([]string, 0)
	for _, scope := range scopes {
		if slices.Contains(allowedScopes, scope) {
			allowed = append(allowed, scope)
		}
	}
	return allowed, nil
}

// Health implements the op.Storage interface
func (s *storage) Health(_ context.Context) error {
	return nil
}

func (s *storage) ValidateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	//TODO for what is this function?!?
	return nil
}

func (s *storage) CreateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	//TODO for what is this function?!?
	return nil
}

func (s *storage) GetPrivateClaimsFromTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) (claims map[string]any, err error) {
	return s.GetPrivateClaimsFromScopes(ctx, request.GetSubject(), request.GetClientID(), request.GetScopes())
}

func (s *storage) SetUserinfoFromTokenExchangeRequest(ctx context.Context, userinfo *oidc.UserInfo, request op.TokenExchangeRequest) error {
	return s.SetUserinfoFromScopes(ctx, userinfo, request.GetSubject(), request.GetClientID(), request.GetScopes())
}

// ============================================
