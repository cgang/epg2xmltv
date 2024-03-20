package crawler

import (
	"context"
	"fmt"
	"strings"

	"github.com/cgang/epg2xmltv/pkg/crawler/brtv"
	"github.com/cgang/epg2xmltv/pkg/crawler/cntv"
	"github.com/cgang/epg2xmltv/pkg/crawler/zhongshu"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

var (
	cache = make(map[string]*xmltv.Program)
)

func Run(ctx context.Context, source string) (program *xmltv.Program, err error) {
	if cached, ok := cache[source]; ok {
		return cached, nil
	}

	crawler, arg, found := strings.Cut(source, "/")
	if !found {
		return nil, fmt.Errorf("invalid source: %s", source)
	}

	switch crawler {
	case "cntv":
		program, err = cntv.GetProgram(ctx, arg)
	case "brtv":
		program, err = brtv.GetProgram(ctx, arg)
	case "zhongshu":
		program, err = zhongshu.GetProgram(ctx, arg)
	default:
		return nil, fmt.Errorf("unknown crawler: %s", crawler)
	}

	if err == nil {
		cache[source] = program
	}
	return
}
