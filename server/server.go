package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/olivere/elastic"
)

type HackerNews struct {
	Url       string `json:"url"`
	Time      int
	Title     string
	Highlight []string
}

func main() {
	e := echo.New()
	e.GET("/list", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		lists := getDocument()
		return c.JSON(http.StatusOK, lists)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func getDocument() []HackerNews {
	var result []HackerNews

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
	if exists {
		q := elastic.NewWildcardQuery("news", "*")
		highlight := elastic.NewHighlight().Field("news")

		searchResult, err := client.Search().
			Index("hackernews"). // search in index "twitter"
			Query(q).
			Highlight(highlight).
			From(0).Size(100).
			Pretty(true). // pretty print request and response JSON
			Do(ctx)       // execute
		if err != nil {
			// Handle error
			panic(err)
		}

		// Here's how you iterate hits with full control.
		if searchResult.Hits.TotalHits > 0 {
			fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index

				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var t HackerNews
				err := json.Unmarshal(*hit.Source, &t)
				if err != nil {
					// Deserialization failed
				}

				t.Highlight = hit.Highlight["news"]
				result = append(result, t)
			}
			return result
		} else {
			// No hits
			fmt.Print("Found no tweets\n")
		}
	}
	return result
}
