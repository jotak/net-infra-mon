package jobs

import (
	"context"

	"github.com/jotak/net-infra-mon/pkg/jobs/vip"
)

var jobs = []func(ctx context.Context){
	vip.Run,
}

func Run(ctx context.Context) {
	// TODO: handle graceful exit
	for i := range jobs {
		go jobs[i](ctx)
	}
}
