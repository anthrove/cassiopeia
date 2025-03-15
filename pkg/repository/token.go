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
	"database/sql"
	"github.com/anthrove/identity/pkg/object"
	"gorm.io/gorm"
)

func CreateToken(ctx context.Context, db *gorm.DB, tenantId string, createToken object.CreateToken) (object.Token, error) {
	token := object.Token{
		TenantID:      tenantId,
		ApplicationID: createToken.ApplicationID,
		UserID: sql.NullString{
			String: createToken.UserID,
			Valid:  true,
		},
		AccessToken:  createToken.AccessToken,
		RefreshToken: createToken.RefreshToken,
		ExpiredAt:    createToken.ExpiredAt,
		Scope:        createToken.Scope,
	}

	err := db.WithContext(ctx).Model(&object.Token{}).Create(&token).Error

	return token, err
}

func KillToken(ctx context.Context, db *gorm.DB, tenantID string, tokenID string) error {
	return db.WithContext(ctx).Delete(&object.Token{}, "id = ? AND tenant_id = ?", tokenID, tenantID).Error
}

func KillTokens(ctx context.Context, db *gorm.DB, tenantID string, tokenIDs []string) error {
	return db.WithContext(ctx).Delete(&object.Token{}, "id in ? AND tenant_id = ?", tokenIDs, tenantID).Error
}

func FindToken(ctx context.Context, db *gorm.DB, tenantID string, tokenID string) (object.Token, error) {
	var token object.Token
	err := db.WithContext(ctx).Take(&token, "id = ? AND tenant_id = ?", tokenID, tenantID).Error
	return token, err
}

func FindTokenByRefresh(ctx context.Context, db *gorm.DB, tenantID string, refreshToken string) (object.Token, error) {
	var token object.Token
	err := db.WithContext(ctx).Take(&token, "refresh_token = ? AND tenant_id = ?", refreshToken, tenantID).Error
	return token, err
}

func FindTokens(ctx context.Context, db *gorm.DB, tenantID string, pagination object.Pagination) ([]object.Token, error) {
	var data []object.Token
	err := db.WithContext(ctx).Scopes(Pagination(pagination)).Where("tenant_id = ?", tenantID).Find(&data).Error
	return data, err
}

func FindUserTokens(ctx context.Context, db *gorm.DB, tenantID string, applicationID string, userID string) ([]object.Token, error) {
	var data []object.Token
	err := db.WithContext(ctx).Where("tenant_id = ? AND application_id = ? AND user_id = ?", tenantID, applicationID, userID).Find(&data).Error
	return data, err
}
