package response

import "github.com/dashug/ldap-admin-platform/model"

type MenuListRsp struct {
	Total int64        `json:"total"`
	Menus []model.Menu `json:"menus"`
}
