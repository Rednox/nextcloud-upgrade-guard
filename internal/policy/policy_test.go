package policy

import (
	"testing"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

func TestParsePolicy(t *testing.T) {
	raw := []byte(`block_on_incompatible: true
block_on_unknown: false
critical_apps:
  - oidc_login
  - twofactor_totp
block_on_critical_risk: true
`)

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !cfg.BlockOnIncompatible || cfg.BlockOnUnknown || !cfg.BlockOnCriticalRisk {
		t.Fatalf("unexpected bool fields: %+v", cfg)
	}
	if len(cfg.CriticalApps) != 2 || cfg.CriticalApps[0] != "oidc_login" {
		t.Fatalf("unexpected critical apps: %+v", cfg.CriticalApps)
	}
}

func TestPolicyEvaluatorExitCodes(t *testing.T) {
	results := []types.CheckResult{
		{AppID: "oidc_login", Enabled: true, Status: types.StatusUnknown},
		{AppID: "calendar", Enabled: true, Status: types.StatusIncompatible},
		{AppID: "files", Enabled: true, Status: types.StatusCompatible},
	}

	tests := []struct {
		name string
		cfg  Config
		want int
	}{
		{
			name: "critical risk first",
			cfg:  Config{BlockOnIncompatible: true, BlockOnUnknown: true, CriticalApps: []string{"oidc_login"}, BlockOnCriticalRisk: true},
			want: 12,
		},
		{
			name: "incompatible block",
			cfg:  Config{BlockOnIncompatible: true, BlockOnUnknown: false, BlockOnCriticalRisk: false},
			want: 10,
		},
		{
			name: "unknown block",
			cfg:  Config{BlockOnIncompatible: false, BlockOnUnknown: true, BlockOnCriticalRisk: false},
			want: 11,
		},
		{
			name: "policy pass",
			cfg:  Config{BlockOnIncompatible: false, BlockOnUnknown: false, BlockOnCriticalRisk: false},
			want: 0,
		},
	}

	evaluator := PolicyEvaluator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluator.Evaluate(tt.cfg, 35, results)
			if got.ExitCode != tt.want {
				t.Fatalf("Evaluate() exit code = %d, want %d", got.ExitCode, tt.want)
			}
		})
	}
}
