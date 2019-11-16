package distribution

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var distributionRegexp = regexp.MustCompile(`[^/]/[^/]{2}/([^/]{2,})/(Perl6/)?([^/]+)\.(?:tar\.gz|tar\.bz2|zip|tgz)`)
var versionRegexp = regexp.MustCompile(`^v?[\d_.]+$`)

type Error struct {
	Message string
}

func (d *Error) Error() string {
	return d.Message
}

type Distribution struct {
	CPANID     string `json:"cpanid"`
	Distvname  string `json:"distvname"`
	Distname   string `json:"distname"`
	MainModule string `json:"main_module"`
	IsRaku     bool   `json:"raku"`
	ID         string `json:"-"`
}

func New(id string, str string) (*Distribution, error) {
	res := distributionRegexp.FindAllStringSubmatch(str, -1)
	if len(res) == 0 {
		msg := fmt.Sprintf("failed to parse string '%s'", str)
		return nil, &Error{msg}
	}

	r := res[0]
	d := Distribution{
		CPANID:    r[1], // SKAJI
		Distvname: r[3], // App-cpm-0.987
		IsRaku:    r[2] == "Perl6/",
		ID:        id,
	}

	parts := strings.Split(d.Distvname, "-")
	for {
		if len(parts) < 2 {
			msg := fmt.Sprintf("%s does not have version", d.Distvname)
			return nil, &Error{msg}
		}
		if versionRegexp.MatchString(parts[len(parts)-1]) {
			break
		}
		parts = parts[:len(parts)-1]
	}

	d.Distname = strings.Join(parts[:len(parts)-1], "-")
	d.MainModule = strings.Join(parts[:len(parts)-1], "::")
	return &d, nil
}

func (d *Distribution) Summary() string {
	var url string
	if d.IsRaku {
		url = fmt.Sprintf("https://modules.raku.org/dist/%s:cpan:%s", d.MainModule, d.CPANID)
	} else {
		url = fmt.Sprintf("https://metacpan.org/release/%s/%s", d.CPANID, d.Distvname)
	}
	return fmt.Sprintf("%s by %s\n%s", d.Distvname, d.CPANID, url)
}

func (d *Distribution) AsJSON() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (d *Distribution) MetaURL() string {
	return fmt.Sprintf(
		"https://cpan.metacpan.org/authors/id/%s/%s/%s/Perl6/%s.meta",
		string(d.CPANID[0:1]),
		string(d.CPANID[0:2]),
		d.CPANID,
		d.Distvname,
	)
}
