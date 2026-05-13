package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fastcdn/internal/app/common"
)

func ConsoleIndex(c *gin.Context) {
	data := common.CommonVer(c)
	c.HTML(http.StatusOK, "backend/console/index.tmpl", data)
}
