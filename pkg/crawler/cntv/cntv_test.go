package cntv

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCNTV(t *testing.T) {
	epg, err := getEpgInfo(context.TODO(), "cctv1", time.Now())
	assert.NoError(t, err, "no error")

	result, err := json.MarshalIndent(epg, "", "  ")
	assert.NoError(t, err, "JSON")
	fmt.Println(string(result))
}
