package util

import (
	"encoding/json"
	"fmt"
	"gitlab.com/koshkaj/macler/parser/types"
	"net/url"
	"strconv"
	"time"
)

func NormalizeData(site string, data *types.IncommingEventMessage) {
	if site == "myhome" {
		cityMappings := map[string]string{
			"თბილისი": "1996871",
			"ბათუმი":  "8742159",
			"გუდაური": "540059524",
		}
		propertyTypeMappings := map[string]string{
			"იყიდება":   "1",
			"ქირავდება": "3",
		}
		adTypeMapping := map[string]string{
			"ბინა":  "1",
			"სახლი": "2",
		}
		data.Filters.ADType = adTypeMapping[data.Filters.ADType]
		data.Filters.PropertyType = propertyTypeMappings[data.Filters.PropertyType]
		data.Filters.City = cityMappings[data.Filters.City]
	}

}

func CheckKeyExists(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}

func ValidateEventMessage(message []byte) (types.IncommingEventMessage, error) {
	var jsonData types.IncommingEventMessage
	err := json.Unmarshal(message, &jsonData)
	return jsonData, err
}

func BuildUrl(provider string, filters types.EventMessageFilter) (string, error) {
	if provider == "myhome" {
		baseURL := "https://www.myhome.ge/ka/s/"
		params := url.Values{}
		keyword := filters.Keyword
		params.Add("AdTypeID", filters.ADType)
		params.Add("PrTypeID", filters.PropertyType)
		params.Add("FCurrencyID", "1")
		params.Add("FPriceFrom", strconv.Itoa(filters.PriceMin))
		params.Add("FPriceTo", strconv.Itoa(filters.PriceMax))
		params.Add("AreaSizeFrom", strconv.Itoa(filters.SquareMin))
		params.Add("AreaSizeTo", strconv.Itoa(filters.SquareMax))
		params.Add("cities", filters.City)
		params.Add("Ajax", "1")

		if filters.EmergencyListing {
			keyword = fmt.Sprintf("%s %s", filters.Keyword, "სასწრაფოდ")
		}
		params.Add("Keyword", keyword)
		u, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}
		u.RawQuery = params.Encode()
		return u.String(), nil
	} else {
		return "", nil
	}
}

func PrintStartMessage() {
	message := `
	   ___  ___ ________ ___ ____
	  / _ \/ _ '/ __(_-</ -_) __/
	 / .__/\_,_/_/ /___/\__/_/   
	/_/                          

	Parser App Has Started
	_________________________________________________________________________`
	fmt.Println(message)
}

func LoadLocalTime() *time.Location {
	gmt4 := time.FixedZone("GMT+4", 4*60*60)
	return gmt4
}
