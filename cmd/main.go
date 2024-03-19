package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cgang/epg2xmltv/pkg/config"
	"github.com/cgang/epg2xmltv/pkg/crawler"
	"github.com/cgang/epg2xmltv/pkg/xmltv"
)

func main() {
	help := flag.Bool("h", false, "help")
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

	cfg, err := config.LoadConfig(*conf)
	if err != nil {
		fmt.Printf("Failed to load %s: %s\n", *conf, err)
		os.Exit(1)
	}

	ctx := context.Background()
	programs := make(map[string]*xmltv.Program)
	for _, ccfg := range cfg.CrawlersConfig {
		if prog, err := crawler.Run(ctx, ccfg); err == nil {
			programs[prog.Channel.Id] = prog
		} else {
			log.Printf("Failed to get program for %s: %s", ccfg.Id, err)
		}
	}

	for _, xcfg := range cfg.OutputsConfig {
		tv := xmltv.NewXml(xcfg.Channels, programs)
		if err = tv.Save(xcfg.Name); err != nil {
			fmt.Printf("Failed to save %s: %s\n", xcfg.Name, err)
		}
	}
}
