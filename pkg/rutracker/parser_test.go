package rutracker

import (
	"bufio"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	f, err := os.Open("search_response.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	doc := bufio.NewReader(f)

	results, err := parseSearchResults(doc)

	if err != nil {
		t.Errorf("Failed to parse search results: %v", err)
	}

	expectedLen := 48
	if len(results) != expectedLen {
		t.Errorf("Results count does not match. Expected: %v, got: %v", expectedLen, len(results))
	}

	expectedFirstResult := Torrent{
		ID:            6229824,
		Status:        "проверено",
		Name:          "Summer International Test Match 2022 / Portugal - Italy / Португалия - Италия / Sky Sport [25.06.2022, Регби / Rugby, WEBRip/720p/30fps, MKV/H.264, IT]",
		Size:          "2.16 GB ↓",
		Seeders:       3,
		DownloadUrl:   "https://rutracker.org/forum/dl.php?t=6229824",
		DownloadToken: "1ebe3a62095d08b0c2bc6563e1bd9fc6",
		URL:           "https://rutracker.org/forum/viewtopic.php?t=6229824",
		Category:      "Регби",
	}

	if !cmp.Equal(results[0], expectedFirstResult) {
		t.Errorf("Result parsed wrong. Expected:\n%+v, got:\n%+v", expectedFirstResult, results[0])
		t.Errorf(cmp.Diff(results[0], expectedFirstResult))
	}
}

func TestExtractFormToken(t *testing.T) {
	token, err := extractFormToken(`
	window.BB = {
		cur_domain: location.hostname.replace(/.*?(([a-z0-9-]+\.){1,2}[a-z0-9-]+)$/, '$1'),
		form_token: '1ebe3a62095d08b0c2bc6563e1bd9fc6',
		opt_js: {"only_new":0,"h_flag":0,"h_av":0,"h_rnk_i":0,"h_post_i":0,"h_smile":0,"h_sig":0,"sp_op":0,"tr_tm":0,"h_cat":"","h_tsp":0,"h_ta":0},
	
		IS_GUEST: !!'',
		IMG_URL: 'https://static.t-ru.org/templates/v1/images',
		SMILES_URL: 'https://static.t-ru.org/smiles',
		catId: 0,
		FORUM_ID: 0,
		parentForumId: 0,
			PG_PER_PAGE: '50',
		PG_BASE_URL: 'tracker.php?search_id=PxGq96MfxnYz',
			COOKIE_MARK: 'bb_mark_read',
		};
	BB.cookie_defaults = {
		domain: '.' + BB.cur_domain,
		path: "/forum/",
	};
	`)
	if err != nil {
		t.Errorf("Failed to extract token: %v", err)
	}

	expectedToken := "1ebe3a62095d08b0c2bc6563e1bd9fc6"
	if token != expectedToken {
		t.Errorf("Bad expected token, expected: %s, got: %s", expectedToken, token)
	}
}
