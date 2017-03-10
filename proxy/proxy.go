package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var port string
var target string

func Start() {
	port = "8080"
	target = "https://parahumans.wordpress.com"

	fmt.Printf("listening on: %s, proxing to: %s\n", port, target)

	http.HandleFunc("/", proxy)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rewriting: ", r.RequestURI)
	reqUrl, err := url.Parse(r.RequestURI)
	fatalIf(err)
	resp, err := http.Get(target + reqUrl.Path)
	fatalIf(err)
	tree, err := html.Parse(resp.Body)
	fatalIf(err)
	rewriteLinks(tree)
	html.Render(w, tree)
}

func rewriteLinks(node *html.Node) {
	targetUrl, err := url.Parse(target)
	fatalIf(err)
	targetHost := targetUrl.Host

	if node == nil {
		return
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		newAttr := make([]html.Attribute, 0)
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				url, err := url.Parse(attr.Val)
				fatalIf(err)
				if url.Host == targetHost {
					attr.Val = url.Path
				}
			}
			newAttr = append(newAttr, attr)
		}
		node.Attr = newAttr
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		rewriteLinks(n)
	}
}

func fatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
