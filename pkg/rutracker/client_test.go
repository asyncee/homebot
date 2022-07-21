package rutracker

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func newClientWithAuth(t *testing.T) (*RutrackerClient, error) {
	var username, password string
	if username = os.Getenv("RUTRACKER_LOGIN"); username == "" {
		t.Skip("Please, set RUTRACKER_LOGIN to run this test")
	}
	if password = os.Getenv("RUTRACKER_PASSWORD"); password == "" {
		t.Skip("Please, set RUTRACKER_PASSWORD to run this test")
	}
	return newClient(Username(username), Password(password))
}

func newClient(username Username, password Password) (*RutrackerClient, error) {
	time.Sleep(5 * time.Second) // Prevent too frequent queries to rutracker.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return NewClient(username, password, jar)
}

func TestLoginEmptyFields(t *testing.T) {
	client, err := newClient("", "")

	if client != nil {
		t.Errorf("Failed to create client: %v:", client)
	}

	if err == nil {
		t.Errorf("Expected: %v, got: %v", "both username and password must be provided", err)
	}
}

func TestLogin(t *testing.T) {
	client, err := newClientWithAuth(t)

	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	err = client.login()
	if err != nil {
		t.Errorf("Failed to login: %v", err)
	}
}

func TestSearchUnauthorizedGoodCredentials(t *testing.T) {
	client, err := newClientWithAuth(t)

	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	_, err = client.Search("test")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSearchUnauthorizedBadCredentials(t *testing.T) {
	client, err := newClient("BAD", "CRED")

	if err != nil {
		t.Fatalf("Failed to create client, expected: %v, got: %v", nil, err)
	}

	_, err = client.Search("test")
	expected := "rutracker.org: bad response to login request: 200 OK"
	if err == nil {
		t.Fatal("search must fail because of bad credentials, but it's not")
	}
	if err.Error() != expected {
		t.Errorf("Failed to do search, expected: %v, got: %v", expected, err)
	}
}

func TestSearch(t *testing.T) {
	// If this method is failing then it should be run separately because rutracker
	// blocks too many login requests per second.
	client, err := newClientWithAuth(t)

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	query := "test"
	results, err := client.Search(query)
	if err != nil {
		t.Fatalf("Failed to do search: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("Search by %v returned zero results", query)
	}
}

func TestDownloadFile(t *testing.T) {
	header := `attachment; filename="[rutracker.org].t6228743.torrent"; filename*=UTF-8''%D0%94%D0%BE%D0%BA%D1%82%D0%BE%D1%80%20%D0%A1%D1%82%D1%80%D1%8D%D0%BD%D0%B4%D0%B6%20%D0%92%20%D0%BC%D1%83%D0%BB%D1%8C%D1%82%D0%B8%D0%B2%D1%81%D0%B5%D0%BB%D0%B5%D0%BD%D0%BD%D0%BE%D0%B9%20%D0%B1%D0%B5%D0%B7%D1%83%D0%BC%D0%B8%D1%8F%20Doctor%20Strange%20in%20the%20Multiverse%20of%20Madness%20%28%D0%A1%D1%8D%D0%BC%20%D0%A0%D1%8D%D0%B9%D0%BC%D0%B8%20Sam%20Raimi%29%20%5B2022%2C%20%D0%A1%D0%A8%D0%90%2C%20%D1%84%D0%B0%D0%BD%D1%82%D0%B0%D1%81%20%5Brutracker-6228743%5D.torrent`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-disposition", header)
		w.Write([]byte("test"))
	}))
	defer ts.Close()

	client, err := newClientWithAuth(t)

	if err != nil {
		t.Errorf("Failed to create client, expected: %v, got: %v", nil, err)
	}

	f, err := client.FetchTorrent(ts.URL, "token")
	if err != nil {
		t.Errorf("Failed to fetch torrent: %v", err)
	}
	defer f.Body.Close()

	expected_filename := "Доктор Стрэндж В мультивселенной безумия Doctor Strange in the Multiverse of Madness (Сэм Рэйми Sam Raimi) [2022, США, фантас [rutracker-6228743].torrent"
	if f.Name != expected_filename {
		t.Errorf("bad filename: expected: %v, got: %v", expected_filename, f.Name)
	}

	body, err := ioutil.ReadAll(f.Body)
	if err != nil {
		t.Errorf("Failed to read file body: %v", err)
	}

	expectedBody := []byte("test")
	if !cmp.Equal(body, expectedBody) {
		t.Errorf("Bad file body, expected: %v, got: %v", err, body)
	}
}
