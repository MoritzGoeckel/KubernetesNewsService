package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetArticle(document *goquery.Document, url string, source string) (Article, error) {
	/*html, err := document.Html()
	if err != nil {
		return Article{}, err
	}*/

	article := Article{Time: time.Now()}

	datapoints := 0
	datapointsStr := ""

	//Source

	article.Source = source
	datapointsStr += "Source "
	datapoints++

	//Title

	if title, exists := getMeta(document, "og:title"); exists {
		article.Headline = title
		datapointsStr += "Title "
		datapoints++
	} else {
		article.Headline = document.Find("h1").First().Text()
		if article.Headline != "" {
			datapointsStr += "Title "
			datapoints++
		}
	}

	//URL

	if ogURL, exists := getMeta(document, "og:url"); exists {
		article.Url = ogURL
	} else {
		article.Url = url
	}
	datapointsStr += "URL "
	datapoints++

	//Description

	if description, exists := getMeta(document, "og:description", "description"); exists {
		article.Description = description
		datapointsStr += "Description "
		datapoints++
	}

	//Image

	if image, exists := getMeta(document, "og:image", "twitter:image"); exists {
		article.Image = image
		datapointsStr += "Image "
		datapoints++
	}

	//Search for content
	//Should we add keywords?

	if datapoints <= 3 {
		return Article{}, errors.New("Not enough datapoints: " + fmt.Sprint(datapoints) + " -> " + datapointsStr)
	}

	return article, nil
}

//Make it take an array as input
func getMeta(document *goquery.Document, names ...string) (string, bool) {
	for _, name := range names {
		val, exists := document.Find("meta[name='" + name + "'], meta[property='" + name + "']").First().Attr("content")
		if exists {
			return val, exists
		}
	}

	return "", false
}

func GetHTML(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Status code error: " + fmt.Sprint(res.StatusCode) + " " + fmt.Sprint(res.Status) + "%s")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
