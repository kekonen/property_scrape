package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

func main() {

	// output
	fName := "ebay_items.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write CSV header
	writer.Write([]string{"Name", "Price", "URL", "Image URL"})

	// Proxies
	dat, err := ioutil.ReadFile("./proxies.txt")
	dat1 := strings.Split(string(dat), "\n")
	rp, err := proxy.RoundRobinProxySwitcher(dat1...)
	if err != nil {
		log.Fatal(err)
	}

	listingCollector := colly.NewCollector(
		colly.AllowedDomains("www.ebay-kleinanzeigen.de"),
		colly.AllowURLRevisit(),
		colly.UserAgent("fwefopwekfopwjeopjwpojweopfjwpeofjopwefjwepofwe"),
	)
	itemCollector := colly.NewCollector(
		colly.AllowedDomains("www.ebay-kleinanzeigen.de"),
		colly.AllowURLRevisit(),
		colly.UserAgent("wepoofppwdofowieojfiowjehiowhfuiewhfuiofwiefpjwpiejfewpifwepjfopwjfopewjowjfopewjew"),
	)
	listingCollector.SetProxyFunc(rp)
	itemCollector.SetProxyFunc(rp)

	// Find and visit all links
	listingCollector.OnHTML("li.lazyload-item h2 a[href]", func(e *colly.HTMLElement) {

		// writer.Write([]string{
		// 	e.ChildAttr("a", "title"),
		// 	e.ChildText("span"),
		// 	e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
		// 	"https" + e.ChildAttr("img", "src"),
		// })
		fmt.Println(e.Attr("href"))
		itemCollector.Visit(fmt.Sprintf("https://www.ebay-kleinanzeigen.de%s", e.Attr("href")))
	})

	itemCollector.OnHTML("#viewad-product", func(e *colly.HTMLElement) {
		fmt.Println("KEK")
		// writer.Write([]string{
		// 	e.ChildAttr("a", "title"),
		// 	e.ChildText("span"),
		// 	e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
		// 	"https" + e.ChildAttr("img", "src"),
		// })
		// e.Request.Visit(e.Attr("href"))
		fmt.Println(e.ChildText("#viewad-title"))
		fmt.Println(e.ChildText("#viewad-price"))
		fmt.Println(e.ChildText("#viewad-locality"))
		fmt.Println(e.ChildText("#viewad-extra-info span"))
		fmt.Println(e.ChildText("#viewad-cntr span"))

	})

	listingCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", RandomString())
	})
	itemCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", RandomString())
	})

	listingCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("ListingRequest URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		listingCollector.Visit(r.Request.URL.String())
	})
	itemCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("ItemRequest URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		itemCollector.Visit(r.Request.URL.String())
	})

	// listingCollector.Visit("http://go-colly.org/")

	categories := [1][2]interface{}{
		{"s-wohnung-mieten", 203},
	}

	for _, c := range categories {
		for p := 1; p < 2; p++ {
			name := c[0]
			cid := c[1]
			fmt.Println(fmt.Sprintf("https://www.ebay-kleinanzeigen.de/%s/seite:%d/c%d", name, p, cid))
			listingCollector.Visit(fmt.Sprintf("https://www.ebay-kleinanzeigen.de/%s/seite:%d/c%d", name, p, cid))
		}

	}

	// err1 := listingCollector.Visit("https://www.ebay-kleinanzeigen.de/s-wohnung-mieten/seite:1/c203")
	// fmt.Println(err1)
	log.Println("listingCollector:\n", listingCollector)
	log.Println("itemCollector:\n", itemCollector)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
