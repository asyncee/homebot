package rutracker

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

var formTokenRe = regexp.MustCompile("form_token: '(.*)'")

func extractFormToken(script string) (string, error) {
	results := formTokenRe.FindStringSubmatch(script)
	if len(results) != 2 {
		return "", fmt.Errorf("failed to extract form token from script %s", script)
	}
	return results[1], nil
}

func parseSearchResults(reader io.Reader) ([]Torrent, error) {
	decoder := charmap.Windows1251.NewDecoder().Reader(reader)
	doc, err := goquery.NewDocumentFromReader(decoder)

	if err != nil {
		return nil, err
	}

	results := []Torrent{}

	// todo: log skipped records w/ reason
	formToken := ""
	doc.Find("head > script").Each(func(_ int, s *goquery.Selection) {
		if formToken != "" {
			return
		}

		if !strings.Contains(s.Text(), "form_token:") {
			return
		}
		token, err := extractFormToken(s.Text())
		if err != nil {
			return
		}
		formToken = token
	})

	doc.Find("table#tor-tbl tbody tr").Each(func(_ int, s *goquery.Selection) {
		titleSel := s.Find("td.t-title-col a")
		if titleSel.Length() == 0 {
			return
		}

		category := s.Find("td.f-name-col").Text()
		if category == "" {
			return
		}

		idText := titleSel.AttrOr("data-topic_id", "")
		if idText == "" {
			return
		}
		id, err := strconv.Atoi(idText)
		if err != nil {
			return
		}

		href := titleSel.AttrOr("href", "")
		if href == "" {
			return
		}

		title := titleSel.Text()
		if title == "" {
			return
		}

		status := s.Find("td.t-ico span").Parent().AttrOr("title", "")

		sizeSel := s.Find("td.tor-size a")
		if sizeSel.Length() == 0 {
			return
		}

		size := sizeSel.Text()
		downloadUrl, downloadUrlexists := sizeSel.Attr("href")
		if !downloadUrlexists {
			return
		}

		seedersText := s.Find("td b.seedmed").Text()
		if seedersText == "" {
			return
		}
		seeders, err := strconv.Atoi(seedersText)
		if err != nil {
			return
		}

		baseUrl := "https://rutracker.org/forum/"

		results = append(results, Torrent{
			ID:            id,
			Status:        normalizeString(status),
			Name:          normalizeString(title),
			Size:          normalizeString(size),
			Seeders:       seeders,
			DownloadUrl:   baseUrl + normalizeString(downloadUrl),
			DownloadToken: formToken,
			URL:           baseUrl + href,
			Category:      normalizeString(category),
		})
	})

	return results, nil
}

func normalizeString(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\u00a0", " ")
}
