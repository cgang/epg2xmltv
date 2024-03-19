package brtv

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetBRTV(t *testing.T) {
	epg, err := getEpgInfo(context.TODO(), "134", time.Now())
	assert.NoError(t, err, "no error")

	result, err := json.MarshalIndent(epg, "", "  ")
	assert.NoError(t, err, "JSON")
	fmt.Println(string(result))
}
