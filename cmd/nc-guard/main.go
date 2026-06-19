package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/collector"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/engine"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/policy"
	"github.com/Rednox/nextcloud-upgrade-guard/internal/reporter"
	localresolver "github.com/Rednox/nextcloud-upgrade-guard/internal/resolver/local"
)

const (
	defaultOccPath = "/var/www/nextcloud/occ"
	defaultPHPBin  = "php"
)

func main() {
	code := run(os.Args[1:])
	os.Exit(code)
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 1
	}

	switch args[0] {
	case "inspect":
		if err := runInspect(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "inspect failed: %v\n", err)
			return 1
		}
		return 0
	case "check":
		if err := runCheck(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "check failed: %v\n", err)
			return 1
		}
		return 0
	case "gate":
		code, err := runGate(args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "gate failed: %v\n", err)
			return 20
		}
		return code
	default:
		usage()
		return 1
	}
}

func runInspect(args []string) error {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	occPath := fs.String("occ-path", defaultOccPath, "Path to Nextcloud occ executable")
	phpBin := fs.String("php-bin", defaultPHPBin, "PHP binary")
	format := fs.String("format", "table", "Output format: table|json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	apps, err := collector.NewOCCCollector(*occPath, *phpBin).Collect(context.Background())
	if err != nil {
		return err
	}

	return reporter.New().RenderInspect(os.Stdout, apps, *format)
}

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	targetMajor := fs.Int("target", 0, "Target Nextcloud major version")
	occPath := fs.String("occ-path", defaultOccPath, "Path to Nextcloud occ executable")
	phpBin := fs.String("php-bin", defaultPHPBin, "PHP binary")
	format := fs.String("format", "table", "Output format: table|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *targetMajor <= 0 {
		return errors.New("--target is required and must be > 0")
	}

	results, err := engine.Check(context.Background(), collector.NewOCCCollector(*occPath, *phpBin), localresolver.New(), *targetMajor)
	if err != nil {
		return err
	}

	report := reporter.BuildCheckReport(*targetMajor, results)
	return reporter.New().RenderCheck(os.Stdout, report, *format)
}

func runGate(args []string) (int, error) {
	fs := flag.NewFlagSet("gate", flag.ContinueOnError)
	targetMajor := fs.Int("target", 0, "Target Nextcloud major version")
	policyPath := fs.String("policy", "", "Path to policy YAML file")
	occPath := fs.String("occ-path", defaultOccPath, "Path to Nextcloud occ executable")
	phpBin := fs.String("php-bin", defaultPHPBin, "PHP binary")
	if err := fs.Parse(args); err != nil {
		return 20, err
	}
	if *targetMajor <= 0 {
		return 20, errors.New("--target is required and must be > 0")
	}
	if *policyPath == "" {
		return 20, errors.New("--policy is required")
	}

	cfg, err := policy.Load(*policyPath)
	if err != nil {
		return 20, err
	}

	results, err := engine.Check(context.Background(), collector.NewOCCCollector(*occPath, *phpBin), localresolver.New(), *targetMajor)
	if err != nil {
		return 20, err
	}

	decision := policy.PolicyEvaluator{}.Evaluate(cfg, *targetMajor, results)
	if err := reporter.New().RenderGate(os.Stdout, decision); err != nil {
		return 20, err
	}
	return decision.ExitCode, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "nc-guard: Nextcloud upgrade compatibility guard")
	fmt.Fprintln(os.Stderr, "Usage: nc-guard <inspect|check|gate> [flags]")
}
