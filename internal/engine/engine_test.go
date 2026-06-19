package engine

import (
	"testing"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

func TestSummarizeResults(t *testing.T) {
	results := []types.CheckResult{
		{Enabled: true, Status: types.StatusCompatible},
		{Enabled: true, Status: types.StatusUpgradableToCompatible},
		{Enabled: true, Status: types.StatusIncompatible},
		{Enabled: true, Status: types.StatusUnknown},
	}

	summary := SummarizeResults(results)
	if summary.Total != 4 || summary.Enabled != 4 {
		t.Fatalf("unexpected totals: %+v", summary)
	}
	if summary.Compatible != 1 || summary.UpgradableToCompatible != 1 || summary.Incompatible != 1 || summary.Unknown != 1 {
		t.Fatalf("unexpected status counts: %+v", summary)
	}
}
