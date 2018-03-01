package main

import (
	"fmt"
	"strconv"

	"golang.org/x/net/html"
)

type class struct {
	name    string
	teacher string
	id      string
	room    string
	year    string
	width   int
}

func main() {

}

func parser() {
	var schedule [4][11]class
	var (
		parse  func(*html.Node)
		period int
		day    int
	)
	period = 0
	parse = func(node *html.Node) {
		var width int

		if node.Type == html.ElementNode {
			if node.Parent
			for _, a := range node.Attr {
				if a.Key == "collspan" {
					width, err := strconv.Atoi(a.Val)
					if day == 0 {
						if err != nil {
							fmt.Println("couldnt parse")
						}
						var stuff = class{
							name:    "day",
							teacher: "none",
							id:      "none",
							room:    "none",
							year:    "none",
							width:   width,
						}
						schedule[day][period] = stuff
					}else{
						var stuff = class{
							name:    "day",
							teacher: "none",
							id:      "none",
							room:    "none",
							year:    "none",
							width:   width,
						}
					}
				}
			}
			for _, a := range node.Attr {
				if a.Key == "nowrap" {

				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}

	}
	parse()

}
