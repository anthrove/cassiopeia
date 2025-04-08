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
	"database/sql"
	"fmt"
	"github.com/anthrove/identity/docs"
	"github.com/anthrove/identity/pkg/logic"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zitadel/oidc/v3/pkg/op"
	"html/template"
	"net/http"
	"time"
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

	v1.POST("/tenant/:tenant_id/provider", identityRoutes.createProvider)
	v1.GET("/tenant/:tenant_id/provider", Pagination(), identityRoutes.findProviders)
	v1.GET("/tenant/:tenant_id/provider/:provider_id", identityRoutes.findProvider)
	v1.PUT("/tenant/:tenant_id/provider/:provider_id", identityRoutes.updateProvider)
	v1.DELETE("/tenant/:tenant_id/provider/:provider_id", identityRoutes.killProvider)
	v1.POST("/tenant/:tenant_id/provider/:provider_id/mail", identityRoutes.SendMail)

	v1.POST("/tenant/:tenant_id/certificate", identityRoutes.createCertificate)
	v1.GET("/tenant/:tenant_id/certificate", Pagination(), identityRoutes.findCertificates)
	v1.GET("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.findCertificate)
	v1.PUT("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.updateCertificate)
	v1.DELETE("/tenant/:tenant_id/certificate/:certificate_id", identityRoutes.killCertificate)

	v1.POST("/tenant/:tenant_id/template", identityRoutes.createMessageTemplate)
	v1.GET("/tenant/:tenant_id/template", Pagination(), identityRoutes.findMessageTemplates)
	v1.GET("/tenant/:tenant_id/template/:template_id", identityRoutes.findMessageTemplate)
	v1.PUT("/tenant/:tenant_id/template/:template_id", identityRoutes.updateMessageTemplate)
	v1.DELETE("/tenant/:tenant_id/template/:template_id", identityRoutes.killMessageTemplate)
	v1.POST("/tenant/:tenant_id/template/:template_id/fill", identityRoutes.fillMessageTemplate)

	v1.POST("/tenant/:tenant_id/application", identityRoutes.createApplication)
	v1.GET("/tenant/:tenant_id/application", Pagination(), identityRoutes.findApplication)
	v1.GET("/tenant/:tenant_id/application/:application_id", identityRoutes.findApplications)
	v1.PUT("/tenant/:tenant_id/application/:application_id", identityRoutes.updateApplication)
	v1.DELETE("/tenant/:tenant_id/application/:application_id", identityRoutes.killApplication)

	v1.POST("/tenant/:tenant_id/resource", identityRoutes.createResource)
	v1.GET("/tenant/:tenant_id/resource", Pagination(), identityRoutes.findResources)
	v1.GET("/tenant/:tenant_id/resource/:resource_id", identityRoutes.findResource)
	v1.DELETE("/tenant/:tenant_id/resource/:resource_id", identityRoutes.killResource)

	v1.POST("/tenant/:tenant_id/model", identityRoutes.createModel)
	v1.GET("/tenant/:tenant_id/model", Pagination(), identityRoutes.findModel)
	v1.GET("/tenant/:tenant_id/model/:model_id", identityRoutes.findModel)
	v1.PUT("/tenant/:tenant_id/model/:model_id", identityRoutes.updateModel)
	v1.DELETE("/tenant/:tenant_id/model/:model_id", identityRoutes.killModel)

	v1.POST("/tenant/:tenant_id/application/:application_id/login", identityRoutes.signIn)

	v1.GET("/cdn/:tenant_id/*file_path", identityRoutes.cdnGetFile)

	r.Any("/favicon.ico", func(context *gin.Context) {})

	r.Any("/auth/:tenant_id/login", func(c *gin.Context) {
		switch c.Request.Method {
		case "GET":
			loginHandler(c.Writer, c.Request)
		case "POST":
			identityRoutes.checkLoginHandler(c)
		default:
			c.Status(http.StatusMethodNotAllowed)
			return
		}
	})
	r.Any("/:tenant_id/*any", identityRoutes.OIDCEndpoints)

}

var (
	loginTmpl, _ = template.New("login").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Login</title>
		</head>
		<body style="display: flex; align-items: center; justify-content: center; height: 100vh;">
			<form method="POST" style="height: 200px; width: 200px;">
				<input type="hidden" name="id" value="{{.ID}}">
				<div>
					<label for="username">Username:</label>
					<input id="username" name="username" style="width: 100%">
				</div>
				<div>
					<label for="password">Password:</label>
					<input id="password" name="password" style="width: 100%">
				</div>
				<p style="color:red; min-height: 1rem;">{{.Error}}</p>
				<button type="submit">Login</button>
			</form>
		</body>
	</html>`)
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.URL.Query().Get("request_id")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	//the oidc package will pass the id of the auth request as query parameter
	//we will use this id through the login process and therefore pass it to the  login page
	renderLogin(w, requestID, nil)
}

func renderLogin(w http.ResponseWriter, id string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg,
	}
	err = loginTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ir IdentityRoutes) checkLoginHandler(c *gin.Context) {
	tenantId := c.Param("tenant_id")

	authRequestID, _ := c.GetPostForm("id")

	authRequest, err := ir.service.FindAuthRequest(c, tenantId, authRequestID)

	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")

	_, user, err := ir.service.SignIn(c, tenantId, authRequest.ApplicationID, object.SignInRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		renderLogin(c.Writer, authRequestID, err)
		return
	}

	err = ir.service.UpdateAuthRequest(c, tenantId, authRequestID, object.UpdateAuthRequest{
		UserID: sql.NullString{
			String: user.ID,
			Valid:  true,
		},
		Authenticated:   true,
		AuthenticatedAt: time.Now(),
	})

	if err != nil {
		renderLogin(c.Writer, authRequestID, err)
		return
	}

	provider, err := GetProvider(c, ir.service, tenantId)
	if err != nil {
		renderLogin(c.Writer, authRequestID, err)
	}

	callbackURL := op.AuthCallbackURL(provider)
	c.Redirect(http.StatusFound, "/"+tenantId+"/"+callbackURL(c, authRequestID))
}
