package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/dashug/ldap-admin-platform/config"
	"github.com/dashug/ldap-admin-platform/model"
	"github.com/dashug/ldap-admin-platform/public/common"
	"github.com/robfig/cron/v3"
)

// 统一定时自动同步：在 InitCron 创建的 cron 实例上注册一个「全量自动同步」任务，
// 触发时同步所有已启用（EnableSync）的来源（部门+用户）。与各平台原有的
// DeptSyncTime/UserSyncTime 互不影响，由「设置-定时同步」页统一管理。
var (
	autoSyncCron    *cron.Cron
	autoSyncEntryID cron.EntryID
	// 与 cron.New(cron.WithSeconds()) 一致的 6 段解析器，用于保存前校验表达式
	cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
)

// ValidateCron 校验 cron 表达式是否合法（6 段，含秒）
func ValidateCron(spec string) error {
	_, err := cronParser.Parse(strings.TrimSpace(spec))
	return err
}

// bindAutoSyncCron 由 InitCron 调用，绑定全局 cron 实例并按当前配置注册任务
func bindAutoSyncCron(c *cron.Cron) {
	autoSyncCron = c
	if err := ReloadAutoSync(); err != nil {
		common.Log.Errorf("注册定时自动同步任务失败: %v", err)
	}
}

// ReloadAutoSync 按当前配置(重新)注册定时自动同步任务，配置变更后调用即可热生效
func ReloadAutoSync() error {
	if autoSyncCron == nil {
		return nil // cron 尚未初始化（InitCron 之前），忽略
	}
	if autoSyncEntryID != 0 {
		autoSyncCron.Remove(autoSyncEntryID)
		autoSyncEntryID = 0
	}
	if config.Conf.System == nil || !config.Conf.System.AutoSyncEnabled {
		return nil
	}
	spec := strings.TrimSpace(config.Conf.System.AutoSyncCron)
	if spec == "" {
		return nil
	}
	id, err := autoSyncCron.AddFunc(spec, runAutoSync)
	if err != nil {
		return err
	}
	autoSyncEntryID = id
	common.Log.Infof("定时自动同步已启用: %s", spec)
	return nil
}

// runAutoSync 同步所有已启用来源（部门+用户）
func runAutoSync() {
	if config.Conf.Ldap != nil && config.Conf.Ldap.EnableSync {
		RunSourceSync("ldap", "auto")
	}
	if config.Conf.DingTalk != nil && config.Conf.DingTalk.EnableSync {
		RunSourceSync("dingtalk", "auto")
	}
	if config.Conf.WeCom != nil && config.Conf.WeCom.EnableSync {
		RunSourceSync("wecom", "auto")
	}
	if config.Conf.FeiShu != nil && config.Conf.FeiShu.EnableSync {
		RunSourceSync("feishu", "auto")
	}
}

// RunSourceSync 执行单个来源的「部门+用户」同步，并记录一条 SyncRun
func RunSourceSync(source, trigger string) {
	start := time.Now()
	msg, ok := syncOneSource(source)
	recordSyncRun(source, trigger, msg, ok, time.Since(start))
}

// syncOneSource 顺序同步指定来源的部门与用户，聚合错误信息
func syncOneSource(source string) (message string, success bool) {
	var errs []string
	addErr := func(label string, e any) {
		if e != nil {
			errs = append(errs, label+": "+fmt.Sprint(e))
		}
	}
	switch source {
	case "ldap":
		_, e1 := OpenLdap.SyncOpenLdapDepts(nil, nil)
		addErr("部门", e1)
		_, e2 := OpenLdap.SyncOpenLdapUsers(nil, nil)
		addErr("用户", e2)
	case "dingtalk":
		_, e1 := DingTalk.SyncDingTalkDepts(nil, nil)
		addErr("部门", e1)
		_, e2 := DingTalk.SyncDingTalkUsers(nil, nil)
		addErr("用户", e2)
	case "wecom":
		_, e1 := WeCom.SyncWeComDepts(nil, nil)
		addErr("部门", e1)
		_, e2 := WeCom.SyncWeComUsers(nil, nil)
		addErr("用户", e2)
	case "feishu":
		_, e1 := FeiShu.SyncFeiShuDepts(nil, nil)
		addErr("部门", e1)
		_, e2 := FeiShu.SyncFeiShuUsers(nil, nil)
		addErr("用户", e2)
	default:
		return "未知同步来源: " + source, false
	}
	if len(errs) == 0 {
		return "同步成功", true
	}
	return strings.Join(errs, "; "), false
}

// recordSyncRun 落一条同步运行记录（DB 未就绪时静默跳过）
func recordSyncRun(source, trigger, message string, success bool, dur time.Duration) {
	if common.DB == nil {
		return
	}
	if len(message) > 500 {
		message = message[:500]
	}
	rec := model.SyncRun{
		Source:   source,
		Trigger:  trigger,
		Success:  success,
		Message:  message,
		Duration: dur.Milliseconds(),
	}
	if err := common.DB.Create(&rec).Error; err != nil {
		common.Log.Warnf("同步运行记录写入失败: %v", err)
	}
}
