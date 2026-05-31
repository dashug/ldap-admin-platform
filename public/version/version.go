package version

import (
	"runtime"
	"runtime/debug"
)

var (
	Version   string // 版本号（如v1.0.0）
	GitCommit string // Git提交哈希（短格式）
	BuildTime string // 构建时间
	GoVersion string // Go编译器版本
	Platform  string // 编译平台（如linux/amd64）
)

// fillDefaults 在未通过 ldflags 注入时，用运行时与构建信息兜底，
// 保证「系统信息」页面不出现空白；ldflags 显式设置时优先保留。
func fillDefaults() {
	if GoVersion == "" {
		GoVersion = runtime.Version()
	}
	if Platform == "" {
		Platform = runtime.GOOS + "/" + runtime.GOARCH
	}
	// Go 1.18+ 在 VCS 仓库内 go build 会自动嵌入提交哈希与提交时间
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if GitCommit == "" && len(s.Value) >= 7 {
					GitCommit = s.Value[:7]
				}
			case "vcs.time":
				if BuildTime == "" {
					BuildTime = s.Value
				}
			}
		}
	}
	if Version == "" {
		Version = "dev"
	}
	if GitCommit == "" {
		GitCommit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "未知"
	}
}

// GetVersion 获取版本信息
func GetVersion() map[string]string {
	fillDefaults()
	return map[string]string{
		"version":   Version,
		"gitCommit": GitCommit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
		"platform":  Platform,
	}
}
