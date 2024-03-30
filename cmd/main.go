package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cgang/epg2xmltv/pkg/config"
	"github.com/cgang/epg2xmltv/pkg/crawler"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

func main() {
	help := flag.Bool("h", false, "help")
	quiet := flag.Bool("q", false, "quiet")
	conf := flag.String("c", "", "config YAML")

	flag.Usage = func() {
		prompt := "Usage: %s -c config.yaml [-h]\n" +
			"A program to crawl EPG from website and save it to XMLTV file(s)\n\n" +
			"Options:\n"
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, prompt, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(out)
	}

	flag.Parse()

	if *help || *conf == "" {
		flag.Usage()
		return
	}

	if *quiet {
		// suppress log in quiet mode
		log.SetOutput(io.Discard)
	}

	cfgs, err := config.LoadConfigs(*conf)
	if err != nil {
		fmt.Printf("Failed to load %s: %s\n", *conf, err)
		os.Exit(1)
	}

	ctx := context.Background()
	for _, xcfg := range cfgs {
		tv := xmltv.NewXml()
		for _, channel := range xcfg.Channels {
			source := channel.Source
			if program, err := crawler.Run(ctx, source); err == nil {
				tv.AddProgram(channel.Id, channel.Name, program)
			} else {
				log.Printf("Failed to get program for %s: %s", source, err)
			}
		}

		if err = tv.Save(xcfg.Name); err != nil {
			fmt.Printf("Failed to save %s: %s\n", xcfg.Name, err)
		}
	}
}
