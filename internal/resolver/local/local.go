package local

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/resolver"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

var _ resolver.Resolver = (*Resolver)(nil)

type Resolver struct{}

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(ctx context.Context, targetMajor int, apps []types.AppInfo) ([]types.CheckResult, error) {
	if targetMajor <= 0 {
		return nil, fmt.Errorf("target major must be > 0")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	results := make([]types.CheckResult, 0)
	for _, app := range apps {
		if !app.Enabled {
			continue
		}

		status, reasons, action := inferStatus(app.InstalledVersion, targetMajor)
		results = append(results, types.CheckResult{
			App:               app,
			AppID:             app.AppID,
			Enabled:           app.Enabled,
			InstalledVersion:  app.InstalledVersion,
			Status:            status,
			Reasons:           reasons,
			RecommendedAction: action,
		})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].AppID < results[j].AppID })
	return results, nil
}

func inferStatus(version string, targetMajor int) (types.CompatibilityStatus, []string, string) {
	if version == "" {
		return types.StatusUnknown, []string{"installed version unavailable from local metadata"}, "verify app compatibility manually"
	}

	major, ok := parseMajor(version)
	if !ok {
		return types.StatusUnknown, []string{"unable to infer compatibility from version format"}, "verify app compatibility manually"
	}

	if major >= targetMajor {
		return types.StatusCompatible, []string{"installed major version appears compatible with target"}, "none"
	}
	if major == targetMajor-1 {
		return types.StatusUpgradableToCompatible, []string{"installed major suggests upgrade path may reach compatibility"}, "update app before core upgrade"
	}
	return types.StatusIncompatible, []string{"installed major version appears below target compatibility threshold"}, "find compatible update or disable app"
}

func parseMajor(version string) (int, bool) {
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return major, true
}

// TODO(v0.2): add remote resolver for App Store compatibility metadata.
// TODO(v0.2): add cache/retry/offline resolver modes.
// TODO(v0.2): add richer compatibility confidence scoring.
