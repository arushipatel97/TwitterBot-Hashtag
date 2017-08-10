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
	log               = &Logger{logrus.New()}
	durRound          = time.Second * 60
	durProgram        = time.Minute * 2
)

var first string
var careAboutPrev bool
var startProgram time.Time

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

//head of linked list
var startList = &HashTag{
	tag:  "",
	freq: 0,
	next: nil,
	prev: nil,
}

//CLI for running without .env file or specifying a different starting word
//(for Command line)
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
		cli.BoolFlag{
			Name:        "not-prev, n",
			Usage:       "if you don't want to search any previously searched hashtag",
			Destination: &careAboutPrev,
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

//sets up API environment & makes first call to finding initial hashtag
func run(first *string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(log)

	//stops running all searches after specified time
	startProgram = time.Now()
	startList.tag = *first
	careAboutPrev = true

	findHashtags(api, *first)
	fmt.Println("Final Order of Hashtags")
	PrintList(*first)
}

//searches for a specific hashtag in Twitter (real-time), and makes a map for
//the other hashtags mentioned in the posts containing the specified hashtag &
//based on the time has passed
func findHashtags(api *anaconda.TwitterApi, first string) {
	startRound := time.Now()
	hashMap := make(map[string]int)
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{first}, //hashtag that is being searched for
	})

	defer stream.Stop()

	for time.Since(startProgram) < durProgram {
		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Warningf("received unexpected value of type %T", v)
				continue //don't want to panic, just take note of
			}
			parseText(t.Text, hashMap)
			if time.Since(startProgram) > durProgram {
				return
			}
			if time.Since(startRound) > durRound {
				roundCheck(hashMap, stream, api)
			}
		}
		if time.Since(startProgram) < durProgram {
			time.Sleep(durRound)
			findHashtags(api, first)
		} else {
			return
		}
	}
	return
}

//checks time the round of searching for a specific hashtag has been running
func roundCheck(hashMap map[string]int, stream *anaconda.Stream, api *anaconda.TwitterApi) {
	//time to move to next hashtag
	nextTag, freq, total := findMaxHashtag(hashMap, first)
	stream.Stop()
	AddToList(nextTag, freq, total) //add most frequent hashtag to linked list
	PrintList(first)
	findHashtags(api, nextTag) //recursively calls itself with next hashtag

}

//parses tweets found with given hashtag to find other hashtags mentioned, &
//places them into a map
func parseText(text string, hashMap map[string]int) {
	parts := strings.Split(text, " ")
	for _, tag := range parts {
		if strings.HasPrefix(tag, "#") { //only care about hashtags
			hashMap[tag]++
		}
	}
}

//goes through current hashMap finding the hashtag with the greatest frequency
//of showing up in posts with the specified hashtag, returning both the most
//frequent hashtag itself, along with its count/frequency
func findMaxHashtag(hashMap map[string]int, first string) (string, int, int) {
	bestFreq := 0
	bestStr := ""
	count := 0
	for tag := range hashMap {
		count += hashMap[tag]
		if hashMap[tag] > bestFreq && !strings.EqualFold(tag, first) && !matchesPrev(tag) {
			bestStr = tag
			bestFreq = hashMap[tag]
		}
	}
	if bestStr == "" {
		bestStr = first
		bestFreq = hashMap[bestStr]
	}
	return bestStr, bestFreq, count
}

//checks if the next hashtag to be searched matches any previous hashtag
//that was searched -- to prevent just going back & forth
//(this is an optional parameter, based on careAboutPrev)
func matchesPrev(tag string) bool {
	if !careAboutPrev {
		return false
	}
	return InList(tag)
}
