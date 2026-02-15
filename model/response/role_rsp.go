package response

import "github.com/dashug/ldap-admin-platform/model"

type RoleListRsp struct {
	Total int64        `json:"total"`
	Roles []model.Role `json:"roles"`
}
