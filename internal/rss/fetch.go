package rss

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"myapp/internal/db"

	"github.com/mmcdole/gofeed"
)

// FetchRssFeeds парсит RSS-ленту и сохраняет новости в базу данных.
func FetchRssFeeds(ctx context.Context, url string, database *db.DB) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return fmt.Errorf("ошибка при получении RSS-ленты: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(feed.Items))

	for _, item := range feed.Items {
		wg.Add(1)
		go func(item *gofeed.Item) {
			defer wg.Done()
			pubTime, err := parsePubTime(item.Published)
			if err != nil {
				errCh <- fmt.Errorf("ошибка парсинга даты для новости '%s': %v", item.Title, err)
				return
			}

			post := db.Post{
				Title:   item.Title,
				Content: item.Description,
				PubTime: pubTime.Unix(),
				Link:    item.Link,
			}

			_, err = database.AddPost(ctx, &post)
			if err != nil {
				errCh <- fmt.Errorf("ошибка добавления новости '%s' в базу данных: %v", item.Title, err)
			}
		}(item)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Println(err)
	}

	return nil
}

// parsePubTime парсит строку времени с учетом возможных форматов.
func parsePubTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}

	dateStr = strings.Replace(dateStr, "GMT", "+0000", 1)

	for _, format := range formats {
		pubTime, err := time.Parse(format, dateStr)
		if err == nil {
			return pubTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить дату: %v", dateStr)
}
