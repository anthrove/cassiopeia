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
	"time"
)

type ProfilePage struct {
	UserID string `json:"user_id" gorm:"primaryKey;type:char(25)" validate:"required,max=100" maxLength:"100"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Fields []ProfilePageField `json:"fields"  gorm:"serializer:json"`

	User User `json:"-"`
}

type ProfilePageField struct {
	Identifier string `json:"identifier"`
	Value      any    `json:"value"`
}

type CreateProfilePage struct {
	Fields []ProfilePageField `json:"fields" validate:"required"`
}

type UpdateProfilePage struct {
	Fields []ProfilePageField `json:"fields" validate:"required"`
}
