package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type response struct {
	GravatarURL string `json:"gravatar_url"`
}

func getGravatarURL(CPANID string) (string, error) {
	res, err := http.Get("https://fastapi.metacpan.org/v1/author/" + CPANID)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	out := new(response)
	if err := json.Unmarshal(body, out); err != nil {
		return "", err
	}
	return out.GravatarURL, nil
}
