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

package main

import (
	"github.com/anthrove/identity/internal/api"
	"github.com/anthrove/identity/pkg/logic"
	"github.com/anthrove/identity/pkg/repository"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	engine, err := repository.GetEngine()

	if err != nil {
		log.Panic("Problem while connecting to database: ", err)
	}

	err = repository.Migrate(engine)

	if err != nil {
		log.Panic("Problem while migrating database: ", err)
	}

	service := logic.NewIdentityService(engine)

	router := gin.Default()
	api.SetupRoutes(router, service)
	err = router.Run(":8080")

	if err != nil {
		log.Panic("Problem while running server: ", err)
		return
	}
}
