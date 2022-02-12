package stream

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/skaji/raku-new-module/pkg/distribution"
	"github.com/skaji/raku-new-module/pkg/log"
)

func NewRaku(ctx context.Context, url string, tick time.Duration) (<-chan *distribution.Distribution, error) {
	var lastSeen time.Time
	{
		log.Debugf("access %s", url)
		dists, err := getRecents(ctx, url)
		if err != nil {
			return nil, err
		}
		if len(dists) == 0 {
			return nil, errors.New("first fetch returns 0 dists!")
		}
		log.Debugf("returns %s %s ~ %s %s", dists[0].ID, dists[0].Published.Local().Format(time.RFC3339), dists[len(dists)-1].ID, dists[len(dists)-1].Published.Local().Format(time.RFC3339))
		dist := dists[len(dists)-1]
		log.Printf("will tweet after %s %s", dist.ID, dist.Published.Local().Format(time.RFC3339))
		lastSeen = dist.Published
	}

	ch := make(chan *distribution.Distribution, 10)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Debugf("access %s", url)
				dists, err := getRecents(ctx, url)
				if err != nil {
					log.Printf("%s %s", url, err.Error())
					continue
				}
				if len(dists) == 0 {
					log.Printf("%s returns 0 dists", url)
					continue
				}
				log.Debugf("returns %s %s ~ %s %s", dists[0].ID, dists[0].Published.Local().Format(time.RFC3339), dists[len(dists)-1].ID, dists[len(dists)-1].Published.Local().Format(time.RFC3339))
				for _, dist := range dists {
					if dist.Published.After(lastSeen) {
						ch <- dist
					}
				}
				lastSeen = dists[len(dists)-1].Published
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

func getRecents(ctx context.Context, url string) ([]*distribution.Distribution, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "raku-new-module; +https://github.com/skaji/raku-new-module")
	req.Close = true
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var resp struct {
		Items []*distribution.Distribution
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	dists := make([]*distribution.Distribution, len(resp.Items))
	for i, dist := range resp.Items {
		j := len(resp.Items) - i - 1
		dists[j] = dist
	}
	return dists, nil
}
