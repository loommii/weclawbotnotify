package config

import (
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/mapping"
)

func TestAuthConfigTimeDuration(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantAcc time.Duration
		wantRef time.Duration
	}{
		{
			name: "秒格式",
			yaml: `
PrivateKeyPath: etc/keys/private.pem
PublicKeyPath: etc/keys/public.pem
AccessExpire: 900s
RefreshExpire: 604800s
`,
			wantAcc: 900 * time.Second,
			wantRef: 604800 * time.Second,
		},
		{
			name: "分钟和小时格式",
			yaml: `
PrivateKeyPath: etc/keys/private.pem
PublicKeyPath: etc/keys/public.pem
AccessExpire: 1h
RefreshExpire: 168h
`,
			wantAcc: 1 * time.Hour,
			wantRef: 168 * time.Hour,
		},
		{
			name: "组合格式",
			yaml: `
PrivateKeyPath: etc/keys/private.pem
PublicKeyPath: etc/keys/public.pem
AccessExpire: 15m
RefreshExpire: 24h
`,
			wantAcc: 15 * time.Minute,
			wantRef: 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c AuthConfig
			err := mapping.UnmarshalYamlBytes([]byte(tt.yaml), &c)
			if err != nil {
				t.Fatalf("UnmarshalYamlBytes 失败: %v", err)
			}
			if c.AccessExpire != tt.wantAcc {
				t.Errorf("AccessExpire = %v, want %v", c.AccessExpire, tt.wantAcc)
			}
			if c.RefreshExpire != tt.wantRef {
				t.Errorf("RefreshExpire = %v, want %v", c.RefreshExpire, tt.wantRef)
			}
		})
	}
}
