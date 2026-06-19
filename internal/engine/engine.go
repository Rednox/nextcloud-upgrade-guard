package engine

import (
	"context"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/collector"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/resolver"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

func Check(ctx context.Context, c collector.Collector, r resolver.Resolver, targetMajor int) ([]types.CheckResult, error) {
	apps, err := c.Collect(ctx)
	if err != nil {
		return nil, err
	}
	return r.Resolve(ctx, targetMajor, apps)
}

func SummarizeResults(results []types.CheckResult) types.CheckSummary {
	summary := types.CheckSummary{}
	for _, result := range results {
		summary.Total++
		if result.Enabled {
			summary.Enabled++
		}
		switch result.Status {
		case types.StatusCompatible:
			summary.Compatible++
		case types.StatusUpgradableToCompatible:
			summary.UpgradableToCompatible++
		case types.StatusIncompatible:
			summary.Incompatible++
		default:
			summary.Unknown++
		}
	}
	return summary
}
