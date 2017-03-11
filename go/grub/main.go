package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
    "time"
)

const (
	BASEURL  = "http://www.bz55.com"
    //DOWNPATH = "/Users/didi/goproj/src/grub/image/"
	DOWNPATH = "/home/neagle/goproj/src/grub/image/"
    CHANNEL_CACHE_LENGTH = 10
    CHANNEL_DOWN_QUEUE_LENGTH = 2
)

var urlChl chan string
var downChl chan string

func findElement(pageUrl string) {
	doc, err := goquery.NewDocument(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("p").Eq(2).Find("img").Each(func(i int, s *goquery.Selection) {
		src, find := s.Attr("src")
		if find {
			urlChl <- src
			logOut(fmt.Sprintf("parse down url: %s, channel length: %d" ,src, len(urlChl)))
		}
	})
}

func downFile(fileUrl string) {
	fullUrl := BASEURL + fileUrl
	//path, filename := filepath.Split(fullUrl)
	_, filename := filepath.Split(fullUrl)
    file := DOWNPATH + filename
    logOut(file + " now is downloading ... ...")
	res, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal("down file error, url:" + fullUrl)
	}
	defer res.Body.Close()
	f, err := os.Create(file)
	if err != nil {
		log.Fatal("create file wrong")
	}
	defer f.Close()
	io.Copy(f, res.Body)
    logOut(file + " now is download finished")
}

func outCacheChannel() {
    i := 0
	for {
        i++
        logOut(fmt.Sprintf("cache loop: %d, channel length:%d", i, len(urlChl)))
		select {
		    case url := <-urlChl:
                downChl <- url
		}
	}
}

func downChannel() {
    i := 0
    for {
        i++
        logOut(fmt.Sprintf("down loop: %d, channel length:%d", i, len(downChl)))
        select {
            case url := <-downChl:
                downFile(url)
        }
    }
}

func logOut(s string) {
    log.Printf("%+v\t%s\n", time.Now().Format("2006-01-02 15:04:05"), s)
}

func main() {
    urls := []string{
        "http://www.bz55.com/meinvbizhi/38925.html",
        "http://www.bz55.com/meinvbizhi/38872.html",
        "http://www.bz55.com/meinvbizhi/38818.html",
        "http://www.bz55.com/meinvbizhi/38726.html",
        "http://www.bz55.com/meinvbizhi/38444.html",
        "http://www.bz55.com/meinvbizhi/38365.html",
        "http://www.bz55.com/meinvbizhi/38039.html",
        "http://www.bz55.com/meinvbizhi/25443.html",
        "http://www.bz55.com/meinvbizhi/37989.html",
        "http://www.bz55.com/meinvbizhi/37692.html",
        "http://www.bz55.com/meinvbizhi/37359.html",
        "http://www.bz55.com/meinvbizhi/36985.html",
        "http://www.bz55.com/meinvbizhi/36876.html",
        "http://www.bz55.com/meinvbizhi/36898.html",
        "http://www.bz55.com/meinvbizhi/36741.html",
    }

	go outCacheChannel()
    go downChannel()
	urlChl = make(chan string, CHANNEL_CACHE_LENGTH)
    downChl = make(chan string, CHANNEL_DOWN_QUEUE_LENGTH)
    for i, url := range urls {
        fmt.Printf("%v\t%#v\n", i, string(url) )
        findElement(string(url))
    }
}
