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

package repository

import (
	"context"
	"github.com/anthrove/identity/pkg/object"
	"gorm.io/gorm"
	"time"
)

func CreateAuthRequest(ctx context.Context, db *gorm.DB, tenantId string, createAuthRequest object.CreateAuthRequest) (object.AuthRequest, error) {
	authRequest := object.AuthRequest{
		TenantID:        tenantId,
		ApplicationID:   createAuthRequest.ApplicationID,
		UserID:          createAuthRequest.UserID,
		CreatedAt:       time.Now(),
		CallbackURI:     createAuthRequest.CallbackURI,
		TransferState:   createAuthRequest.TransferState,
		Prompt:          createAuthRequest.Prompt,
		LoginHint:       createAuthRequest.LoginHint,
		MaxAuthAge:      createAuthRequest.MaxAuthAge,
		Scopes:          createAuthRequest.Scopes,
		ResponseType:    createAuthRequest.ResponseType,
		ResponseMode:    createAuthRequest.ResponseMode,
		Nonce:           createAuthRequest.Nonce,
		CodeChallenge:   createAuthRequest.CodeChallenge,
		Authenticated:   false,
		AuthenticatedAt: time.Time{},
	}

	err := db.WithContext(ctx).Model(&object.AuthRequest{}).Create(&authRequest).Error

	return authRequest, err
}

func UpdateAuthRequest(ctx context.Context, db *gorm.DB, tenantID string, authRequestID string, updateAuthRequest object.UpdateAuthRequest) error {
	authRequest := object.AuthRequest{
		TenantID:        tenantID,
		UserID:          updateAuthRequest.UserID,
		Authenticated:   updateAuthRequest.Authenticated,
		AuthenticatedAt: time.Now(),
	}

	err := db.WithContext(ctx).Model(&object.AuthRequest{}).Where("id = ? AND tenant_id = ?", authRequestID, tenantID).Updates(&authRequest).Error

	return err
}

func KillAuthRequest(ctx context.Context, db *gorm.DB, tenantID string, authRequestID string) error {
	return db.WithContext(ctx).Delete(&object.AuthRequest{}, "id = ? AND tenant_id = ?", authRequestID, tenantID).Error
}

func FindAuthRequest(ctx context.Context, db *gorm.DB, tenantID string, authRequestID string) (object.AuthRequest, error) {
	var authRequest object.AuthRequest
	err := db.WithContext(ctx).Take(&authRequest, "id = ? AND tenant_id = ?", authRequestID, tenantID).Error
	return authRequest, err
}

func FindAuthRequests(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.AuthRequest, error) {
	var data []object.AuthRequest
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}
