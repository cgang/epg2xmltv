package tvmao

import (
	"context"
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTVMao(t *testing.T) {
	items, err := getEpgInfo(context.TODO(), "CCTV1", time.Now())
	assert.NoError(t, err, "no error")

	xml.NewEncoder(os.Stdout).Encode(items)
}
