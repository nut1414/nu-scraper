package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
	"github.com/gocolly/colly"
)

type novel struct {
  Name              string 
  ChapterCount  int
  UpdatesFrequency   string
  Readers           int
  Reviews           int
  LastUpdated        string
  ImageUrl          string
  Url               string
  OriginLanguage    string
  Rating            double
}

func main() {
  c := colly.NewCollector(
    colly.AllowedDomains("www.novelupdates.com"),
    )
 
  transport, _ := cfrt.New(&http.Transport{
	  Proxy: http.ProxyFromEnvironment,
	  DialContext: (&net.Dialer{
	  	Timeout:   30 * time.Second,
		  KeepAlive: 30 * time.Second,
		  DualStack: true,
	  }).DialContext,
	  MaxIdleConns:          100,
	  IdleConnTimeout:       90 * time.Second,
	  TLSHandshakeTimeout:   10 * time.Second,
	  ExpectContinueTimeout: 1 * time.Second,
  })

  c.WithTransport(transport)

  c.OnHTML(".search_main_box_nu", func(e *colly.HTMLElement) {
    count, _ := strconv.Atoi(strings.Split(e.ChildText(".search_stats > span:nth-child(1)"), " ")[0])
    readers, _ := strconv.Atoi(strings.Split(e.ChildText(".search_stats > span:nth-child(3)"), " ")[0])
    reviews, _ := strconv.Atoi(strings.Split(e.ChildText(".search_stats > span:nth-child(4)"), " ")[0])

    novel := novel{
      Name: e.ChildText(".search_title > a"),
      Url:  e.ChildAttr(".search_title > a", "href"),
      ChapterCount: count,
      UpdatesFrequency:  e.ChildText(".search_stats > span:nth-child(2)"),
      Readers:  readers,
      Reviews:  reviews,
      LastUpdated: e.ChildText(".search_stats > span:nth-child(5)"),
      ImageUrl: e.ChildAttr(".search_img_nu > img", "src"),
    }  
    fmt.Println(novel)

  })

  c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())

	})
  
  c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.StatusCode)
		fmt.Println(r.Request.Headers)
	})

  c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
		fmt.Println(r.Request.Headers)
	})

  c.Visit("https://www.novelupdates.com/genre/comedy/")
}
