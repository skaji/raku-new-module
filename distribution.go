package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var distributionRegexp = regexp.MustCompile(`[^/]/[^/]{2}/([^/]+)/(Perl6/)?([^/]+)\.(?:tar\.gz|tar\.bz2|zip|tgz)`)
var versionRegexp = regexp.MustCompile(`^v?[\d_.]+$`)

// DistributionError is
type DistributionError struct {
	Message string
}

func (d *DistributionError) Error() string {
	return d.Message
}

// NewDistribution is
func NewDistribution(body string) (*Distribution, error) {
	res := distributionRegexp.FindAllStringSubmatch(body, -1)
	if len(res) == 0 {
		return nil, &DistributionError{"failed to parse"}
	}

	r := res[0]
	d := Distribution{
		CPANID:    r[1], // SKAJI
		Distvname: r[3], // App-cpm-0.987
		IsPerl6:   false,
		Pathname:  r[0], // S/SK/SKAJI/App-cpm-0.987.tar.gz
	}
	if r[2] == "Perl6/" {
		d.IsPerl6 = true
	}

	parts := strings.Split(d.Distvname, "-")
	for {
		if len(parts) < 2 {
			return nil, &DistributionError{fmt.Sprintf("%s does not have version", d.Distvname)}
		}
		if versionRegexp.MatchString(parts[len(parts)-1]) {
			break
		}
		parts = parts[:len(parts)-1]
	}

	d.Distname = strings.Join(parts[:len(parts)-1], "-")
	d.MainModule = strings.Join(parts[:len(parts)-1], "::")
	if d.IsPerl6 {
		name, err := d.guessPerl6MainModule()
		if err == nil {
			d.MainModule = name
		} else {
			log.Println(err)
		}
	}
	return &d, nil
}

// Distribution is
type Distribution struct {
	CPANID     string `json:"cpanid"`
	Distvname  string `json:"distvname"`
	Distname   string `json:"distname"`
	MainModule string `json:"main_module"`
	IsPerl6    bool   `json:"perl6"`
	Pathname   string `json:"pathname"`
}

func (d *Distribution) fetchPerl6Meta() ([]byte, error) {
	if len(d.CPANID) < 2 {
		return nil, errors.New("too short CPANID")
	}
	url := fmt.Sprintf(
		"https://cpan.metacpan.org/authors/id/%s/%s/%s/Perl6/%s.meta",
		string(d.CPANID[0:1]),
		string(d.CPANID[0:2]),
		d.CPANID,
		d.Distvname,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s, %s", res.Status, url)
	}
	return body, nil
}

type distributionMeta struct {
	Name string `json:"name"`
}

func (d *Distribution) guessPerl6MainModule() (string, error) {
	body, err := d.fetchPerl6Meta()
	if err != nil {
		return "", err
	}
	var meta distributionMeta
	if err := json.Unmarshal(body, &meta); err != nil {
		return "", err
	}
	if name := meta.Name; name != "" {
		return name, nil
	}
	return "", errors.New("cannot find suitable main module name from 'name' in meta")
}

// Summary is
func (d *Distribution) Summary() string {
	var url string
	if d.IsPerl6 {
		url = fmt.Sprintf("https://modules.perl6.org/dist/%s:cpan:%s", d.MainModule, d.CPANID)
	} else {
		url = fmt.Sprintf("https://metacpan.org/release/%s/%s", d.CPANID, d.Distvname)
	}
	return fmt.Sprintf("%s by %s\n%s", d.Distvname, d.CPANID, url)
}

// AsJSON is
func (d *Distribution) AsJSON() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
