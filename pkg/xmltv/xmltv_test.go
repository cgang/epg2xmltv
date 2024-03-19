package xmltv

import (
	"encoding/xml"
	"os"
	"testing"
	"time"
)

func TestXml(t *testing.T) {
	var tv XmlTv
	tv.Channels = append(tv.Channels, Channel{
		Id: "id1",
	})
	tv.Programmes = append(tv.Programmes, Programme{
		Start:   Timestamp(time.Now()),
		Stop:    Timestamp(time.Now()),
		Channel: "channel",
	})
	xml.NewEncoder(os.Stdout).Encode(&tv)
}
