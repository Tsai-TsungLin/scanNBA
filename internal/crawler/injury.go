package crawler

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchInjuryMap() map[string][]string {
	result := make(map[string][]string)

	res, err := http.Get("https://www.espn.com/nba/injuries")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".Table__league-injuries").Each(func(i int, s *goquery.Selection) {
		team := s.Find(".injuries__teamName").Text()
		var injuries []string

		s.Find(".Table__even").Each(func(j int, g *goquery.Selection) {
			name := strings.TrimSpace(g.Find(".AnchorLink").Text())
			status := strings.TrimSpace(g.Find(".col-stat").Text())
			comment := strings.TrimSpace(g.Find(".col-desc").Text())

			if name != "" && status != "" {
				injuries = append(injuries, name+" "+status+" "+comment)
			}
		})

		result[team] = injuries
	})

	return result
}
