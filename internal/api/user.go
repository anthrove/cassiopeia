package api

import (
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary	Creates a new User
// @Tags		User API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path		string								true	"Tenant ID"
//
// @Param		"User"		body		object.CreateUser					true	"Create User Data"
// @Success	200			{object}	HttpResponse{data=object.User{}}	"User"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/user [post]
func (ir IdentityRoutes) createUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateUser
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := ir.service.CreateUser(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: user,
	})
}

// @Summary	Update an existing User
// @Tags		User API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path	string				true	"Tenant ID"
// @Param		user_id		path	string				true	"User ID"
//
// @Param		"User"		body	object.UpdateUser	true	"Update User Data"
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/user/{user_id} [put]
func (ir IdentityRoutes) updateUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	var body object.UpdateUser
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateUser(c, tenantID, userID, body)

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

// @Summary	Delete an existing User
// @Tags		User API
// @Accept		json
// @Produce	json
//
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Param		user_id		path	string	true	"User ID"
//
// @Success	204
// @Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/user/{user_id} [delete]
func (ir IdentityRoutes) killUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	err := ir.service.KillUser(c, tenantID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary	Get an existing User
// @Tags		User API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Param		user_id		path		string								true	"User ID"
// @Success	200			{object}	HttpResponse{data=object.User{}}	"User"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/user/{user_id} [get]
func (ir IdentityRoutes) findUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("user_id")

	user, err := ir.service.FindUser(c, tenantID, userID)
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

// @Summary	Get existing Users
// @Tags		User API
// @Accept		json
// @Produce	json
//
// @Param		page		query		string								false	"Page"
// @Param		page_limit	query		string								false	"Page Limit"
//
// @Param		tenant_id	path		string								true	"Tenant ID"
// @Success	200			{object}	HttpResponse{data=[]object.User{}}	"User"
// @Failure	400			{object}	HttpResponse{data=nil}				"Bad Request"
// @Router		/api/v1/tenant/{tenant_id}/user [get]
func (ir IdentityRoutes) findUsers(c *gin.Context) {
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

	users, err := ir.service.FindUsers(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: users,
	})
}
