package controller

import (
	"github.com/dashug/ldap-admin-platform/logic"
	"github.com/dashug/ldap-admin-platform/model/request"

	"github.com/gin-gonic/gin"
)

type ApiKeyController struct{}

// List API Key 列表
func (m *ApiKeyController) List(c *gin.Context) {
	req := new(request.ApiKeyListReq)
	Run(c, req, func() (any, any) {
		return logic.ApiKey.List(c, req)
	})
}

// Create 创建 API Key（返回的 key 仅此一次）
func (m *ApiKeyController) Create(c *gin.Context) {
	req := new(request.ApiKeyCreateReq)
	Run(c, req, func() (any, any) {
		return logic.ApiKey.Create(c, req)
	})
}

// Delete 删除 API Key
func (m *ApiKeyController) Delete(c *gin.Context) {
	req := new(request.ApiKeyDeleteReq)
	Run(c, req, func() (any, any) {
		return logic.ApiKey.Delete(c, req)
	})
}
