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

package repository

import (
	"github.com/anthrove/identity/pkg/object"
	"gorm.io/gorm"
)

const (
	MaxPageSize     = 20
	DefaultPageSize = 10
)

// Pagination applies pagination, sorting, and counting to a GORM query.
//
// Parameters:
// - pagination: A pointer to a Pagination struct containing page, limit, sort, total rows, and total pages information.
//
// The function calculates the offset based on the Page and Limit of the Pagination struct, and applies the offset and limit to the GORM DB instance.
func Pagination(pagination object.Pagination) func(db *gorm.DB) *gorm.DB {
	switch {
	case pagination.Limit > MaxPageSize:
		pagination.Limit = MaxPageSize
	case pagination.Limit <= 0:
		pagination.Limit = DefaultPageSize
	}

	offset := (pagination.Page - 1) * pagination.Limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pagination.Limit)
	}
}
