package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type class struct {
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
	Id      string `json:"j"`
	Room    string `json:"location"`
	Year    string `json:"j"`
	Width   int    `json:"j"`
	State   string `json:"state"`
}

func main() {
	file, err := os.Open("shitty.html")
	parser(file)
	if err != nil {
		fmt.Println(err.Error)
	}
}

func parser(scheduleHtml io.ReadCloser) {
	var schedule [10][5][]class
	var (
		parse                  func(*html.Node)
		period                 int
		day                    int
		cWidth, lWidth, tWidth int
		err                    error
	)

	period = 0
	body, err := html.Parse(scheduleHtml)
	if err != nil {
		fmt.Println(err.Error)
	}
	parse = func(node *html.Node) {

		if node.Type == html.ElementNode {
			//fmt.Println(node.Data)
			for _, a := range node.Attr {

				if a.Key == "colspan" {
					//fmt.Println("correct key")
					var stuff class
					lWidth = 0
					cWidth, err = strconv.Atoi(a.Val)
					for _, c := range schedule[period][day] {
						lWidth += c.Width
					}
					tWidth = cWidth + lWidth
					//fmt.Printf("%d-%d\n", day, period)
					//fmt.Println(tWidth)
					if period == 0 {
						if err != nil {
							fmt.Println("couldnt parse")
						}
						stuff = class{
							Name:    "day",
							Teacher: "none",
							Id:      "none",
							Room:    "none",
							Year:    "none",
							Width:   cWidth,
						}
						//fmt.Println(day)
						//fmt.Println(period)
					} else {

						fmt.Println("--------------------------parsing class--------------------------")
						if lWidth >= 12 && schedule[0][0][0].Width/lWidth >= 1 {
							day++
							if day == 5 {
								period++
								day = 0
							}
						}
						if node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild == nil {
							fmt.Println("empty stuff")
						} else {
							//html.Render(os.Stdout, node.FirstChild.FirstChild.LastChild.LastChild.PrevSibling)
							//fmt.Println(node.FirstChild.FirstChild.LastChild.Data)

							stuff = class{
								Name:    node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild.Data,
								Teacher: node.FirstChild.FirstChild.LastChild.FirstChild.FirstChild.FirstChild.Data,
								Id:      "none",
								Room:    node.FirstChild.FirstChild.LastChild.LastChild.PrevSibling.FirstChild.FirstChild.Data,
								Year:    "none",
								Width:   cWidth,
							}
							fmt.Println(node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Data)
							if stuff.Name == "strike" {
								stuff.State = "canceled"
								stuff.Name = node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.Data
								stuff.Teacher = node.FirstChild.FirstChild.LastChild.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild.Data
							} else {
								for _, a := range node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Attr {
									if a.Key == "color" {
										if strings.ToUpper(a.Val) == "#FF0000" {
											stuff.State = "changed"
										}
									}
								}
							}

						}
					}
					periods := 1
					//html.Render(os.Stdout, node)
					for _, a := range node.Attr {
						if a.Key == "rowspan" {
							g, _ := strconv.Atoi(a.Val)
							periods = g / 2
						}
					}
					for i := 0; i < periods; i++ {
						schedule[period+i][day] = append(schedule[period+i][day], stuff)
					}
					fmt.Println(tWidth)
					if math.Mod(float64(tWidth), float64(schedule[0][0][0].Width)) == 0 {
						day++
						if day == 5 {
							period++
							day = 0
						}
					}
					node.FirstChild = nil
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
	table, _ := json.MarshalIndent(schedule[1:], "", "    ")
	tableComplete, _ := json.MarshalIndent(schedule, "", "    ")

	ioutil.WriteFile("scedule.json", table, 0644)
	fmt.Println(string(tableComplete))
}
