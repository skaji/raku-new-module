package main

import (
	"encoding/json"
	"fmt"
)

type Distribution struct {
	CPANID     string `json:"cpanid"`
	Distvname  string `json:"distvname"`
	Distname   string `json:"distname"`
	MainModule string `json:"main_module"`
	IsPerl6    bool   `json:"perl6"`
	Pathname   string `json:"pathname"`
}

func (d *Distribution) Summary() string {
	var url string
	if d.IsPerl6 {
		url = fmt.Sprintf("https://modules.perl6.org/dist/%s:cpan:%s", d.MainModule, d.CPANID)
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
