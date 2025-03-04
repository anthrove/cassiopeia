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
	"errors"
	"github.com/anthrove/identity/pkg/object"
	"github.com/gin-gonic/gin"
	"net/http"
)

//	@Summary	Creates a new Certificate
//	@Tags		Certificate API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id		path		string									true	"Tenant ID"
//
//	@Param		"Certificate"	body		object.CreateCertificate				true	"Create Certificate Data"
//	@Success	200				{object}	HttpResponse{data=object.Certificate{}}	"Certificate"
//	@Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/certificate [post]
func (ir IdentityRoutes) createCertificate(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var body object.CreateCertificate
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	certificate, err := ir.service.CreateCertificate(c, tenantID, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, HttpResponse{
		Data: certificate,
	})
}

//	@Summary	Update an existing Certificate
//	@Tags		Certificate API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id		path	string						true	"Tenant ID"
//	@Param		certificate_id	path	string						true	"Certificate ID"
//
//	@Param		"Certificate"	body	object.UpdateCertificate	true	"Create Certificate Data"
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/certificate/{certificate_id} [put]
func (ir IdentityRoutes) updateCertificate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	certificateID := c.Param("certificate_id")

	var body object.UpdateCertificate
	err := c.ShouldBind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	err = ir.service.UpdateCertificate(c, tenantID, certificateID, body)

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

//	@Summary	Kill an existing Certificate
//	@Tags		Certificate API
//	@Accept		json
//	@Produce	json
//
//	@Param		tenant_id		path	string	true	"Tenant ID"
//	@Param		certificate_id	path	string	true	"Certificate ID"
//
//	@Success	204
//	@Failure	400	{object}	HttpResponse{data=nil}	"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/certificate/{certificate_id} [delete]
func (ir IdentityRoutes) killCertificate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	certificateID := c.Param("certificate_id")

	err := ir.service.KillCertificate(c, tenantID, certificateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

//	@Summary	Get an existing Certificate
//	@Tags		Certificate API
//	@Accept		json
//	@Produce	json
//	@Param		tenant_id		path		string									true	"Tenant ID"
//	@Param		certificate_id	path		string									true	"Certificate ID"
//	@Success	200				{object}	HttpResponse{data=object.Certificate{}}	"Certificate"
//	@Failure	400				{object}	HttpResponse{data=nil}					"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/certificate/{certificate_id} [get]
func (ir IdentityRoutes) findCertificate(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	certificateID := c.Param("certificate_id")

	certificate, err := ir.service.FindCertificate(c, tenantID, certificateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: certificate,
	})
}

//	@Summary	Get existing Certificates
//	@Tags		Certificate API
//	@Accept		json
//	@Produce	json
//
//	@Param		page		query		string										false	"Page"
//	@Param		page_limit	query		string										false	"Page Limit"
//
//	@Param		tenant_id	path		string										true	"Tenant ID"
//	@Success	200			{object}	HttpResponse{data=[]object.Certificate{}}	"Certificate"
//	@Failure	400			{object}	HttpResponse{data=nil}						"Bad Request"
//	@Router		/api/v1/tenant/{tenant_id}/certificate [get]
func (ir IdentityRoutes) findCertificates(c *gin.Context) {
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

	certificates, err := ir.service.FindCertificates(c, tenantID, paginationObj)

	if err != nil {
		c.JSON(http.StatusInternalServerError, HttpResponse{
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, HttpResponse{
		Data: certificates,
	})
}
