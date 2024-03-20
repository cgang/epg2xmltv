package xmltv

import (
	"encoding/xml"
	"os"
	"time"
)

const (
	DocType    = "<!DOCTYPE tv SYSTEM \"xmltv.dtd\">\n"
	TimeLayout = "20060102150405 -0700"
)

type LocalizedText struct {
	Lang string `xml:"lang,attr,omitempty"`
	Data string `xml:",chardata"`
}

func NewText(lang, data string) LocalizedText {
	return LocalizedText{Lang: lang, Data: data}
}

type ChannelIcon struct {
	Src    string `xml:"src,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type Channel struct {
	Id           string          `xml:"id,attr"`
	DisplayNames []LocalizedText `xml:"display-name"`
	Icons        []ChannelIcon   `xml:"icon"`
}

type Timestamp time.Time

func (ts Timestamp) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	tm := time.Time(ts)
	if tm.IsZero() {
		return xml.Attr{}, nil
	}
	value := tm.Format(TimeLayout)
	return xml.Attr{Name: name, Value: value}, nil
}

type Programme struct {
	Start        Timestamp       `xml:"start,attr"`
	Stop         Timestamp       `xml:"stop,attr"`
	Channel      string          `xml:"channel,attr"`
	Title        LocalizedText   `xml:"title"`
	Descriptions []LocalizedText `xml:"desc,omitempty"`
}

func NewProgramme(start, stop time.Time, title string) Programme {
	return Programme{
		Start: Timestamp(start),
		Stop:  Timestamp(stop),
		Title: LocalizedText{Data: title},
	}
}

type Program struct {
	ChannelNames []LocalizedText // channel name from crawler
	Items        []Programme
}

func NewProgram(names ...LocalizedText) *Program {
	return &Program{ChannelNames: names}
}

func (p *Program) AddItems(items []Programme) {
	p.Items = append(p.Items, items...)
}

type XmlTv struct {
	XMLName           string      `xml:"tv"`
	SourceInfoUrl     string      `xml:"source-info-url,attr,omitempty"`
	SourceInfoName    string      `xml:"source-info-name,attr,omitempty"`
	GeneratorInfoName string      `xml:"generator-info-name,attr,omitempty"`
	GeneratorInfoUrl  string      `xml:"generator-info-url,attr,omitempty"`
	Channels          []Channel   `xml:"channel"`
	Programmes        []Programme `xml:"programme"`
}

func NewXml() *XmlTv {
	return &XmlTv{
		GeneratorInfoName: "epg2xmltv",
		GeneratorInfoUrl:  "https://github.com/cgang/epg2xmltv",
	}
}

func (t *XmlTv) Save(name string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write([]byte(xml.Header))
	f.Write([]byte(DocType))

	return xml.NewEncoder(f).Encode(t)
}

func (t *XmlTv) AddProgram(id, name string, program *Program) {
	names := program.ChannelNames
	if len(names) == 0 {
		names = append(names, NewText("", name))
	}

	t.Channels = append(t.Channels, Channel{Id: id, DisplayNames: names})
	for _, item := range program.Items {
		item.Channel = id
		t.Programmes = append(t.Programmes, item)
	}
}
