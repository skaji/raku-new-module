package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IFTTT struct {
	Event string
	Key   string
}

func NewIFTTT(event string, key string) *IFTTT {
	return &IFTTT{Event: event, Key: key}
}

type body struct {
	Value1 string `json:"value1"`
}

func (i *IFTTT) Request(value string) error {
	b, err := json.Marshal(body{Value1: value})
	if err != nil {
		return nil
	}
	url := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", i.Event, i.Key)
	res, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
