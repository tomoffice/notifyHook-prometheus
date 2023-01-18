package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func MakeAlert(u string) (string, error) {
	url := u
	method := "POST"
	testJson := `
	{
		"receiver": "api-receiver",
		"status": "resolved",
		"alerts": [
			{
				"status": "resolved",
				"labels": {
					"alertname": "goApp flip",
					"instance": "golang:9001",
					"job": "goApp",
					"severity": "warm"
				},
				"annotations": {
					"description": "Failed to scrape goApp on golang:9001 for more than 3 minutes. Node seems down.",
					"title": "flip golang:9001 is more 5"
				},
				"startsAt": "2023-01-12T08:09:31.656Z",
				"endsAt": "2023-01-12T08:10:01.656Z",
				"generatorURL": "http://07712047c817:9090/graph?g0.expr=goApp_flip+%3E+5\u0026g0.tab=1",
				"fingerprint": "989b0a76f5d22d7f"
			}
		],
		"groupLabels": {
			"alertname": "goApp flip"
		},
		"commonLabels": {
			"alertname": "goApp flip",
			"instance": "golang:9001",
			"job": "goApp",
			"severity": "warm"
		},
		"commonAnnotations": {
			"description": "Failed to scrape goApp on golang:9001 for more than 3 minutes. Node seems down.",
			"title": "flip golang:9001 is more 5"
		},
		"externalURL": "http://28fdaf41f3ea:9093",
		"version": "4",
		"groupKey": "{}/{alertname=\"goApp flip\"}:{alertname=\"goApp flip\"}",
		"truncatedAlerts": 0
	}`
	payload := strings.NewReader(testJson)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))
	return string(body), nil
}
