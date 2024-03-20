package zhongshu

import (
	"context"
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZhongshu(t *testing.T) {
	items, err := getEpgInfo(context.TODO(), "cetv1", time.Now(), 0)
	assert.NoError(t, err, "no error")

	xml.NewEncoder(os.Stdout).Encode(items)
}
