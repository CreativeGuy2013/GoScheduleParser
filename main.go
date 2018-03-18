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
	Width   int    `json:"width"`
	State   string `json:"state"`
}

func (object class) MarshalJSON() ([]byte, error) {
	marschall := struct {
		Name    string `json:"name"`
		Teacher string `json:"teacher"`
		Room    string `json:"location"`
		State   string `json:"state"`
	}{
		Name:    object.Name,
		Teacher: object.Teacher,
		Room:    object.Room,
		State:   object.State,
	}

	json, _ := json.MarshalIndent(marschall, "", "    ")
	return json, nil
}

func main() {
	var a []chan int
	f, err := os.Open("./YTP")
	fis, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err.Error())
	}
	i := 0
	for v, fi := range fis {

		fmt.Println(fi.Name())
		file, err := os.Open(fmt.Sprintf("./YTP/%s", fi.Name()))
		c := make(chan int)
		a = append(a, c)
		go func(c chan int) {

			if err != nil {
				fmt.Println(err.Error())
			}
			defer file.Close()
			parser(file, file.Name())
			file.Close()
			c <- 1
		}(a[v])
		i++
	}

	for cannel := 0; cannel < i; cannel++ {
		_ = <-a[cannel]
	}
}

func parser(scheduleHtml io.ReadCloser, fileName string) {
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
		fmt.Println(err.Error())
	}
	parse = func(node *html.Node) {

		if node.Type == html.ElementNode {
			//fmt.Println(node.Data)
			for _, a := range node.Attr {

				if a.Key == "colspan" {
					fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------------")
					var stuff class
					lWidth = 0
					cWidth, err = strconv.Atoi(a.Val)
					for _, c := range schedule[period][day] {
						lWidth += c.Width
						fmt.Print("hey")
						fmt.Println(c.Width)
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
						fmt.Println(lWidth >= 12 && schedule[0][0][0].Width/lWidth >= 1)
						for lWidth >= 12 && schedule[0][0][0].Width/lWidth >= 1 {
							day++
							fmt.Println(day)
							if day == 5 {
								period++
								fmt.Println(period)

								day = 0
							}
							lWidth = 0
							for _, c := range schedule[period][day] {
								lWidth += c.Width
								fmt.Print("hey")
								fmt.Println(c.Width)
							}
						}

						if node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild == nil {
							fmt.Println("empty stuff")
							stuff = class{
								Name:    "none",
								Teacher: "none",
								Id:      "none",
								Room:    "none",
								Year:    "none",
								Width:   cWidth,
							}
						} else {
							//html.Render(os.Stdout, node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild)
							//fmt.Println(node.FirstChild.FirstChild.LastChild.Data)
							//fmt.Println(node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Data)

							var (
								name string
							)
							if node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.Data == "\nstart\n" {
								name = "start Toetsweek"
							} else {
								name = node.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild.Data
							}
							fmt.Println(name)

							stuff = class{
								Name:    name,
								Teacher: node.FirstChild.FirstChild.LastChild.FirstChild.FirstChild.FirstChild.Data,
								Id:      "none",
								Room:    node.FirstChild.FirstChild.LastChild.LastChild.PrevSibling.FirstChild.FirstChild.Data,
								Year:    "none",
								Width:   cWidth,
							}
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
					//fmt.Println(schedule)
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
	table, err := json.MarshalIndent(schedule, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	} //tableComplete, _ := json.MarshalIndent(schedule, "", "    ")

	ioutil.WriteFile(fmt.Sprintf("%s.json", fileName), table, 0644)
	fmt.Println(string(table))
}
