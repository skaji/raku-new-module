package distribution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RetryableError struct {
	Message string
}

func (e *RetryableError) Error() string {
	return e.Message
}

type Perl6Fetcher struct {
}

func NewPerl6Fetcher() *Perl6Fetcher {
	return &Perl6Fetcher{}
}

func (f *Perl6Fetcher) fetchMeta(ctx context.Context, d *Distribution) ([]byte, error) {
	url := fmt.Sprintf(
		"https://cpan.metacpan.org/authors/id/%s/%s/%s/Perl6/%s.meta",
		string(d.CPANID[0:1]),
		string(d.CPANID[0:2]),
		d.CPANID,
		d.Distvname,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%s, %s", res.Status, url)
	if res.StatusCode == http.StatusOK {
		return body, nil
	} else if res.StatusCode == http.StatusNotFound {
		return nil, &RetryableError{Message: msg}
	} else {
		return nil, errors.New(msg)
	}
}

type perl6Meta struct {
	Name string `json:"name"`
}

func (f *Perl6Fetcher) FetchName(ctx context.Context, d *Distribution) (string, error) {
	if !d.IsPerl6 {
		panic("not perl6!")
	}
	if len(d.CPANID) < 2 {
		return "", errors.New("too short CPANID")
	}
	body, err := f.fetchMeta(ctx, d)
	if err != nil {
		return "", err
	}
	var meta perl6Meta
	if err := json.Unmarshal(body, &meta); err != nil {
		return "", err
	}
	if name := meta.Name; name != "" {
		return name, nil
	}
	return "", errors.New("cannot find suitable main module name from 'name' in meta")
}
