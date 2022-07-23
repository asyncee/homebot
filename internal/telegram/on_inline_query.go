package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/internal/torrents/domain"
	"github.com/asyncee/homebot/pkg/logging"
	"go.uber.org/fx"
	tele "gopkg.in/telebot.v3"
)

type InlineQueryHandler struct {
	fx.In
	Logger logging.Logger
	Repo   application.TorrentRepository
}

func (h *InlineQueryHandler) Handle(c tele.Context) error {
	text := c.Query().Text
	h.Logger.Debug("inline_query", text)

	if text == "" {
		return nil
	}

	// TODO: extract FindTorrentsByNameQuery

	torrents, err := h.Repo.FindByName(context.TODO(), text)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка: %s", err.Error()))
	}

	results := make(tele.Results, len(torrents))
	for i, torrent := range torrents {
		result := h.articleResult(torrent)
		results[i] = result
	}

	return c.Answer(&tele.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
}

func (h *InlineQueryHandler) articleResult(torrent domain.Torrent) *tele.ArticleResult {
	dsc := fmt.Sprintf("%s · %s · %s · %d на раздаче", torrent.Category, torrent.Size, torrent.Status, torrent.Seeders)
	thumbnail := "https://cdn-icons-png.flaticon.com/512/2521/2521768.png"
	movieCategories := []string{
		"классика мирового кинематографа", "фильм", "кино",
	}
	movieCategoryThumbnail := "https://cdn-icons-png.flaticon.com/512/3507/3507102.png"

	for _, movieCategory := range movieCategories {
		if strings.Contains(strings.ToLower(torrent.Category), movieCategory) {
			thumbnail = movieCategoryThumbnail
		}
	}

	result := &tele.ArticleResult{
		Title:       torrent.Name,
		Description: dsc,
		Text:        torrent.Name,
		URL:         torrent.URL,
		HideURL:     false,
		ThumbURL:    thumbnail,
	}
	result.SetResultID(strconv.Itoa(torrent.ID))
	return result
}
