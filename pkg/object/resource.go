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

type Resource struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" `
	TenantID string `json:"tenant_id" gorm:"type:char(25)"`

	CreatedAt time.Time `json:"createdAt" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time"`

	ProviderID string `json:"provider_id" gorm:"type:char(25)"`
	Tag        string `json:"tag"`
	MimeType   string `json:"mime_type"`
	FilePath   string `json:"file_path"`
	FileSize   int64  `json:"file_size"`
	Format     string `json:"format"`
	Url        string `json:"url"`
	Hash       string `json:"hash"`
}

func (base *Resource) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateResource struct {
	ProviderID string `json:"provider_id"`
	Tag        string `json:"tag"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `json:"mime_type"`
}
