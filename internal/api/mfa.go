package api

import (
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary	Creates a new MFA for a User
// @Tags		MFA API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string							true	"Tenant ID"
// @Param		user_id		path		string							true	"User ID"
// @Param		"MFA"		body		object.CreateMFA				true	"Create MFA Data"
// @Success	200			{object}	HttpResponse{data=object.MFA{}}	"MFA"
// @Failure	400			{object}	HttpResponse{data=nil}			"Bad Request"
// @Router		/tenant/{tenant_id}/user/{user_id}/mfa [post]
func (ir IdentityRoutes) createMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	var body object.CreateMFA
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	mfa, err := ir.service.CreateMFA(c, tenantID, userID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: mfa,
	})
}

// @Summary	Update an existing MFA for a User
// @Tags		MFA API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string				true	"Tenant ID"
// @Param		user_id		path	string				true	"User ID"
// @Param		mfa_id		path	string				true	"MFA ID"
// @Param		"MFA"		body	object.UpdateMFA	true	"Update MFA Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [put]
func (ir IdentityRoutes) updateMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	var body object.UpdateMFA
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateMFA(c, tenantID, userID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: gin.H{},
	})
}

// @Summary	Delete an existing MFA for User
// @Tags		MFA API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		user_id		path	string	true	"User ID"
// @Param		mfa_id		path	string	true	"MFA ID"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [delete]
func (ir IdentityRoutes) killMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	err := ir.service.KillMFA(c, tenantID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing MFA for User
// @Tags		MFA API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string							true	"Tenant ID"
// @Param		user_id		path		string							true	"User ID"
// @Param		mfa_id		path		string							true	"MFA ID"
// @Success	200			{object}	HttpResponse{data=object.MFA{}}	"MFA"
// @Failure	400			{object}	HttpResponse{data=nil}			"Bad Request"
// @Router		/tenant/{tenant_id}/user/{user_id}/mfa/{mfa_id} [get]
func (ir IdentityRoutes) findMFA(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	user, err := ir.service.FindMFA(c, tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: user,
	})
}

// @Summary	Get existing MFAs
// @Tags		MFA API
// @Accept		json
// @Produce	json
// @Param		page		query		string								false	"Page"
// @Param		page_limit	query		string								false	"Page Limit"
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Param		user_id		path		string								true	"User ID"
// @Success	200			{object}	HttpResponse{data=[]object.MFA{}}	"MFA"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/tenant/{tenant_id}/user/{user_id}/mfa [get]
func (ir IdentityRoutes) findMFAs(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	pagination, ok := c.Get("pagination")
	if !ok {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: "pagination parameter is missing",
		})
		return
	}

	paginationObj, ok := pagination.(object.Pagination)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.New("pagination parameter cant be converted to object.Pagination"))
		return
	}

	users, err := ir.service.FindMFAs(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: users,
	})
}
