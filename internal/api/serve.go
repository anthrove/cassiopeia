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
	v1Auth := v1.Group("", identityRoutes.Authorization())
	v1.POST("/tenant", identityRoutes.createTenant)
	v1.GET("/tenant", Pagination(), identityRoutes.findTenants)
	v1Auth.GET("/tenant/:tenant_id", identityRoutes.findTenant)
	v1Auth.PUT("/tenant/:tenant_id", identityRoutes.updateTenant)
	v1.DELETE("/tenant/:tenant_id", identityRoutes.killTenant)

	v1Auth.POST("/tenant/:tenant_id/group", identityRoutes.createGroup)
	v1Auth.GET("/tenant/:tenant_id/group", Pagination(), identityRoutes.findGroups)
	v1Auth.GET("/tenant/:tenant_id/group/:group_id", identityRoutes.findGroup)
	v1Auth.PUT("/tenant/:tenant_id/group/:group_id", identityRoutes.updateGroup)
	v1Auth.DELETE("/tenant/:tenant_id/group/:group_id", identityRoutes.killGroup)

	v1Auth.POST("/tenant/:tenant_id/user", identityRoutes.createUser)
	v1Auth.GET("/tenant/:tenant_id/user", Pagination(), identityRoutes.findUsers)
	v1Auth.GET("/tenant/:tenant_id/user/:user_id", identityRoutes.findUser)
	v1Auth.PUT("/tenant/:tenant_id/user/:user_id", identityRoutes.updateUser)
	v1Auth.DELETE("/tenant/:tenant_id/user/:user_id", identityRoutes.killUser)

	// TODO: Add VerifieMFA endpoint
	v1Auth.POST("/tenant/:tenant_id/user/:user_id/mfa", identityRoutes.createMFA)
	v1Auth.GET("/tenant/:tenant_id/user/:user_id/mfa", Pagination(), identityRoutes.findMFAs)
	v1Auth.GET("/tenant/:tenant_id/user/:user_id/mfa/:mfa_id", identityRoutes.findMFA)
	v1Auth.PUT("/tenant/:tenant_id/user/:user_id/mfa/:mfa_id", identityRoutes.updateMFA)
	v1Auth.POST("/tenant/:tenant_id/user/:user_id/mfa/:mfa_id/verify", identityRoutes.verifyMFA)
	v1Auth.DELETE("/tenant/:tenant_id/user/:user_id/mfa/:mfa_id", identityRoutes.killMFA)

	v1Auth.POST("/tenant/:tenant_id/provider", identityRoutes.createProvider)
	v1Auth.GET("/tenant/:tenant_id/provider", Pagination(), identityRoutes.findProviders)
	v1Auth.GET("/tenant/:tenant_id/provider/category", Pagination(), identityRoutes.findProviderCategories)
	v1Auth.GET("/tenant/:tenant_id/provider/category/:category", Pagination(), identityRoutes.findProviderCategoryTypes)
	v1Auth.GET("/tenant/:tenant_id/provider/category/:category/:type", Pagination(), identityRoutes.findProviderCategoryTypeConfiguration)
	v1Auth.GET("/tenant/:tenant_id/provider/:provider_id", identityRoutes.findProvider)
	v1Auth.PUT("/tenant/:tenant_id/provider/:provider_id", identityRoutes.updateProvider)
	v1Auth.DELETE("/tenant/:tenant_id/provider/:provider_id", identityRoutes.killProvider)
	v1Auth.POST("/tenant/:tenant_id/provider/:provider_id/mail", identityRoutes.SendMail)

	v1Auth.POST("/tenant/:tenant_id/certificate", identityRoutes.createCertificate)
	v1Auth.GET("/tenant/:tenant_id/certificate", Pagination(), identityRoutes.findCertificates)
	v1Auth.GET("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.findCertificate)
	v1Auth.PUT("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.updateCertificate)
	v1Auth.DELETE("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.killCertificate)

	v1Auth.POST("/tenant/:tenant_id/template", identityRoutes.createMessageTemplate)
	v1Auth.GET("/tenant/:tenant_id/template", Pagination(), identityRoutes.findMessageTemplates)
	v1Auth.GET("/tenant/:tenant_id/template/:template_id", identityRoutes.findMessageTemplate)
	v1Auth.PUT("/tenant/:tenant_id/template/:template_id", identityRoutes.updateMessageTemplate)
	v1Auth.DELETE("/tenant/:tenant_id/template/:template_id", identityRoutes.killMessageTemplate)
	v1Auth.POST("/tenant/:tenant_id/template/:template_id/fill", identityRoutes.fillMessageTemplate)

	v1Auth.POST("/tenant/:tenant_id/application", identityRoutes.createApplication)
	v1Auth.GET("/tenant/:tenant_id/application", Pagination(), identityRoutes.findApplications)
	v1Auth.GET("/tenant/:tenant_id/application/:application_id", identityRoutes.findApplication)
	v1Auth.PUT("/tenant/:tenant_id/application/:application_id", identityRoutes.updateApplication)
	v1Auth.DELETE("/tenant/:tenant_id/application/:application_id", identityRoutes.killApplication)

	v1Auth.POST("/tenant/:tenant_id/resource", identityRoutes.createResource)
	v1Auth.GET("/tenant/:tenant_id/resource", Pagination(), identityRoutes.findResources)
	v1Auth.GET("/tenant/:tenant_id/resource/:resource_id", identityRoutes.findResource)
	v1Auth.DELETE("/tenant/:tenant_id/resource/:resource_id", identityRoutes.killResource)

	v1Auth.POST("/tenant/:tenant_id/model", identityRoutes.createModel)
	v1Auth.GET("/tenant/:tenant_id/model", Pagination(), identityRoutes.findModels)
	v1Auth.GET("/tenant/:tenant_id/model/:model_id", identityRoutes.findModel)
	v1Auth.PUT("/tenant/:tenant_id/model/:model_id", identityRoutes.updateModel)
	v1Auth.DELETE("/tenant/:tenant_id/model/:model_id", identityRoutes.killModel)

	v1Auth.POST("/tenant/:tenant_id/adapter", identityRoutes.createAdapter)
	v1Auth.GET("/tenant/:tenant_id/adapter", Pagination(), identityRoutes.findAdapters)
	v1Auth.GET("/tenant/:tenant_id/adapter/:adapter_id", identityRoutes.findAdapter)
	v1Auth.PUT("/tenant/:tenant_id/adapter/:adapter_id", identityRoutes.updateAdapter)
	v1Auth.DELETE("/tenant/:tenant_id/adapter/:adapter_id", identityRoutes.killAdapter)

	v1Auth.POST("/tenant/:tenant_id/permission", identityRoutes.createPermission)
	v1Auth.GET("/tenant/:tenant_id/permission", Pagination(), identityRoutes.findPermissions)
	v1Auth.GET("/tenant/:tenant_id/permission/:permission_id", identityRoutes.findPermission)
	v1Auth.PUT("/tenant/:tenant_id/permission/:permission_id", identityRoutes.updatePermission)
	v1Auth.DELETE("/tenant/:tenant_id/permission/:permission_id", identityRoutes.killPermission)

	v1Auth.POST("/tenant/:tenant_id/enforcer", identityRoutes.createEnforcer)
	v1Auth.GET("/tenant/:tenant_id/enforcer", Pagination(), identityRoutes.findEnforcers)
	v1Auth.GET("/tenant/:tenant_id/enforcer/:enforcer_id", identityRoutes.findEnforcer)
	v1Auth.PUT("/tenant/:tenant_id/enforcer/:enforcer_id", identityRoutes.updateEnforcer)
	v1Auth.DELETE("/tenant/:tenant_id/enforcer/:enforcer_id", identityRoutes.killEnforcer)
	v1Auth.POST("/tenant/:tenant_id/enforcer/:enforcer_id/enforce", identityRoutes.enforce)

	v1.GET("/tenant/:tenant_id/application/:application_id/login/begin", identityRoutes.signInBegin)
	v1.POST("/tenant/:tenant_id/application/:application_id/login", identityRoutes.signInSubmit)

	v1Auth.GET("/profile", identityRoutes.getProfileFields)
	v1Auth.POST("/profile", identityRoutes.upsertProfileFields)
	v1Auth.POST("/profile/mfa", identityRoutes.profileCreateMFA)
	v1Auth.POST("/profile/mfa/:mfa_id/verify", identityRoutes.profileVerifyMFA)
	v1Auth.POST("/profile/mfa/:mfa_id", identityRoutes.profileUpdateMFA)
	v1Auth.GET("/profile/mfa", Pagination(), identityRoutes.profileGetMFAs)
	v1Auth.DELETE("/profile/mfa/:mfa_id", identityRoutes.profileGetMFAs)

	v1.GET("/cdn/:tenant_id/*file_path", identityRoutes.cdnGetFile)

	r.Any("/favicon.ico", func(context *gin.Context) {})

	r.Any("/:tenant_id/*any", identityRoutes.OIDCEndpoints)

}
