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
	"gorm.io/gorm"
	"time"
)

// User represents a user entity in the system.
// It contains information about the user such as their ID, organisation ID, timestamps, username, display name, email, password details, and associated groups.
type User struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)"`
	TenantID string `json:"tenant_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`

	Username               string `json:"username"  gorm:"type:varchar(100);index"`
	DisplayName            string `json:"display_name"  gorm:"type:varchar(100)"`
	Email                  string `json:"email" gorm:"type:varchar(100);index"`
	EmailVerified          bool   `json:"email_verified"`
	EmailVerificationToken string `json:"-" gorm:"type:char(6)"`

	Groups []Group `json:"groups" gorm:"many2many:user_groups;"`
}

// BeforeCreate is a GORM hook that is called before a new user record is inserted into the database.
// It generates a unique ID for the user if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *User) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

// CreateUser represents the data required to create a new user.
// It includes the username, display name, email, and password, all of which are required and have validation constraints.
type CreateUser struct {
	Username    string `json:"username" validate:"required,max=100"`
	DisplayName string `json:"display_name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,max=100,email"`
	Password    string `json:"password" validate:"required,max=100"`
}

// UpdateUser represents the data required to update an existing user's display name.
// It includes the display name, which is required and has a maximum length of 100 characters.
type UpdateUser struct {
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100"`
}

// UpdateUserPassword represents the data required to update an existing user's password.
// It includes the password, which is required and has a maximum length of 100 characters.
type UpdateUserPassword struct {
	Password string `json:"password" validate:"required,max=100"`
}
