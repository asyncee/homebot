package rutracker

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	stdjar "net/http/cookiejar"
	"net/url"
	"sort"
	"time"

	"github.com/asyncee/homebot/pkg/logging"
	cookiejar "github.com/juju/persistent-cookiejar"
)

type Torrent struct {
	ID            int
	Status        string
	Name          string
	Size          string
	Seeders       int
	DownloadUrl   string
	DownloadToken string
	URL           string
	Category      string
}

type RutrackerClient struct {
	Username   string
	Password   string
	Jar        http.CookieJar
	persistent bool
	logger     logging.Logger
}

func NewClient(username, password string, persistent bool) (*RutrackerClient, error) {
	if username == "" || password == "" {
		return nil, errors.New("both username and password must be provided")
	}
	var jar http.CookieJar

	if persistent {
		persistentJar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		jar = persistentJar
	} else {
		inmemoryJar, err := stdjar.New(nil)
		if err != nil {
			return nil, err
		}
		jar = inmemoryJar
	}

	return &RutrackerClient{
		Username:   username,
		Password:   password,
		Jar:        jar,
		logger:     logging.GetLogger(),
		persistent: persistent,
	}, nil
}

func (c *RutrackerClient) login() error {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: c.Jar,
	}

	resp, err := client.PostForm("https://rutracker.org/forum/login.php", url.Values{
		"login_username": {c.Username},
		"login_password": {c.Password},
		"login":          {"%E2%F5%EE%E4"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return fmt.Errorf("rutracker.org: bad response to login request: %s, redirect expected", resp.Status)
	}

	if c.persistent {
		persistentJar := c.Jar.(*cookiejar.Jar)
		persistentJar.Save()
	}

	c.logger.Info("msg", "successfully logged-in to rutracker.org")
	return nil
}

func (c *RutrackerClient) Search(query string) ([]Torrent, error) {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: c.Jar,
	}

	searchUrl, err := url.Parse("https://rutracker.org/forum/tracker.php")
	if err != nil {
		return nil, err
	}
	q := searchUrl.Query()
	q.Add("nm", query)

	resp, err := client.PostForm(searchUrl.String(), url.Values{
		"max": {"1"},
		"nm":  {query},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		c.logger.Debug("msg", "rutracker search failed: not authorized")

		err = c.login()
		if err != nil {
			return nil, err
		}
		return c.Search(query)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("rutracker.org: bad response to search request: %s", resp.Status)
	}

	c.logger.Debug("msg", "rutracker search succeeded")
	results, err := parseSearchResults(resp.Body)
	if err != nil {
		c.logger.Error("msg", "failed to parse search results")
		return nil, err
	}
	c.logger.Debug("msg", "parsed %d search result(s)", len(results))

	sort.Slice(results, func(i, j int) bool {
		return results[i].Seeders > results[j].Seeders
	})

	return results, nil
}

type TorrentFile struct {
	Name string
	Body io.ReadCloser
}

func (c *RutrackerClient) FetchTorrent(downloadUrl, formToken string) (*TorrentFile, error) {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: c.Jar,
	}

	resp, err := client.PostForm(downloadUrl, url.Values{
		"form_token": {formToken},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 302 {
		c.logger.Debug("msg", "rutracker download file failed: not authorized")

		err = c.login()
		if err != nil {
			return nil, err
		}
		return c.FetchTorrent(downloadUrl, formToken)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("rutracker.org: bad response to download request: %s", resp.Status)
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, err
	}

	f := TorrentFile{
		Name: params["filename"],
		Body: resp.Body,
	}

	return &f, nil
}
