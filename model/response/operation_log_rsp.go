package response

import "github.com/dashug/ldap-admin-platform/model"

type LogListRsp struct {
	Total int64                `json:"total"`
	Logs  []model.OperationLog `json:"logs"`
}
