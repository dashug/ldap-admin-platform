package common

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 全局CasbinEnforcer
// 使用 SyncedEnforcer：其 Enforce/策略增删等操作内部已用读写锁保证并发安全，
// 因此鉴权中间件无需再用全局互斥锁串行化请求，提升并发吞吐。
var CasbinEnforcer *casbin.SyncedEnforcer

// 初始化casbin策略管理器
func InitCasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		Log.Panicf("初始化Casbin失败：%v", err)
		panic(fmt.Sprintf("初始化Casbin失败：%v", err))
	}

	CasbinEnforcer = e
	Log.Info("初始化Casbin完成!")
}

var casbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

func mysqlCasbin() (*casbin.SyncedEnforcer, error) {
	a, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(casbinModel)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
