package logic

import "testing"

// appendCSV 向逗号分隔字符串追加值，应去空段与去重。
func TestAppendCSV(t *testing.T) {
	cases := []struct {
		csv, val, want string
	}{
		{"", "5", "5"},
		{"1,2", "3", "1,2,3"},
		{"1,2", "2", "1,2"},        // 去重
		{",1,", "2", "1,2"},        // 去掉空段与前导逗号
		{"1, 2 ,3", "2", "1,2,3"},  // 去空白后去重
		{"1,2,3", "", "1,2,3"},     // 追加空值不改变
	}
	for _, c := range cases {
		if got := appendCSV(c.csv, c.val); got != c.want {
			t.Errorf("appendCSV(%q,%q)=%q, want %q", c.csv, c.val, got, c.want)
		}
	}
}
