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

package logic

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

// IdentityService provides methods to interact with the identity database.
type IdentityService struct {
	db *gorm.DB
}

// NewIdentityService initializes a new IdentityService with the given database connection.
// It also sets up the validator instance used for validating structs.
//
// Parameters:
//   - db: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An initialized IdentityService instance.
func NewIdentityService(db *gorm.DB) IdentityService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return IdentityService{
		db: db,
	}
}
