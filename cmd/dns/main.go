package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ralsnet/go-dns"
	"github.com/urfave/cli/v2"
)

const (
	OptConfig     = "config"
	OptSubdomains = "subdomains"
	OptRecursive  = "recursive"
	OptFormat     = "format"
)

func loadConfigFromFile(file string) *dns.Config {
	config := &dns.Config{}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return config
	}

	// Read file
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return config
	}
	defer f.Close()

	json.NewDecoder(f).Decode(config)

	return config
}

func action(c *cli.Context) error {

	// Load config

	config := loadConfigFromFile(c.String(OptConfig))

	for _, host := range strings.Split(c.String(OptSubdomains), ",") {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		if slices.Contains(config.Hosts, host) {
			continue
		}
		config.Hosts = append(config.Hosts, host)
	}

	config.Recursive = c.Bool(OptRecursive)

	// Lookup

	if c.Args().Len() == 0 {
		return errors.New("domain required")
	}

	domains := dns.Run(c.Args().First(), config)

	// Print results

	if c.String(OptFormat) == "json" {
		data, _ := json.MarshalIndent(domains, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, d := range domains {
			fmt.Println(d)
		}
	}

	return nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defaultConfigFile := home + "/.config/dns/dns.json"

	app := &cli.App{
		Name:  "dns",
		Usage: "DNS lookup tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    OptConfig,
				Aliases: []string{"c"},
				Usage:   "config file",
				Value:   defaultConfigFile,
			},
			&cli.StringFlag{
				Name:    OptSubdomains,
				Aliases: []string{"s"},
				Usage:   "subdomain hosts to lookup (comma separated) e.g. www,mail,ftp",
			},
			&cli.BoolFlag{
				Name:    OptRecursive,
				Aliases: []string{"r"},
				Usage:   "recursive lookup",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    OptFormat,
				Aliases: []string{"f"},
				Usage:   "output format",
				Value:   "text",
			},
		},
		Action: action,
	}

	app.Run(os.Args)
}
