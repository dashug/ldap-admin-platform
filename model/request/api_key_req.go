package request

// ApiKeyListReq 列表（分页）
type ApiKeyListReq struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// ApiKeyCreateReq 创建 API Key（仅名称，密钥由服务端生成并仅返回一次）
type ApiKeyCreateReq struct {
	Name string `json:"name" validate:"required,min=1,max=64"`
}

// ApiKeyDeleteReq 删除
type ApiKeyDeleteReq struct {
	ID uint `json:"id" validate:"required"`
}
