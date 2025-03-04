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

type MessageTemplate struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" `
	TenantID string `json:"tenant_id"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	DisplayName  string `json:"display_name" validate:"required,max=100" maxLength:"100"`
	TemplateType int    `json:"template_type" validate:"required"`
	Template     string `json:"template" validate:"required"`
}

func (base *MessageTemplate) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateMessageTemplate struct {
	DisplayName  string `json:"display_name" validate:"required,max=100" maxLength:"100"`
	TemplateType int    `json:"template_type" validate:"required"`
	Template     string `json:"template" validate:"required"`
}

type UpdateMessageTemplate struct {
	DisplayName string `json:"display_name" validate:"required,max=100" maxLength:"100"`
	Template    string `json:"template" validate:"required"`
}
