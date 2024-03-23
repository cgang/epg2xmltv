package tvmao

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

var (
	CST        = locale.MustGet("Asia/Shanghai")
	httpClient = &http.Client{}
)

type ProgramItem struct {
	Name string `json:"name"` // "纪录片瞬间中国-农耕探文明3"
	Time string `json:"time"` // "00:06"
}

type ProgramGuide struct {
	EpgName string        `json:"epgName"`
	EpgCode string        `json:"epgCode"`
	Items   []ProgramItem `json:"pro"`
}

func getEpgInfo(ctx context.Context, id string, dt time.Time) ([]xmltv.Programme, error) {
	urlStr := fmt.Sprintf("https://lighttv.tvmao.com/qa/qachannelschedule?epgCode=%s&op=getProgramByChnid&epgName=&isNew=on&day=%d",
		id, dt.Weekday())
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
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %s", resp.Status, data)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []json.RawMessage
	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, err
	} else if len(data) < 3 {
		log.Println(string(raw))
		return nil, fmt.Errorf("unrecognized payload")
	}

	var epg ProgramGuide
	if err = json.Unmarshal(data[2], &epg); err != nil {
		return nil, err
	}

	year, month, day := dt.Date()
	var items []xmltv.Programme
	for _, item := range epg.Items {
		tm, err := time.ParseInLocation("2006 1 2 15:04", fmt.Sprintf("%d %d %d %s", year, month, day, item.Time), CST)
		if err != nil {
			log.Printf("Failed to parse date: %s", err)
			continue
		}

		items = append(items, xmltv.Programme{
			Start: xmltv.Timestamp(tm),
			Stop:  xmltv.Timestamp(tm.Add(5 * time.Minute)), // placeholder
			Title: xmltv.NewText("", item.Name),
		})
	}

	// reset stop time
	for idx, item := range items {
		if idx < len(items)-1 {
			item.Stop = items[idx+1].Start
			items[idx] = item
		}
	}

	return items, nil
}

func GetProgram(ctx context.Context, arg string) (*xmltv.Program, error) {
	today := time.Now()
	items, err := getEpgInfo(ctx, arg, today)
	if err != nil {
		return nil, err
	}

	program := xmltv.NewProgram()
	program.AddItems(items)

	if items, err := getEpgInfo(ctx, arg, today.Add(24*time.Hour)); err == nil {
		program.AddItems(items)
	} else {
		log.Printf("Failed to get EPG for tomorrow: %s", err)
	}

	return program, nil
}
