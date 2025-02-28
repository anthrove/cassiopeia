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
	_ "github.com/lib/pq"

	"github.com/anthrove/identity/internal/config"
	"github.com/caarlos0/env/v11"
	"xorm.io/xorm"
)

func GetEngine() (*xorm.Engine, error) {
	dbConfig, err := env.ParseAs[config.Database]()

	if err != nil {
		return nil, err
	}

	engine, err := xorm.NewEngine(dbConfig.Driver, dbConfig.DataSource)

	if err != nil {
		return nil, err
	}

	return engine, nil
}

func Migrate(engine *xorm.Engine) error {

	return nil
}
