package resolver

import (
	"context"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

type Resolver interface {
	Resolve(ctx context.Context, targetMajor int, apps []types.AppInfo) ([]types.CheckResult, error)
}
