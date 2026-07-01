package model

import (
	"testing"
	"time"
)

func TestApiKeyIsUsable(t *testing.T) {
	tru, fal := true, false
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)

	cases := []struct {
		name string
		key  ApiKey
		want bool
	}{
		{"默认(未设启用/过期)可用", ApiKey{}, true},
		{"显式启用可用", ApiKey{Enabled: &tru}, true},
		{"停用不可用", ApiKey{Enabled: &fal}, false},
		{"已过期不可用", ApiKey{ExpiresAt: &past}, false},
		{"未过期可用", ApiKey{ExpiresAt: &future}, true},
		{"启用且未过期可用", ApiKey{Enabled: &tru, ExpiresAt: &future}, true},
		{"停用即便未过期也不可用", ApiKey{Enabled: &fal, ExpiresAt: &future}, false},
	}
	for _, c := range cases {
		if got := c.key.IsUsable(); got != c.want {
			t.Errorf("%s: IsUsable()=%v, want %v", c.name, got, c.want)
		}
	}
}
