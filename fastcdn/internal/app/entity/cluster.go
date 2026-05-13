package entity

import (
	"fastcdn/internal/model"
)

type ClusterNodeEntityList struct {
	model.ClusterNode
	IpList []model.ClusterNodeIpaddr `json:"ip_list"`
}
