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
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             string `json:"id" gorm:"primaryKey;type:char(25)"`
	OrganisationID string `json:"organisation_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	Username      string `json:"username;type:varchar(100)"`
	DisplayName   string `json:"display_name;type:varchar(100)"`
	Email         string `json:"email;type:varchar(100);index"`
	EmailVerified bool   `json:"email_verified"`
	PasswordHash  string `json:"password_hash;type:varchar(150)"`
	PasswordSalt  string `json:"password_salt;type:varchar(100)"`
	PasswordType  string `json:"password_type;type:varchar(100)"`

	Groups []Group `json:"groups" gorm:"many2many:user_groups;"`
}
