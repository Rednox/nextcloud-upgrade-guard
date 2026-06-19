package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

type OCCCollector struct {
	OccPath string
	PHPBin  string
}

func NewOCCCollector(occPath, phpBin string) *OCCCollector {
	return &OCCCollector{OccPath: occPath, PHPBin: phpBin}
}

func (c *OCCCollector) Collect(ctx context.Context) ([]types.AppInfo, error) {
	cmd := exec.CommandContext(ctx, c.PHPBin, c.OccPath, "app:list", "--output=json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("collect apps via occ: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	var payload map[string]map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		return nil, fmt.Errorf("parse occ app:list json: %w", err)
	}

	apps := make([]types.AppInfo, 0)
	for _, enabled := range []bool{true, false} {
		bucket := "disabled"
		if enabled {
			bucket = "enabled"
		}
		for appID, rawVersion := range payload[bucket] {
			apps = append(apps, types.AppInfo{
				AppID:            appID,
				Enabled:          enabled,
				InstalledVersion: extractVersion(rawVersion),
				Source:           "occ",
			})
		}
	}

	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Enabled != apps[j].Enabled {
			return apps[i].Enabled
		}
		return apps[i].AppID < apps[j].AppID
	})

	return apps, nil
}

func extractVersion(raw any) string {
	switch v := raw.(type) {
	case string:
		return v
	case map[string]any:
		if version, ok := v["version"].(string); ok {
			return version
		}
	}
	return ""
}
