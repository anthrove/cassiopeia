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
	ID       string `json:"id" gorm:"primaryKey;type:char(25)"`
	TenantID string `json:"tenant_id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ParentGroupID *string `json:"parent_group_id"`
	ParentGroup   *Group  `json:"-"`

	DisplayName string `json:"displayName" gorm:"type:varchar(100)"`
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

type CreateGroup struct {
	DisplayName   string  `json:"display_name" validate:"required,max=100"`
	ParentGroupID *string `json:"parent_group_id" validate:"omitnil,len=25"`
}

type UpdateGroup struct {
	DisplayName   string  `json:"display_name" validate:"required,max=100"`
	ParentGroupID *string `json:"parent_group_id" validate:"omitnil,len=25"`
}
