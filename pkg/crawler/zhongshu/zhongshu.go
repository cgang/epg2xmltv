package zhongshu

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/cgang/epg2xmltv/pkg/locale"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

var (
	dataExpr   = regexp.MustCompile(`epgs\[\d+\]=new Array\(\"(\d+)\",\"(\d+)\",\"(\d+:\d+)\", \"(.+?)\",.+?\)`)
	CST        = locale.MustGet("Asia/Shanghai")
	httpClient = &http.Client{}
)

func getEpgInfo(ctx context.Context, id string, begin time.Time, week int) ([]xmltv.Programme, error) {
	urlStr := fmt.Sprintf("http://epg.tv.cn/epg/%s/live/index.php?week=%d",
		id, week)
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

	var items []xmltv.Programme
	year, month, _ := begin.Date()
	for _, m := range dataExpr.FindAllSubmatch(data, -1) {
		tm, err := time.ParseInLocation("2006 1 2 15:04", fmt.Sprintf("%d %s %s %s", year, m[1], m[2], m[3]), CST)
		if err != nil {
			log.Printf("Failed to parse date: %s", err)
			continue
		}

		if tm.Month() < month { // beginning of new year
			tm = tm.AddDate(1, 0, 0)
		}
		if tm.Before(begin) {
			continue
		}

		items = append(items, xmltv.Programme{
			Start: xmltv.Timestamp(tm),
			Title: xmltv.NewText("", string(m[4])),
		})
	}
	return items, nil
}

func GetProgram(ctx context.Context, arg string) (*xmltv.Program, error) {
	begin := time.Now().Add(-time.Hour)
	items, err := getEpgInfo(ctx, arg, begin, 0)
	if err != nil {
		return nil, err
	}

	program := xmltv.NewProgram()
	program.AddItems(items)

	if begin.Weekday() == time.Sunday {
		if items, err := getEpgInfo(ctx, arg, begin, 1); err == nil {
			program.AddItems(items)
		} else {
			log.Printf("Failed to get EPG for next week: %s", err)
		}
	}

	return program, nil
}
