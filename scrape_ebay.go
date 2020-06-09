package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
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
		colly.UserAgent(browser.Random()),
		colly.Async(true),
		colly.IgnoreRobotsTxt(),
	)
	itemCollector := colly.NewCollector(
		colly.AllowedDomains("www.ebay-kleinanzeigen.de"),
		colly.AllowURLRevisit(),
		colly.UserAgent(browser.Random()),
		colly.Async(true),
		colly.IgnoreRobotsTxt(),
	)
	listingCollector.SetProxyFunc(rp)
	itemCollector.SetProxyFunc(rp)

	listingCollector.SetRequestTimeout(time.Second * 25)
	itemCollector.SetRequestTimeout(time.Second * 25)

	listingCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
	})
	itemCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
	})

	// Find and visit all links
	listingCollector.OnHTML("li.lazyload-item h2 a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Success", e.Attr("href"))

		// ctx := colly.NewContext()
		// ctx.Put("category", name)
		url := fmt.Sprintf("https://www.ebay-kleinanzeigen.de%s", e.Attr("href"))
		listingCollector.Request("GET", url, nil, e.Request.Ctx, nil)
		// itemCollector.Visit()
	})

	itemCollector.OnHTML("#viewad-product", func(e *colly.HTMLElement) {
		fmt.Println("KEK", e.Request.Ctx.Get("category"), e.ChildText("#viewad-title"), e.ChildText("#viewad-price"), e.ChildText("#viewad-locality"), e.ChildText("#viewad-extra-info span"), e.ChildText("#viewad-cntr span"))

		// TODO: parse url

		// fmt.Println(e.ChildText("#viewad-title"))
		// fmt.Println(e.ChildText("#viewad-price"))
		// fmt.Println(e.ChildText("#viewad-locality"))
		// fmt.Println(e.ChildText("#viewad-extra-info span"))
		// fmt.Println(e.ChildText("#viewad-cntr span"))

		// TODO: Get details
		// TODO: Ausstattung
		// TODO: Beschreibung
		// TODO: Ad owner

		// writer.Write([]string{
		// 	e.ChildAttr("a", "title"),
		// 	e.ChildText("span"),
		// 	e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
		// 	"https" + e.ChildAttr("img", "src"),
		// })
		// e.Request.Visit(e.Attr("href"))
	})

	listingCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", browser.Random())
	})
	itemCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", browser.Random())
	})

	listingCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("UA L:", r.Request.Headers.Clone().Get("User-Agent"))
		fmt.Println("ListingRequest URL:", r.Request.URL, "failed with response:", r, "\nError:", err) //, r.Headers.Get("User-Agent")

		r.Request.Visit(r.Request.URL.String())
	})
	itemCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("UA I:", r.Request.Headers.Clone().Get("User-Agent"))
		fmt.Println("ItemRequest URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Visit(r.Request.URL.String())
	})

	// listingCollector.Visit("http://go-colly.org/")

	categories := [1][2]interface{}{
		{"s-wohnung-mieten", 203},
	}

	for _, c := range categories {
		for p := 1; p < 6; p++ {
			name := c[0]
			cid := c[1]
			url := fmt.Sprintf("https://www.ebay-kleinanzeigen.de/%s/seite:%d/c%d", name, p, cid)
			ctx := colly.NewContext()
			ctx.Put("category", name)
			listingCollector.Request("GET", url, nil, ctx, nil)
		}

	}

	// err1 := listingCollector.Visit("https://www.ebay-kleinanzeigen.de/s-wohnung-mieten/seite:1/c203")
	// fmt.Println(err1)
	listingCollector.Wait()
	itemCollector.Wait()

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
