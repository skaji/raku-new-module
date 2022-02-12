package distribution

import (
	"fmt"
	"strings"
	"time"
)

type Distribution struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Published time.Time `json:"date_published"`
	URL       string    `json:"url"`
}

func (d *Distribution) Summary() string {
	return fmt.Sprintf("%s by %s\n%s", d.Title, d.Auth(), d.URL)
}

func (d *Distribution) IsZef() bool {
	return strings.HasPrefix(d.Auth(), "zef:")
}

func (d *Distribution) IsCPAN() bool {
	return strings.HasPrefix(d.Auth(), "cpan:")
}

func (d *Distribution) Auth() string {
	parts := strings.Split(d.ID, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
