package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/koshkaj/macler/parser/types"
	"gitlab.com/koshkaj/macler/parser/util"
	"io"
	"log"
	"net/http"
	"time"
)

func parseMyHome(ctx context.Context, url string, set map[string]struct{}, ch chan<- types.DataNormalized) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var jsonResponse types.Response

	body, _ := io.ReadAll(resp.Body)

	resp.Body.Close()
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Fatal(err)

	}
	layout := "2006-01-02 15:04:05"
	loc := util.LoadLocalTime()
	iterateData := &jsonResponse.Data.Prs

	for _, element := range *iterateData {
		t, err := time.Parse(layout, element.CreatedAt)
		t = t.In(loc)
		if err != nil {
			log.Fatal(err)
		}
		postedToday := t.YearDay() == time.Now().YearDay() && t.Year() == time.Now().Year()
		if !util.CheckKeyExists(set, element.ProductId) && postedToday && element.OwnerTypeID == "1" { // Ar gvaq nanaxi, Dges ari dadebuli, mesakutre ari
			element.Site = "myhome"
			element.Link = fmt.Sprintf("https://www.myhome.ge/ka/pr/%s/", element.ProductId)
		}
		ch <- element
	}
}

func RunMyhomeParser(ctx context.Context, seenProductIDs map[string]struct{}, ch chan<- types.DataNormalized, params types.IncommingEventMessage) {
	// get the amount of user
	url, err := util.BuildUrl("myhome", params.Filters)
	if err != nil {
		log.Fatal(err)
	}
	// TODO marto axlebs vugzavnit anu pirvel gverdze rac gaichiteba
	parseMyHome(ctx, url, seenProductIDs, ch)
}
