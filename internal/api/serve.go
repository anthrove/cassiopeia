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

package api

import (
	"github.com/anthrove/identity/docs"
	"github.com/anthrove/identity/pkg/logic"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// IdentityRoutes defines the routes related to identity management.
// It holds a reference to the IdentityService which contains the business logic.
type IdentityRoutes struct {
	service logic.IdentityService
}

// SetupRoutes initializes the API routes for identity management.
// It sets up the routes for creating and updating tenants.
//
// Parameters:
//   - r: a gin.Engine instance representing the HTTP router.
//   - service: an instance of IdentityService containing the business logic for identity management.
func SetupRoutes(r *gin.Engine, service logic.IdentityService) {
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	identityRoutes := &IdentityRoutes{service}

	v1 := r.Group("/api/v1")
	v1.POST("/tenant", identityRoutes.createTenant)
	v1.GET("/tenant", Pagination(), identityRoutes.findTenants)
	v1.GET("/tenant/:tenant_id", identityRoutes.findTenant)
	v1.PUT("/tenant/:tenant_id", identityRoutes.updateTenant)
	v1.DELETE("/tenant/:tenant_id", identityRoutes.killTenant)

	v1.POST("/tenant/:tenant_id/group", identityRoutes.createGroup)
	v1.GET("/tenant/:tenant_id/group", Pagination(), identityRoutes.findGroups)
	v1.GET("/tenant/:tenant_id/group/:group_id", identityRoutes.findGroup)
	v1.PUT("/tenant/:tenant_id/group/:group_id", identityRoutes.updateGroup)
	v1.DELETE("/tenant/:tenant_id/group/:group_id", identityRoutes.killGroup)

	v1.POST("/tenant/:tenant_id/user", identityRoutes.createUser)
	v1.GET("/tenant/:tenant_id/user", Pagination(), identityRoutes.findUsers)
	v1.GET("/tenant/:tenant_id/user/:user_id", identityRoutes.findUser)
	v1.PUT("/tenant/:tenant_id/user/:user_id", identityRoutes.updateUser)
	v1.DELETE("/tenant/:tenant_id/user/:user_id", identityRoutes.killUser)

}
