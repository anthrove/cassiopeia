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

// Group represents a group entity in the system.
// It contains information about the group such as its ID, tenant ID, timestamps, parent group, display name, and status.
type Group struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	ParentGroupID *string `json:"parent_group_id" maxLength:"25" minLength:"25" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	ParentGroup   *Group  `json:"-"`

	DisplayName string `json:"display_name" gorm:"type:varchar(100)" maxLength:"100" example:"Tenant Title"`

	Users []User `json:"-" gorm:"many2many:user_groups;"`
}

// BeforeCreate is a GORM hook that is called before a new group record is inserted into the database.
// It generates a unique ID for the group if it is not already set.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue generating the unique ID.
func (base *Group) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

// CreateGroup represents the data required to create a new group.
// It includes the display name and an optional parent group ID, both of which have validation constraints.
type CreateGroup struct {
	DisplayName   string  `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Tenant Title"`
	ParentGroupID *string `json:"parent_group_id" validate:"omitnil,len=25" maxLength:"25" minLength:"25"`
}

// UpdateGroup represents the data required to update an existing group.
// It includes the display name and an optional parent group ID, both of which have validation constraints.
type UpdateGroup struct {
	DisplayName   string  `json:"display_name" validate:"required,max=100" maxLength:"100" example:"Tenant Title"`
	ParentGroupID *string `json:"parent_group_id" validate:"omitnil,len=25" maxLength:"25" minLength:"25"`
}
