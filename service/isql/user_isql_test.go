package isql

import (
	"strings"
	"testing"
)

// deptIDFilter 应按「整段 token」匹配逗号分隔的 department_id，避免 "1" 命中 "12"。
func TestDeptIDFilter(t *testing.T) {
	cond, args := deptIDFilter([]uint{1})
	wantCond := "(department_id = ? OR department_id LIKE ? OR department_id LIKE ? OR department_id LIKE ?)"
	if cond != wantCond {
		t.Errorf("单部门条件不符:\n got: %s\nwant: %s", cond, wantCond)
	}
	want := []any{"1", "1,%", "%,1", "%,1,%"}
	if len(args) != len(want) {
		t.Fatalf("参数个数=%d, want %d", len(args), len(want))
	}
	for i := range want {
		if args[i] != want[i] {
			t.Errorf("args[%d]=%v, want %v", i, args[i], want[i])
		}
	}

	// 多部门用 OR 连接，各产生 4 个占位参数
	condN, argsN := deptIDFilter([]uint{1, 12})
	if strings.Count(condN, " OR department_id LIKE ?") != 6 || strings.Count(condN, ") OR (") != 1 {
		t.Errorf("多部门条件结构异常: %s", condN)
	}
	if len(argsN) != 8 {
		t.Errorf("多部门参数个数=%d, want 8", len(argsN))
	}
}
