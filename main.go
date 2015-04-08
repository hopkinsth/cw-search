package main

import (
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	cwl "github.com/awslabs/aws-sdk-go/service/cloudwatchlogs"
	"github.com/codegangsta/cli"
	. "github.com/tj/go-debug"
	"os"
	"time"
)

type filterFn func(out *cwl.OutputLogEvent) bool

var nop = func(out *cwl.OutputLogEvent) bool { return true }

func main() {
	app := cli.NewApp()
	app.Name = "cw-search"
	app.Usage = "hey there"

	now := time.Now().UTC()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "logGroup, lg",
			Value: "",
			Usage: "log group you're searching through",
		},
		cli.StringFlag{
			Name:  "logStream, ls",
			Value: "",
			Usage: "log stream you're searching",
		},
		cli.StringFlag{
			Name:  "region, rg",
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
			).Format(time.Stamp),
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
			).Format(time.Stamp),
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

	infoOut(
		"running with start time", startTime.Format(time.Stamp),
		"and end time", endTime.Format(time.Stamp),
	)

	// all these aws.<Type> methods are used
	// because this struct wants pointers to everything
	// not entirely sure why! but it does!
	i := &cwl.GetLogEventsInput{
		LogGroupName:  aws.String("logGroup"),
		LogStreamName: aws.String("logStream"),
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

	next := out.NextForwardToken

	for next != nil {
		debug("have forward token, going to loop again")
		for _, v := range out.Events {
			if filter(v) {
				fmt.Println(*v.Message)
			}
		}
		debug("printed stuff")
		i.NextToken = next
		out, err = cl.GetLogEvents(i)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func parseTime(in string) time.Time {
	t, err := time.Parse(time.Stamp, in)
	var res time.Time
	if err != nil {
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
