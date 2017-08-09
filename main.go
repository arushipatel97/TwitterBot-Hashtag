package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	log               = &logger{logrus.New()}
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type dummyServer struct {
}

func (ds dummyServer) ServeHTTP(http.ResponseWriter, *http.Request) {
}

var first string
var dur time.Duration

func main() {
	app := cli.NewApp()
	app.Name = "twitterbot"
	app.Usage = "Twitterbot practice"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "twitter-consumer-key, k",
			Usage:       "twitter consumer key for twitter API",
			EnvVar:      "TWITTER_CONSUMER_KEY",
			Destination: &consumerKey,
		},
		cli.StringFlag{
			Name:        "twitter-consumer-secret, c",
			Usage:       "twitter consumer key for twitter API",
			EnvVar:      "TWITTER_CONSUMER_SECRET",
			Destination: &consumerSecret,
		},
		cli.StringFlag{
			Name:        "twitter-access-token, t",
			Usage:       "twitter access token for twitter API",
			EnvVar:      "TWITTER_ACCESS_TOKEN",
			Destination: &accessToken,
		},
		cli.StringFlag{
			Name:        "twitter-access-token-secret, s",
			Usage:       "twitter access token secret for twitter API",
			EnvVar:      "TWITTER_ACCESS_TOKEN_SECRET",
			Destination: &accessTokenSecret,
		},
		cli.StringFlag{
			Name:        "first-search, f",
			Usage:       "name of first hashtag to search",
			Value:       "#love",
			Destination: &first,
		},
	}
	app.Action = func(c *cli.Context) error {
		run(&first)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("error Running cli", err.Error())
	}
}

func run(first *string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	api.SetLogger(log)
	findHashtags(api, *first)

}

func findHashtags(api *anaconda.TwitterApi, first string) {
	startTime := time.Now()
	hashMap := make(map[string]int)
	dur = time.Second * 5
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{first},
	})

	defer stream.Stop()

	for time.Since(startTime) < dur {
		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Warningf("received unexpected value of type %T", v)
				continue
			}
			hashMap[t.Text]++
			fmt.Println(t.Text)
		}
	}
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
