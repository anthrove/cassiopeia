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

type Enforcer struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	ModelID   string  `json:"model_id"`
	Model     Model   `json:"-"`
	AdapterID string  `json:"adapter_id"`
	Adapter   Adapter `json:"-"`
}

func (base *Enforcer) BeforeCreate(db *gorm.DB) error {
	if base.ID == "" {
		id, err := gonanoid.New(25)
		if err != nil {
			return err
		}

		base.ID = id
	}

	return nil
}

type CreateEnforcer struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	ModelID   string `json:"model_id"`
	AdapterID string `json:"adapter_id"`
}

type UpdateEnforcer struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	ModelID   string `json:"model_id"`
	AdapterID string `json:"adapter_id"`
}
