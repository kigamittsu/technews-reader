package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/olivere/elastic"
)

type Article struct {
	By          string
	Descendants int
	Id          int
	Kids        []int
	Score       int
	Time        int
	Title       string
	Url         string
}

type HackerNews struct {
	News  string `json:"news"`
	Url   string `json:"url"`
	Time  int
	Title string
}

func GetTopStories() []string {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/newstories.json")
	if err != nil {
		// handle error
		fmt.Println("error!", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error!", err)
	}
	bodyString := strings.Replace(string(body), "[", "", -1)
	bodyString = strings.Replace(bodyString, "]", "", -1)
	bodyArray := strings.Split(bodyString, ",")
	return bodyArray[0:30]
}

func GetUrl(topStories []string) []Article {
	var topUrls []Article
	for i := 0; i < len(topStories); i++ {
		resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", topStories[i]))
		if err != nil {
			// handle error
			fmt.Println("error!", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var target Article
		json.Unmarshal(body, &target)
		article := Article{Title: target.Title, Time: target.Time, Url: target.Url}
		topUrls = append(topUrls, article)
		//fmt.Println(target.Url)
		time.Sleep(1 * time.Second)
	}
	return topUrls
}

func GetArticle(articles []Article) []HackerNews {
	var result []HackerNews
	// c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})
	for i := 0; i < len(articles); i++ {
		c := colly.NewCollector()
		// Find and visit all links
		c.OnHTML("body", func(e *colly.HTMLElement) {
			body := HackerNews{Title: articles[i].Title, Time: articles[i].Time, Url: articles[i].Url, News: e.Text}
			result = append(result, body)
			// e.ForEach("div", func(_ int, elem *colly.HTMLElement) {
			// 	classname := elem.Attr("class")
			// 	if strings.Contains(classname, "container") {
			// 		body := HackerNews{Title: articles[i].Title, Time: articles[i].Time, Url: articles[i].Url, News: elem.Text}
			// 		result = append(result, body)
			// 	}
			// 	// fmt.Printf("%d", len(strings.Trim(elem.Text, " ")))
			// 	// fmt.Println(strings.Contains(classname, "contain"))
			// 	// if strings.Contains(elem.Attr("class"), "contain") {
			// 	// 	fmt.Println(strings.Trim(elem.Text, " "))
			// 	// 	result = append(result, elem.Text)
			// 	// }
			// })
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.Visit(articles[i].Url)
	}
	return result
}

func storeES(result []HackerNews, url []Article) {
	fmt.Println(len(result))
	fmt.Println(len(url))
	mapping := `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"hackernews":{
			"properties":{
				"title": {
					"type":"text"
				},
				"news": {
					"type":"text"
				},
				"url": {
					"type":"text"
				},
				"time": {
					"type":"text"
				}
			}
		}
	}
}`

	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("hackernews").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("hackernews").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	for i := 0; i < len(result); i++ {
		// body := `{"hackernews" : {"news" : ` + result[i] + `}}`
		// body := HackerNews{result[i]}
		// id := strconv.Itoa(i + 1)
		put, err := client.Index().
			Index("hackernews").
			Type("hackernews").
			BodyJson(result[i]).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed hackernews %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	}

}

func main() {
	topStories := GetTopStories()
	urls := GetUrl(topStories)
	result := GetArticle(urls)
	// var result []string
	storeES(result, urls)
}
