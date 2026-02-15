package logic

import (
	"fmt"

	"github.com/dashug/ldap-admin-platform/model/request"
	"github.com/dashug/ldap-admin-platform/model/response"
	"github.com/dashug/ldap-admin-platform/public/tools"
	"github.com/dashug/ldap-admin-platform/service/isql"

	"github.com/gin-gonic/gin"
)

type ApiKeyLogic struct{}

// List 列表
func (l ApiKeyLogic) List(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiKeyListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	page := tools.NewPageOption(r.PageNum, r.PageSize)
	list, err := isql.ApiKey.List(page.PageNum, page.PageSize)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取 API Key 列表失败: %s", err.Error()))
	}
	total, err := isql.ApiKey.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取总数失败"))
	}
	items := make([]response.ApiKeyItemRsp, 0, len(list))
	for _, k := range list {
		items = append(items, response.ApiKeyItemRsp{
			ID:        k.ID,
			Name:      k.Name,
			KeyPrefix: k.KeyPrefix,
			CreatedAt: k.CreatedAt,
		})
	}
	return response.ApiKeyListRsp{Total: total, Items: items}, nil
}

// Create 创建；返回的 Key 仅此一次，需提示用户保存
func (l ApiKeyLogic) Create(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiKeyCreateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	rawKey, err := tools.GenerateRandomApiKey()
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("生成密钥失败: %s", err.Error()))
	}
	ak, err := isql.ApiKey.Add(r.Name, rawKey)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建 API Key 失败: %s", err.Error()))
	}
	return response.ApiKeyCreateRsp{
		ID:        ak.ID,
		Name:      ak.Name,
		KeyPrefix: ak.KeyPrefix,
		Key:       rawKey,
		CreatedAt: ak.CreatedAt,
	}, nil
}

// Delete 删除
func (l ApiKeyLogic) Delete(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.ApiKeyDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	if err := isql.ApiKey.Delete(r.ID); err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除失败: %s", err.Error()))
	}
	return nil, nil
}
