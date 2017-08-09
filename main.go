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
var careAboutPrev bool
var startProgram time.Time

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

//block type for linked list
type hashTag struct {
	tag   string
	freq  int
	total int
	next  *hashTag
	prev  *hashTag
}

//head of linked list
var startList = &hashTag{
	tag:  "",
	freq: 0,
	next: nil,
	prev: nil,
}

//CLI for running without .env file or specifying a different starting word
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
			Usage:       "if you don't want it to search previous hashtag",
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
	printList(*first)
}

//searches for a specific hashtag in Twitter (real-time), and makes a map for
//the other hashtags mentioned in the posts containing the specified hashtag &
//based on the time has passed, will recursively call itself with the new
//specified hashtag being the one with the greatest frequency in previous map
func findHashtags(api *anaconda.TwitterApi, first string) {
	startRound := time.Now()
	hashMap := make(map[string]int)
	stream := api.PublicStreamFilter(url.Values{
		"track": []string{first}, //hashtag that is being searched for
	})
	defer stream.Stop()

	for true {
		//fixes time spent looking for each hashtag
		for v := range stream.C {
			t, ok := v.(anaconda.Tweet)
			if !ok {
				log.Warningf("received unexpected value of type %T", v)
				continue //don't want to panic, just take note of
			}
			parseText(t.Text, hashMap)
			if time.Since(startProgram) > durProgram {
				//fixes approximate total time program spends searching for hashtags
				return
			}
			if time.Since(startRound) > durRound {
				//time to move to next hashtag
				nextTag, freq, total := findMaxHashtag(hashMap, first)
				if nextTag == "" {
					nextTag = first
				}
				stream.Stop()
				addToList(nextTag, freq, total) //add most frequent hashtag to linked list
				printList(first)
				findHashtags(api, nextTag) //recursively calls itself with next hashtag
			}
		}
	}
	return
}

//goes through linked list printing the most popular hashtags/order of searching
//with frequency
func printList(first string) {
	count := 1
	grammar1, grammar2 := "tweets", "tweets"
	for temp := startList; temp != nil; temp = temp.next {
		if temp.prev != nil {
			prev := temp.prev.tag
			if temp.freq == 1 {
				grammar1 = "tweet"
			}
			if temp.total == 1 {
				grammar2 = "tweet"
			}
			fmt.Printf("%d.) %d %s had %s of the %d %s that had %s \n", count, temp.freq, grammar1, temp.tag, temp.total, grammar2, prev)
			count++
		}
	}
}

//adds next hashtag to be searched in linked list
func addToList(text string, frequency int, total int) {
	block := &hashTag{
		tag:   text,
		freq:  frequency,
		total: total,
		next:  nil,
		prev:  nil,
	}
	var temp *hashTag
	for temp = startList; temp.next != nil; temp = temp.next {
	}
	temp.next = block
	block.prev = temp
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
		if hashMap[tag] > bestFreq && !strings.EqualFold(tag, first) && !prev(tag) {
			bestStr = tag
			bestFreq = hashMap[tag]
		}
	}
	return bestStr, bestFreq, count
}

//makes sure that the next hashtag to be searched isn't any previous hashtag
//that was searched to prevent just going back & forth
//(this is an optional parameter, based on careAboutPrev)
func prev(tag string) bool {
	if !careAboutPrev {
		return false
	}
	var temp *hashTag
	for temp = startList; temp.next != nil; temp = temp.next {
		if strings.EqualFold(tag, temp.tag) {
			return true
		}
	}
	return false
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
