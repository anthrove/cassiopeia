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
	"errors"
	"github.com/anthrove/identity/internal/config"
	"github.com/anthrove/identity/pkg/object"
	"github.com/caarlos0/env/v11"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetEngine initializes and returns a database connection based on the configuration.
// It supports multiple database drivers such as PostgreSQL and SQLite.
//
// Returns:
//   - A gorm.DB instance representing the database connection.
//   - An error if there is any issue during the initialization.
func GetEngine() (*gorm.DB, error) {
	dbConfig, err := env.ParseAs[config.Database]()

	if err != nil {
		return nil, err
	}

	var db *gorm.DB

	switch dbConfig.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbConfig.DataSource), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbConfig.DataSource), &gorm.Config{})

	default:
		return nil, errors.New("unknown database driver: " + dbConfig.Driver)
	}

	if err != nil {
		return nil, err
	}

	if dbConfig.Debug {
		db = db.Debug()
	}

	return db, nil
}

// Migrate applies the database schema migrations for the specified entities.
// It ensures that the database schema is up-to-date with the defined models.
//
// Parameters:
//   - engine: a gorm.DB instance representing the database connection.
//
// Returns:
//   - An error if there is any issue during the migration.
func Migrate(engine *gorm.DB) error {
	return engine.AutoMigrate(&object.Tenant{},
		&object.Group{},
		&object.User{},
		&object.Provider{},
		&object.Certificate{},
		&object.MessageTemplate{},
		&object.Application{},
		&object.AuthRequest{},
		&object.Token{},
		&object.Resource{},
		&object.Model{},
		&object.Enforcer{},
		&object.Permission{},
		&object.Credentials{})
}
