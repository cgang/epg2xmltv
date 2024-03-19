package cntv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/cgang/epg2xmltv/pkg/config"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

const (
	oneDay = 24 * time.Hour
)

var (
	httpClient = &http.Client{}
	jsonpExpr  = regexp.MustCompile(`[a-zA-Z]+\((.*)\);`)
)

type Timestamp int

func (t Timestamp) toTime() time.Time {
	return time.Unix(int64(t), 0)
}

type ProgramItem struct {
	StartTime Timestamp `json:"startTime"` // 1710777900
	EndTime   Timestamp `json:"endTime"`   // 1710780090
	Length    int       `json:"length"`    // 2190
	ShowTime  string    `json:"showTime"`  // 00:05
	Title     string    `json:"title"`     //
}

func (i ProgramItem) toProgramme() xmltv.Programme {
	return xmltv.NewProgramme(i.StartTime.toTime(), i.EndTime.toTime(), i.Title)
}

type ProgramGuide struct {
	ChannelName string        `json:"channelName"`
	List        []ProgramItem `json:"list"`
}

func (g *ProgramGuide) toProgrammes() []xmltv.Programme {
	var result []xmltv.Programme
	for _, item := range g.List {
		result = append(result, item.toProgramme())
	}
	return result
}

func getEpgInfo(ctx context.Context, arg string, dt time.Time) (*ProgramGuide, error) {
	urlStr := fmt.Sprintf("http://api.cntv.cn/epg/getEpgInfoByChannelNew?c=%s&serviceId=tvcctv&d=%s&t=jsonp&cb=set",
		arg, dt.Format("20060102"))
	log.Printf("URL: %s", urlStr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, data)
	}

	if m := jsonpExpr.FindSubmatch(data); m != nil {
		var data = struct {
			Data map[string]*ProgramGuide `json:"data"`
		}{
			Data: make(map[string]*ProgramGuide),
		}

		if err = json.Unmarshal(m[1], &data); err == nil {
			return data.Data[arg], nil
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("JSONP not found")
	}
}

func GetProgram(ctx context.Context, cfg config.CrawlerConfig) (*xmltv.Program, error) {
	arg := cfg.ArgOrId()
	dt := time.Now()
	current, err := getEpgInfo(ctx, arg, dt)
	if err != nil {
		return nil, err
	}

	program := xmltv.NewProgram(cfg.Id, xmltv.NewText("zh", cfg.Name))
	program.AddItems(current.toProgrammes())

	if next, err := getEpgInfo(ctx, arg, dt.Add(oneDay)); err == nil {
		program.AddItems(next.toProgrammes())
	} else {
		log.Printf("Failed to get EPG for tomorrow: %s", err)
	}

	return program, nil
}
