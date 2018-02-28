package main

import "golang.org/x/net/html"

func main() {

}

func parser() {
	var (
		parse func(*html.Node)
	)

	parse = func(node *html.Node) {
		if node.Type == html.ElementNode {
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
