package collector

import (
	"context"

	"github.com/Rednox/nextcloud-upgrade-guard/internal/types"
)

type Collector interface {
	Collect(ctx context.Context) ([]types.AppInfo, error)
}
