package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

func fetchUrlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	RSSFeed := RSSFeed{}

	err = xml.Unmarshal(dat, &RSSFeed)
	if err != nil {
		return RSSFeed, err
	}

	return RSSFeed, nil

}
