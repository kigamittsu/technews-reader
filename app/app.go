package main

import (
	"encoding/json"
	"net/http"

	"github.com/gopherjs/vecty/prop"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type MyComponent struct {
	vecty.Core
}

type HackerNews struct {
	Url       string `json:"url"`
	Time      int
	Title     string
	Highlight []string
}

func (c *MyComponent) Render() vecty.ComponentOrHTML {
	req, _ := http.NewRequest("GET", "http://localhost:1323/list", nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// Handle error
		panic(err)
	}
	var d []HackerNews
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		// Handle error
		panic(err)
	}
	defer resp.Body.Close()

	vecty.SetTitle("HackingHackerNews")
	var array []vecty.MarkupOrChild
	for i := 0; i < len(d); i++ {
		array = append(array, elem.Div(
			vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Anchor(
					vecty.Markup(vecty.Class("card-title")),
					vecty.Markup(prop.Href(d[i].Url)),
					vecty.Text(d[i].Title),
				),
				elem.Paragraph(
					vecty.Markup(vecty.Class("card-text")),
					vecty.Text(d[i].Highlight[1]),
				),
			),
		),
		)
	}
	// 	<nav class="navbar navbar-light bg-light">
	//   <span class="navbar-brand mb-0 h1">Navbar</span>
	// </nav>
	return elem.Body(
		elem.Navigation(
			vecty.Markup(vecty.Class("navbar")),
			vecty.Markup(vecty.Class("navbar-light")),
			vecty.Markup(vecty.Class("bg-light")),
			elem.Span(vecty.Text("Hacking Hacker News")),
		),
		elem.Div(
			vecty.Markup(vecty.Class("container")),
			vecty.Markup(vecty.Style("margin-bottom", "50px")),
			vecty.Markup(vecty.Style("margin-top", "50px")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					array...,
				),
			),
		),
	)
}

func main() {
	vecty.AddStylesheet("https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css")
	top := &MyComponent{}
	vecty.RenderBody(top)
}
