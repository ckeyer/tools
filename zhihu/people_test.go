package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func getDoc(t *testing.T, name string) *goquery.Document {
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestDecode(t *testing.T) {
	doc := getDoc(t, "about.html")
	classnames := []string{
		"name",
		"location",
		"business",
		"gender",
		"employment",
		"position",
		"education",
		"education-extra",
		"zm-profile-header-description .content",
	}

	doc.Find(".zm-profile-header-main").Each(func(i int, s *goquery.Selection) {
		for k, v := range classnames {
			fmt.Printf("%d %s: %s\n", k, v, s.Find("."+v).Text())
		}
		// name := s.Find(".name").Text()
		// location := s.Find(".location").Text()
		// business := s.Find(".business").Text()
		// gender := s.Find(".gender").Text()

		// fmt.Printf("name: %s \n", name)
		// fmt.Printf("location: %s \n", location)
		// fmt.Printf("business: %s \n", business)
		// fmt.Printf("gender: %s \n", gender)
	})

	countSuffs := []string{"asks", "answers", "posts", "collections", "logs"}
	doc.Find(".zu-main-content-inner .profile-navbar a").Each(func(i int, s *goquery.Selection) {
		for _, suff := range countSuffs {
			if href, ok := s.Attr("href"); ok {
				if strings.HasSuffix(href, suff) {
					fmt.Printf("%s count: %s\n", suff, s.Find(".num").Text())
				}
			}
		}
	})

	doc.Find(".zm-profile-header-operation").Each(func(i int, s *goquery.Selection) {
		for k, v := range []string{"zm-profile-header-user-agree", "zm-profile-header-user-thanks"} {
			fmt.Printf("%d %s: %s\n", k, v, s.Find("."+v).Find("strong").Text())
		}
	})

	doc.Find(".zm-profile-details-wrap .zm-profile-module-desc span").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "赞同") {
			fmt.Printf("%d 赞同:%s\n", i, s.Find("strong").Text())
		}
		if strings.Contains(s.Text(), "感谢") {
			fmt.Printf("%d 感谢:%s\n", i, s.Find("strong").Text())
		}
		if strings.Contains(s.Text(), "收藏") {
			fmt.Printf("%d 收藏:%s\n", i, s.Find("strong").Text())
		}
		if strings.Contains(s.Text(), "分享") {
			fmt.Printf("%d 分享:%s\n", i, s.Find("strong").Text())
		}
	})

	t.Error("...")
}

func TestAsks(t *testing.T) {
	doc := getDoc(t, "asks.html")
	doc.Find("selector")

}
