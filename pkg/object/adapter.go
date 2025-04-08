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
)

type Adapter struct {
	ID       string `json:"id" gorm:"primaryKey;type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`
	TenantID string `json:"tenant_id" gorm:"type:char(25)" example:"BsOOg4igppKxYwhAQQrD3GCRZ"`

	Name       string `json:"name"`
	TableName  string `json:"table_name"`
	ExternalDB bool   `json:"external_db"`

	Driver string `json:"driver"`

	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}

func (base *Adapter) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateAdapter struct {
	Name       string `json:"name" validate:"required"`
	TableName  string `json:"table_name" validate:"required"`
	ExternalDB bool   `json:"external_db" validate:"required"`

	Driver string `json:"driver"  validate:"required_if=ExternalDB true"  example:"mysql"`

	Host         string `json:"host" validate:"required_if=ExternalDB true" example:"127.0.0.1"`
	Port         string `json:"port" validate:"required_if=ExternalDB true"  example:"3306"`
	Username     string `json:"username"  validate:"required_if=ExternalDB true" example:"root"`
	Password     string `json:"password"  validate:"required_if=ExternalDB true" example:"password"`
	DatabaseName string `json:"database_name"  validate:"required_if=ExternalDB true"  example:"test"`
}

type UpdateAdapter struct {
	Name       string `json:"name" validate:"required"`
	TableName  string `json:"table_name" validate:"required"`
	ExternalDB bool   `json:"external_db" validate:"required"`

	Driver string `json:"driver"  validate:"required_if=ExternalDB true"  example:"mysql"`

	Host         string `json:"host" validate:"required_if=ExternalDB true" example:"127.0.0.1"`
	Port         string `json:"port" validate:"required_if=ExternalDB true"  example:"3306"`
	Username     string `json:"username"  validate:"required_if=ExternalDB true" example:"root"`
	Password     string `json:"password"  validate:"required_if=ExternalDB true" example:"password"`
	DatabaseName string `json:"database_name"  validate:"required_if=ExternalDB true"  example:"test"`
}
