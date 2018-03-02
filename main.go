package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strconv"
	"io"
	"math"
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
	file, err := os.Open("shitty.html")
	parser(file)
	if (err!=nil){
		fmt.Println(err.Error)
	}
}

func parser(scheduleHtml io.ReadCloser) {
	var schedule [4][11][]class
	var (
		parse          func(*html.Node)
		period         int
		day            int
		cWidth, lWidth int
		err error
	)

	period = 0
	body, err := html.Parse(scheduleHtml)
	if (err!=nil){
		fmt.Println(err.Error)
	}
	parse = func(node *html.Node) {

		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == "colspan" {
					fmt.Println("correct key")
					cWidth, err = strconv.Atoi(a.Val)
					tWidth := cWidth + lWidth
					if period == 0 {
						if err != nil {
							fmt.Println("couldnt parse")
						}
						var stuff = class{
							name:    "day",
							teacher: "none",
							id:      "none",
							room:    "none",
							year:    "none",
							width:   cWidth,
						}
						fmt.Println(day)
						fmt.Println(period)
						schedule[day][period] = append(schedule[day][period],stuff)
					} else {
						html.Render(os.Stdout, node.FirstChild.FirstChild)
						var stuff = class{
							name:    node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Data,
							teacher: node.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.Data,
							id:      "none",
							room:    node.FirstChild.FirstChild.LastChild.FirstChild.FirstChild.Data,
							year:    "none",
							width:   cWidth,
						}
						schedule[day][period] = append(schedule[day][period],stuff)
					}
					if math.Mod(float64(tWidth), float64(schedule[0][0][0].width)) == 0{
						day++
						if day == 4{
							period++
							day=0
						}
					}
					node.FirstChild = nil
					lWidth = tWidth
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
	parse(body)
	//html.Render(os.Stdout, body)
	fmt.Println(schedule)

}
