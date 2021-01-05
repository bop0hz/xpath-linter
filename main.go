package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/antchfx/xmlquery"
	"github.com/bop0hz/xpath-linter/lint"
	"github.com/bop0hz/xpath-linter/report"
	"github.com/bop0hz/xpath-linter/version"
	"gopkg.in/yaml.v2"
)

var (
	targets, condition, query string
	fileName, cfg             string
	failed, ci, showVer, must bool
	app                       = version.Version()
)

type Config struct {
	Rules []*lint.Rule `yaml: "rules"`
}

func init() {
	flag.StringVar(&targets, "targets", "/", "XPath query to find target nodes, default is root node")
	flag.StringVar(&condition, "having", "", "Optional XPath condition which must be true to proceed with a query. Works like a filter")
	flag.BoolVar(&must, "must", true, "Whether target node must have or not 'contain' nodes")
	flag.StringVar(&query, "contain", "/", "XPath query to evaluate. Cheatsheet https://devhints.io/xpath")
	flag.BoolVar(&ci, "ci", false, "Continous integration mode. Evalutates rules from the config against the target file")
	flag.StringVar(&cfg, "cfg", "config.yaml", "File with linting rules")
	flag.BoolVar(&showVer, "version", false, "Show app version")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS] [FILE] \n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
	}
	if showVer {
		fmt.Fprint(os.Stdout, app.Print())
		os.Exit(0)
	}

	fileName = flag.Args()[0]
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "File not found: %s\n", fileName)
		os.Exit(1)
	}
	defer f.Close()

	reporter := &report.ColorReport{}
	rules := []*lint.Rule{}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		fmt.Printf("%s %s\n", fileName, reporter.Compile("XML validator", 0, err.Error(), ""))
		os.Exit(3)
	}

	if ci {
		data, err := ioutil.ReadFile(cfg)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Could not open config file %s\n", cfg)
		}
		config := Config{}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Could not load config: %s", err)
			os.Exit(3)
		}
		rules = config.Rules
	} else {
		// Ad-hoc query
		rule := lint.Rule{Targets: targets, Having: condition, Must: must, Contain: query}
		name := fmt.Sprintf("Targets: %s, having: %s, must: %v, contain: %s",
			rule.Targets, rule.Having, rule.Must, rule.Contain)
		rule.Name = name
		rules = append(rules, &rule)
	}
	results, err := lint.CheckAll(doc, rules, reporter)
	if err != nil {
		fmt.Printf("Something went wrong during rules evaluation: %s", err)
		os.Exit(3)
	}

	for _, result := range results {
		fmt.Printf("%s %s\n", fileName, result)
	}
}
