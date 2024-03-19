package crawler

import (
	"context"
	"fmt"

	"github.com/cgang/epg2xmltv/pkg/config"
	"github.com/cgang/epg2xmltv/pkg/crawler/brtv"
	"github.com/cgang/epg2xmltv/pkg/crawler/cntv"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

func Run(ctx context.Context, cfg config.CrawlerConfig) (*xmltv.Program, error) {
	switch cfg.Type {
	case "cntv":
		return cntv.GetProgram(ctx, cfg.Id, cfg.Arg)
	case "brtv":
		return brtv.GetProgram(ctx, cfg.Id, cfg.Arg)
	default:
		return nil, fmt.Errorf("unknown type: %s", cfg.Type)
	}
}
