package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary	Get an file
// @Tags		CDN API
// @Accept		json
// @Produce	json
// @Param		tenant_id	path	string	true	"Tenant ID"
// @Router		/api/v1/cdn/{tenant_id}/{file_path} [get]
func (ir IdentityRoutes) cdnGetFile(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	filePath := c.Param("file_path")

	localFilePath, err := ir.service.ServeResource(c, tenantID, filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.File(localFilePath)
}
