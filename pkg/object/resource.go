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
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	CreatedAt time.Time `json:"createdAt" format:"date-time" example:"2025-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" format:"date-time" example:"2025-01-01T00:00:00Z"`

	ProviderID string `json:"provider_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	Tag        string `json:"tag" example:"example-tag"`
	MimeType   string `json:"mime_type" example:"image/png"`
	FilePath   string `json:"file_path" example:"/path/to/file.png"`
	FileSize   int64  `json:"file_size" example:"1024"`
	Format     string `json:"format" example:"png"`
	Url        string `json:"url" example:"https://domain.tld/files/file.png"`
	Hash       string `json:"hash" example:"d41d8cd98f00b204e9800998ecf8427e"`
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
	ProviderID string `json:"provider_id" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	Tag        string `json:"tag" example:"example-tag"`
	FileName   string `json:"file_name" example:"file.png"`
	FileSize   int64  `json:"file_size" example:"1024"`
	MimeType   string `json:"mime_type" example:"image/png"`
}
