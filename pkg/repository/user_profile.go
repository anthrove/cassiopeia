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
)

func CreateUserProfile(ctx context.Context, db *gorm.DB, tenantId string, userID string, createUserProfile object.CreateProfilePage) (object.ProfilePage, error) {
	page := object.ProfilePage{
		UserID: userID,
		Fields: createUserProfile.Fields,
	}

	err := db.WithContext(ctx).Model(&object.ProfilePage{}).Create(&page).Error

	return page, err
}

func UpdatePageProfile(ctx context.Context, db *gorm.DB, tenantID string, userID string, updateUserProfile object.UpdateProfilePage) error {
	profilePage := object.ProfilePage{
		UserID: userID,
		Fields: updateUserProfile.Fields,
	}

	err := db.WithContext(ctx).Model(&object.ProfilePage{}).Where("user_id = ?", userID).Updates(&profilePage).Error

	return err
}

func KillProfilePage(ctx context.Context, db *gorm.DB, tenantID string, userID string) error {
	return db.WithContext(ctx).Delete(&object.ProfilePage{}, "user_id = ?", userID).Error
}

func FindProfilePage(ctx context.Context, db *gorm.DB, tenantID string, userID string) (object.ProfilePage, error) {
	var profilePage object.ProfilePage
	err := db.WithContext(ctx).Take(&profilePage, "user_id = ?", userID).Error
	return profilePage, err
}
