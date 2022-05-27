package strategy

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"logconvert/internal/logcollecttool"
	"logconvert/internal/logcollecttool/util/templates"
	"logconvert/internal/pkg/soc"
	"logconvert/log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	cmdutil "logconvert/internal/logcollecttool/cmd/util"

	"logconvert/cli/genericclioptions"
	rootoptions "logconvert/internal/pkg/options"
)

type Options struct {
	EsNodes []string
	From    string
	To      string
	Limit   int
	genericclioptions.IOStreams
}

func NewOptions(ioStreams genericclioptions.IOStreams) *Options {
	return &Options{
		From:      "now-6h/h",
		To:        "now/h",
		Limit:     20,
		IOStreams: ioStreams,
	}
}

var validateExample = templates.Examples(`logcollecttool strategy`)

// NewCmdValidate returns new initialized instance of 'validate' sub command.
func NewCmdValidate(rootOptions *rootoptions.ServerRunOptions, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "strategy",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "get the longest-running strategys",
		TraverseChildren:      true,
		Long:                  "get the longest-running strategys",
		Example:               validateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(rootOptions, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.From, "from", o.From, `start with an anchor date, which can either be now, or a date string ending with ||.
		 This anchor date can optionally be followed by one or more maths expressions:
			+1h: Add one hour
			-1d: Subtract one day
			/d: Round down to the nearest day`)
	cmd.Flags().StringVar(&o.To, "to", o.To, `start with an anchor date, which can either be now, or a date string ending with ||.
	       This anchor date can optionally be followed by one or more maths expressions:
			+1h: Add one hour
			-1d: Subtract one day
			/d: Round down to the nearest day`)

	cmd.Flags().IntVar(&o.Limit, "limit", o.Limit, `return  top records which specify by option limit after sort desc`)

	return cmd
}

// Complete completes all the required options.
func (o *Options) Complete(rootOptions *rootoptions.ServerRunOptions, cmd *cobra.Command, args []string) error {
	isaGlobal := soc.NewConfig(rootOptions.IsaGlobalConfigPath)
	o.EsNodes = isaGlobal.Component.EsNodes
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *Options) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a validate sub command using the specified options.
func (o *Options) Run(args []string) error {
	var addresses []string

	for _, addr := range o.EsNodes {
		headAddr := "http://" + addr
		addresses = append(addresses, headAddr)
	}

	var data [][]string
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 1200,
			MaxConnsPerHost:       10,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client:%s", err)
	}

	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"gte": o.From,
								"lte": o.To,
							},
						},
					},
				},
			},
		},

		"aggs": map[string]interface{}{
			"strategy": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "parse_strategy",
					"size":  o.Limit,
					"order": map[string]interface{}{
						"cost_stats.avg": "desc",
					},
				},
				"aggs": map[string]interface{}{
					"cost_stats": map[string]interface{}{
						"stats": map[string]interface{}{
							"field": "parse_cost",
						},
					},
				},
			},
		},
	}

	requestBody, err := json.Marshal(query)
	if err != nil {
		log.Errorf("Marshal query error:%s", err)
	}

	indexes := "logstash-event*-" + time.Now().Format("20060102")
	req := esapi.SearchRequest{
		Index: []string{indexes},
		Body:  strings.NewReader(string(requestBody)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Errorf("Error getting response:%s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Errorf("[%s] Error searching,request:%s", res.Status(), requestBody)
	} else {
		var r Result
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Errorf("Error parsing the response body: %s", err)
		} else {
			for _, bucket := range r.Aggregations.Strategy.Buckets {

				var avgStr string
				avg := bucket.CostStats.Avg
				if avg > 1000000 {
					avgStr = color.RedString(strconv.FormatFloat(avg, 'f', 0, 64))
				} else {
					avgStr = color.GreenString(strconv.FormatFloat(avg, 'f', 0, 64))
				}

				data = append(data, []string{
					bucket.Key,
					strconv.FormatFloat(bucket.CostStats.Count, 'f', 0, 64),
					avgStr,
					strconv.FormatFloat(bucket.CostStats.Min, 'f', 0, 64),
					strconv.FormatFloat(bucket.CostStats.Max, 'f', 0, 64),
				},
				)
			}

		}
	}

	table := tablewriter.NewWriter(o.Out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(logcollecttool.TableWidth)
	table.SetHeader([]string{"name", "count", "avg(ns)", "min(ns)", "max(ns)"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return nil
}

type Stats struct {
	Count float64 `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
}

type Bucket struct {
	Key       string `json:"key"`
	CostStats Stats  `json:"cost_stats"`
}

type Aggregations struct {
	Strategy Strategy `json:"strategy"`
}

type Strategy struct {
	Buckets []Bucket `json:"buckets"`
}

type Result struct {
	Aggregations Aggregations `json:"aggregations"`
}
