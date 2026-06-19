package types

import "time"

type AppInfo struct {
	AppID            string `json:"app_id"`
	Enabled          bool   `json:"enabled"`
	InstalledVersion string `json:"installed_version,omitempty"`
	Source           string `json:"source"`
}

type CompatibilityStatus string

const (
	StatusCompatible             CompatibilityStatus = "compatible"
	StatusUpgradableToCompatible CompatibilityStatus = "upgradable_to_compatible"
	StatusIncompatible           CompatibilityStatus = "incompatible"
	StatusUnknown                CompatibilityStatus = "unknown"
)

type CheckResult struct {
	App               AppInfo             `json:"-"`
	AppID             string              `json:"app_id"`
	Enabled           bool                `json:"enabled"`
	InstalledVersion  string              `json:"installed_version,omitempty"`
	Status            CompatibilityStatus `json:"status"`
	Reasons           []string            `json:"reasons"`
	RecommendedAction string              `json:"recommended_action"`
}

type CheckSummary struct {
	Total                  int `json:"total"`
	Enabled                int `json:"enabled"`
	Compatible             int `json:"compatible"`
	UpgradableToCompatible int `json:"upgradable_to_compatible"`
	Incompatible           int `json:"incompatible"`
	Unknown                int `json:"unknown"`
}

type CheckReport struct {
	TargetMajor int           `json:"target_major"`
	GeneratedAt time.Time     `json:"generated_at"`
	Apps        []CheckResult `json:"apps"`
	Summary     CheckSummary  `json:"summary"`
}

type GateDecision struct {
	TargetMajor int          `json:"target_major"`
	Pass        bool         `json:"pass"`
	ExitCode    int          `json:"exit_code"`
	Reason      string       `json:"reason"`
	FailingApps []string     `json:"failing_apps"`
	Summary     CheckSummary `json:"summary"`
}
