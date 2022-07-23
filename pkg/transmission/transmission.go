package transmission

import (
	"context"
	"fmt"

	"github.com/hekmon/transmissionrpc/v2"
)

type TransmissionClient struct {
	client *transmissionrpc.Client
}

type Host string
type User string
type Password string

func NewClient(
	host Host,
	user User,
	password Password,
) (*TransmissionClient, error) {
	client, err := transmissionrpc.New(string(host), string(user), string(password), nil)
	if err != nil {
		return nil, err
	}
	return &TransmissionClient{
		client: client,
	}, nil
}

func (s *TransmissionClient) TestConnection() error {
	ctx, cancel := context.WithCancel(context.TODO())
	ok, serverVersion, serverMinimumVersion, err := s.client.RPCVersion(ctx)
	defer cancel()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(
			"remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion,
			transmissionrpc.RPCVersion,
			serverMinimumVersion,
		)
	}
	return nil
}

func (s *TransmissionClient) ListTorrents() ([]transmissionrpc.Torrent, error) {
	torrents, err := s.client.TorrentGetAll(context.TODO())
	if err != nil {
		return nil, err
	}
	return torrents, nil
}

type Torrent struct {
	ID   int64
	Name string
}

func (s *TransmissionClient) DownloadTorrent(filepath string) (*Torrent, error) {
	torrent, err := s.client.TorrentAddFile(context.TODO(), filepath)
	if err != nil {
		return nil, err
	}
	return &Torrent{
		ID:   *torrent.ID,
		Name: *torrent.Name,
	}, nil
}

type TorrentStatus int

const (
	TorrentDone    = 1
	TorrentRunning = 2
	TorrentMissing = 3
	TorrentError   = 4
)

func (s *TransmissionClient) TorrentIsDone(torrentId int64) (TorrentStatus, error) {
	torrents, err := s.client.TorrentGet(context.TODO(), []string{"status"}, []int64{torrentId})
	if err != nil {
		return TorrentError, err
	}

	if len(torrents) == 0 {
		return TorrentMissing, nil
	}

	status := *torrents[0].Status
	if status == transmissionrpc.TorrentStatusSeed || status == transmissionrpc.TorrentStatusSeedWait || status == transmissionrpc.TorrentStatusStopped {
		return TorrentDone, nil
	}

	return TorrentRunning, nil
}
