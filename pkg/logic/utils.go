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
	"context"
	"gorm.io/gorm"
)

func saveDBConn(ctx context.Context, db *gorm.DB) context.Context {
	if ctx.Value("db_conn") == nil {
		return context.WithValue(ctx, "db_conn", db)
	}

	return ctx
}

func (is IdentityService) getDBConn(ctx context.Context) (*gorm.DB, bool) {
	if ctx.Value("db_conn") == nil {
		return is.db, false
	}

	dbVal, ok := ctx.Value("db_conn").(*gorm.DB)

	if !ok {
		panic("failed to convert db connection from context")
	}

	return dbVal, true
}
