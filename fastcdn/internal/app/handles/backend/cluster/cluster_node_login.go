package cluster

import (
	"errors"

	"github.com/gin-gonic/gin"

	"fastcdn/internal/app/common"
	"fastcdn/internal/app/form"
	"fastcdn/internal/db"
)

func PostNodeLoginAdd(c *gin.Context) {
	var field form.ClusterNodeLoginAdd
	if err := c.ShouldBind(&field); err != nil {
		common.ErrorResp(c, err, -1)
		return
	}

	if field.NodeID < 1 {
		common.ErrorResp(c, errors.New("add exception[node_id]!"), -1)
		return
	}

	err := db.ClusterNodeLoginAddOrUpdate(field.NodeID, field.SshHost, field.SshPort, field.SshID)
	if err != nil {
		common.ErrorResp(c, err, -1)
		return
	}
	common.SuccessResp(c)
}
