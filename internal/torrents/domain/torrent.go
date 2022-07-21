package domain

type Torrent struct {
	ID                string
	Status            string
	Name              string
	Size              string
	Seeders           int
	DownloadUrl       string
	DownloadCsrfToken string
	URL               string
	Category          string
}
