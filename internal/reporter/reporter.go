package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/engine"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

type Reporter interface {
	RenderInspect(w io.Writer, apps []types.AppInfo, format string) error
	RenderCheck(w io.Writer, report types.CheckReport, format string) error
	RenderGate(w io.Writer, decision types.GateDecision) error
}

type ConsoleReporter struct{}

func New() *ConsoleReporter { return &ConsoleReporter{} }

func BuildCheckReport(targetMajor int, results []types.CheckResult) types.CheckReport {
	return BuildCheckReportAt(targetMajor, time.Now().UTC(), results)
}

func BuildCheckReportAt(targetMajor int, generatedAt time.Time, results []types.CheckResult) types.CheckReport {
	copied := make([]types.CheckResult, len(results))
	copy(copied, results)
	sort.Slice(copied, func(i, j int) bool { return copied[i].AppID < copied[j].AppID })
	return types.CheckReport{
		TargetMajor: targetMajor,
		GeneratedAt: generatedAt,
		Apps:        copied,
		Summary:     engine.SummarizeResults(copied),
	}
}

func (ConsoleReporter) RenderInspect(w io.Writer, apps []types.AppInfo, format string) error {
	sorted := make([]types.AppInfo, len(apps))
	copy(sorted, apps)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Enabled != sorted[j].Enabled {
			return sorted[i].Enabled
		}
		return sorted[i].AppID < sorted[j].AppID
	})

	switch format {
	case "json":
		return renderJSON(w, sorted)
	case "table":
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "APP ID\tENABLED\tINSTALLED VERSION\tSOURCE")
		for _, app := range sorted {
			fmt.Fprintf(tw, "%s\t%t\t%s\t%s\n", app.AppID, app.Enabled, blankIfEmpty(app.InstalledVersion), app.Source)
		}
		return tw.Flush()
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func (ConsoleReporter) RenderCheck(w io.Writer, report types.CheckReport, format string) error {
	switch format {
	case "json":
		return renderJSON(w, report)
	case "table":
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "Target major:\t%d\n", report.TargetMajor)
		fmt.Fprintf(tw, "Totals:\tcompatible=%d upgradable_to_compatible=%d incompatible=%d unknown=%d\n\n", report.Summary.Compatible, report.Summary.UpgradableToCompatible, report.Summary.Incompatible, report.Summary.Unknown)
		fmt.Fprintln(tw, "APP ID\tSTATUS\tVERSION\tRECOMMENDED ACTION")
		for _, app := range report.Apps {
			if app.Status == types.StatusCompatible {
				continue
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", app.AppID, app.Status, blankIfEmpty(app.InstalledVersion), app.RecommendedAction)
		}
		return tw.Flush()
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func (ConsoleReporter) RenderGate(w io.Writer, decision types.GateDecision) error {
	status := "PASS"
	if !decision.Pass {
		status = "BLOCK"
	}
	fmt.Fprintf(w, "%s target=%d exit_code=%d reason=%s\n", status, decision.TargetMajor, decision.ExitCode, decision.Reason)
	if len(decision.FailingApps) > 0 {
		fmt.Fprintf(w, "Failing apps: %s\n", strings.Join(decision.FailingApps, ", "))
	}
	return nil
}

func renderJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func blankIfEmpty(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
