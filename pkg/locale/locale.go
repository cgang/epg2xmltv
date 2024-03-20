package locale

import (
	"fmt"
	"time"
)

func MustGet(name string) *time.Location {
	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		return loc
	} else {
		panic(fmt.Sprintf("failed to load location: %s", err))
	}
}
