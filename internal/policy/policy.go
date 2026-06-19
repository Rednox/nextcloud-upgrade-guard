package policy

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BlockOnIncompatible bool     `yaml:"block_on_incompatible" json:"block_on_incompatible"`
	BlockOnUnknown      bool     `yaml:"block_on_unknown" json:"block_on_unknown"`
	CriticalApps        []string `yaml:"critical_apps" json:"critical_apps"`
	BlockOnCriticalRisk bool     `yaml:"block_on_critical_risk" json:"block_on_critical_risk"`
}

type Evaluator interface {
	Evaluate(config Config, targetMajor int, results []types.CheckResult) types.GateDecision
}

type PolicyEvaluator struct{}

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read policy file %q: %w", path, err)
	}
	return Parse(content)
}

func Parse(content []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse policy yaml: %w", err)
	}
	for i, appID := range cfg.CriticalApps {
		cfg.CriticalApps[i] = strings.TrimSpace(appID)
	}
	return cfg, nil
}

func (PolicyEvaluator) Evaluate(config Config, targetMajor int, results []types.CheckResult) types.GateDecision {
	summary := summarize(results)
	critical := map[string]struct{}{}
	for _, appID := range config.CriticalApps {
		if appID != "" {
			critical[appID] = struct{}{}
		}
	}

	criticalFailing := make([]string, 0)
	incompatibleFailing := make([]string, 0)
	unknownFailing := make([]string, 0)

	for _, result := range results {
		switch result.Status {
		case types.StatusIncompatible:
			incompatibleFailing = append(incompatibleFailing, result.AppID)
			if _, ok := critical[result.AppID]; ok {
				criticalFailing = append(criticalFailing, result.AppID)
			}
		case types.StatusUnknown:
			unknownFailing = append(unknownFailing, result.AppID)
			if _, ok := critical[result.AppID]; ok {
				criticalFailing = append(criticalFailing, result.AppID)
			}
		}
	}

	sort.Strings(criticalFailing)
	sort.Strings(incompatibleFailing)
	sort.Strings(unknownFailing)

	if config.BlockOnCriticalRisk && len(criticalFailing) > 0 {
		return decision(targetMajor, summary, 12, "critical app risk detected", criticalFailing)
	}
	if config.BlockOnIncompatible && len(incompatibleFailing) > 0 {
		return decision(targetMajor, summary, 10, "incompatible apps blocked by policy", incompatibleFailing)
	}
	if config.BlockOnUnknown && len(unknownFailing) > 0 {
		return decision(targetMajor, summary, 11, "unknown apps blocked by policy", unknownFailing)
	}
	return decision(targetMajor, summary, 0, "policy passed", nil)
}

func decision(target int, summary types.CheckSummary, code int, reason string, failing []string) types.GateDecision {
	return types.GateDecision{
		TargetMajor: target,
		Pass:        code == 0,
		ExitCode:    code,
		Reason:      reason,
		FailingApps: failing,
		Summary:     summary,
	}
}

func summarize(results []types.CheckResult) types.CheckSummary {
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
