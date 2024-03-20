package brtv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cgang/epg2xmltv/pkg/locale"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

const (
	oneDay = 24 * time.Hour
)

var (
	CST        = locale.MustGet("Asia/Shanghai")
	httpClient = &http.Client{}
)

type LocalTime string

func toTime(dt time.Time, clock string) time.Time {
	year, month, day := dt.Date()
	if tm, err := time.Parse("15:04", clock); err == nil {
		h, m, s := tm.Clock()
		return time.Date(year, month, day, h, m, s, 0, CST)
	} else {
		return time.Time{}
	}
}

type ProgramItem struct {
	StartTime LocalTime `json:"startTime"` // clock 12:05
	EndTime   LocalTime `json:"endTime"`   // clock 12:30
	Title     string    `json:"name"`      //
}

func (i ProgramItem) toProgramme(dt time.Time) xmltv.Programme {
	start := toTime(dt, string(i.StartTime))
	stop := toTime(dt, string(i.EndTime))
	if stop.Before(start) {
		stop = stop.AddDate(0, 0, 1)
	}
	return xmltv.NewProgramme(start, stop, i.Title)
}

type ProgramGuide struct {
	Id        string        `json:"id"`
	Programes []ProgramItem `json:"programes"`
}

func (g *ProgramGuide) toProgrammes(dt time.Time) []xmltv.Programme {
	var result []xmltv.Programme
	for _, item := range g.Programes {
		result = append(result, item.toProgramme(dt))
	}
	return result
}

func getEpgInfo(ctx context.Context, id string, dt time.Time) (*ProgramGuide, error) {
	urlStr := fmt.Sprintf("https://dynamic.rbc.cn/bvradio_app/service/LIVE?functionName=getCurrentChannel&channelId=%s&curdate=%s",
		id, dt.Format("2006-01-02"))
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

	var body = struct {
		Channel ProgramGuide `json:"channel"`
	}{}

	if err = json.Unmarshal(data, &body); err == nil {
		return &body.Channel, nil
	} else {
		return nil, err
	}
}

func GetProgram(ctx context.Context, arg string) (*xmltv.Program, error) {
	dt := time.Now()
	current, err := getEpgInfo(ctx, arg, dt)
	if err != nil {
		return nil, err
	}

	program := xmltv.NewProgram()
	program.AddItems(current.toProgrammes(dt))

	dt = dt.Add(oneDay)
	if next, err := getEpgInfo(ctx, arg, dt); err == nil {
		program.AddItems(next.toProgrammes(dt))
	} else {
		log.Printf("Failed to get EPG for tomorrow: %s", err)
	}

	return program, nil
}
