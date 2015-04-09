package main

import (
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	cwl "github.com/awslabs/aws-sdk-go/service/cloudwatchlogs"
	"github.com/codegangsta/cli"
	. "github.com/tj/go-debug"
	"os"
	"strconv"
	"strings"
	"time"
)

type filterFn func(out *cwl.OutputLogEvent) bool
type formatter interface {
	Format(string, []string) string
}

var nop = func(out *cwl.OutputLogEvent) bool { return true }

const timeFormat = "2006-01-02 15:04:05"

func main() {
	app := cli.NewApp()
	app.Name = "cw-search"
	app.Usage = "cw-search [log-group] [log-stream]"
	app.Version = "1.0.0"

	now := time.Now().UTC()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region, r",
			Value: "us-east-1",
			Usage: "AWS region you need to hit",
		},
		cli.StringFlag{
			Name: "start, s",
			Value: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				0,
				0,
				0,
				0,
				time.FixedZone("UTC", 0),
			).Format(timeFormat),
			Usage: "start time for the search",
		},
		cli.StringFlag{
			Name: "end, e",
			Value: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				23,
				59,
				59,
				0,
				time.FixedZone("UTC", 0),
			).Format(timeFormat),
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "log line format. valid values are: json",
		},
		cli.StringFlag{
			Name:  "fields",
			Usage: "used with --format; filters the fields shown in the output",
		},
	}

	app.Action = func(c *cli.Context) {
		tail(c, nop)
	}

	app.Run(os.Args)
}

func getCwl(c *cli.Context) *cwl.CloudWatchLogs {
	cfg := &aws.Config{
		Credentials: aws.DetectCreds(
			"",
			"",
			"",
		),
		Region: c.String("region"),
	}

	cl := cwl.New(cfg)

	return cl
}

func tail(c *cli.Context, filter filterFn) {
	debug := Debug("tail")

	cl := getCwl(c)

	startTime := parseTime(c.String("start"))
	endTime := parseTime(c.String("end"))
	logGroup := c.Args().Get(0)
	logStream := c.Args().Get(1)

	infoOut(
		"running with start time", startTime.Format(time.Stamp),
		"and end time", endTime.Format(time.Stamp),
	)

	debug("start time is", strconv.FormatInt(startTime.Unix(), 10))
	debug("end time is", strconv.FormatInt(endTime.Unix(), 10))
	// all these aws.<Type> methods are used
	// because this struct wants pointers to everything
	// not entirely sure why! but it does!
	i := &cwl.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		StartTime:     aws.Long(startTime.Unix() * 1000),
		EndTime:       aws.Long(endTime.Unix() * 1000),
		StartFromHead: aws.Boolean(true),
	}

	out, err := cl.GetLogEvents(i)
	debug("got first set of events")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var f formatter
	fields := strings.Split(c.String("fields"), ",")
	switch c.String("format") {
	case "json":
		f = newJsonFormatter()
	}

	var lastFwdToken *string

	for {
		for _, v := range out.Events {
			if filter(v) == true {
				if f != nil {
					fmt.Println(f.Format(*v.Message, fields))
				} else {
					fmt.Println(*v.Message)
				}

			} else {
				debug("filter failed")
			}
		}
		debug("printed stuff")

		if lastFwdToken != nil && *out.NextForwardToken == *lastFwdToken {
			debug("done with stream")
			break
		}

		debug("have forward token, going to loop again")

		i.NextToken = out.NextForwardToken
		lastFwdToken = i.NextToken

		out, err = cl.GetLogEvents(i)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func parseTime(in string) time.Time {
	debug := Debug("timeParser")
	t, err := time.Parse(timeFormat, in)
	var res time.Time
	if err != nil {
		debug("time parsing failed w/input", in)
		res = time.Now().UTC()
		return res
	}

	res = t.UTC()

	return res
}

// prints stuff to stderror
func infoOut(stuff ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, stuff...)
}
