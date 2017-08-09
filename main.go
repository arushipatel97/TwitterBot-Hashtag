package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
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
	durRound          = time.Second * 10
	durProgram        = time.Second * 50
)

var first string
var startProgram time.Time

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type hashTag struct {
	tag  string
	freq int
	next *hashTag
	prev *hashTag
}

var startList = &hashTag{
	tag:  "",
	freq: 0,
	next: nil,
	prev: nil,
}

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
			Value:       "#food",
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

	startProgram = time.Now()
	startList.tag = *first

	findHashtags(api, *first)
	fmt.Println("Final Order of Hashtags")
	printList()
}

func findHashtags(api *anaconda.TwitterApi, first string) {
	startRound := time.Now()
	hashMap := make(map[string]int)
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{first},
	})

	defer stream.Stop()

	for time.Since(startRound) < durRound {
		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Warningf("received unexpected value of type %T", v)
				continue
			}
			parseText(t.Text, hashMap)
			if time.Since(startProgram) > durProgram {
				return
			}
			if time.Since(startRound) > durRound {
				nextTag, freq := findMaxHashtag(hashMap, first)
				if nextTag == "" {
					nextTag = first
				}
				stream.Stop()
				addToList(nextTag, freq)
				findHashtags(api, nextTag)
			}
		}
	}
	return
}

func printList() {
	count := 1
	for temp := startList; temp != nil; temp = temp.next {
		fmt.Printf("%d.) %s with a frequency of %d \n", count, temp.tag, temp.freq)
		count++
	}
}

func addToList(text string, frequency int) {
	block := &hashTag{
		tag:  text,
		freq: frequency,
		next: nil,
	}
	var temp *hashTag
	for temp = startList; temp.next != nil; temp = temp.next {
	}
	temp.next = block
	block.prev = temp
	printList()
}

func parseText(text string, hashMap map[string]int) {
	parts := strings.Split(text, " ")
	for _, tag := range parts {
		if strings.HasPrefix(tag, "#") {
			hashMap[tag]++
		}
	}
}

func findMaxHashtag(hashMap map[string]int, first string) (string, int) {
	bestFreq := 0
	bestStr := ""
	for tag := range hashMap {
		if hashMap[tag] > bestFreq && !strings.EqualFold(tag, first) && notPrev(tag) {
			bestStr = tag
			bestFreq = hashMap[tag]
		}
	}
	return bestStr, bestFreq
}

func notPrev(tag string) bool {
	var temp *hashTag
	for temp = startList; temp.next != nil; temp = temp.next {
	}
	if temp.prev == nil {
		return true
	}
	return temp.prev.tag != tag
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
