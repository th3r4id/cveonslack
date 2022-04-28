package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type CveData struct {
	Id          string
	Description string
	Link        string
}

func main() {
	freshData := fetchCveData()
	newCVEs := findNewCVEs(freshData)
	sendSlackAlert(newCVEs)
	appendNewCVEs(newCVEs)
}

func fetchCveData() (output []CveData) {
	resp, err := http.Get("https://www.cvedetails.com/json-feed.php")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(body)
	}
	//fmt.Println(string(body))
	var tmp []map[string]interface{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Fatal(err)
	}

	output = make([]CveData, 0)

	for _, item := range tmp {
		output = append(output, CveData{
			Id:          item["cve_id"].(string),
			Description: item["summary"].(string),
			Link:        item["url"].(string),
		})
	}
	return
}

func findNewCVEs(freshData []CveData) (output []CveData) {
	dataRead, err := os.ReadFile("cvesearch.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(dataRead), "\n")
	knownCVEs := map[string]bool{}

	for _, cveId := range lines {
		if len(cveId) == 0 {
			continue
		}
		knownCVEs[cveId] = true
	}

	output = make([]CveData, 0)
	for _, item := range freshData {
		if _, ok := knownCVEs[item.Id]; !ok {
			output = append(output, item)
		}
	}
	return
}

func sendSlackAlert(newCVEs []CveData) {
	SlackAccessToken := "add-token-here"
	api := slack.New(SlackAccessToken)
	for _, item := range newCVEs {
		_, _, err := api.PostMessage(
			"channel-id",
			slack.MsgOptionText(item.Id+"\n"+item.Description+"\n"+item.Link, false),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func appendNewCVEs(newCVEs []CveData) {
	f, err := os.OpenFile("cvesearch.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range newCVEs {
		if _, err := f.Write([]byte(item.Id + "\n")); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
