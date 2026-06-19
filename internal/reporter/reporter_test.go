package reporter

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

func TestCheckReportJSONGolden(t *testing.T) {
	results := []types.CheckResult{
		{
			AppID:             "calendar",
			Enabled:           true,
			InstalledVersion:  "34.2.1",
			Status:            types.StatusUpgradableToCompatible,
			Reasons:           []string{"installed major suggests upgrade path may reach compatibility"},
			RecommendedAction: "update app before core upgrade",
		},
		{
			AppID:             "oidc_login",
			Enabled:           true,
			Status:            types.StatusUnknown,
			Reasons:           []string{"installed version unavailable from local metadata"},
			RecommendedAction: "verify app compatibility manually",
		},
	}

	report := BuildCheckReportAt(35, time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC), results)
	var got bytes.Buffer
	if err := New().RenderCheck(&got, report, "json"); err != nil {
		t.Fatalf("RenderCheck() error = %v", err)
	}

	goldenPath := filepath.Join("testdata", "check_report.golden.json")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}
	if got.String() != string(want) {
		t.Fatalf("json output mismatch\n--- got ---\n%s\n--- want ---\n%s", got.String(), string(want))
	}
}
