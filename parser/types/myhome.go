package types


type MyHomeData struct {
	Site          string      `json:"site"`
	Link          string      `json:"link"`
	ProductId     string      `json:"product_id"`
	OwnerTypeID   string      `json:"owner_type_id"`
	Lat           string      `json:"map_lat"`
	Lon           string      `json:"map_lon"`
	Comment       string      `json:"comment"`
	StreetAddress string      `json:"street_address"`
	CreatedAt     string      `json:"order_date"`
	Rooms         string      `json:"rooms"`
	Bedrooms      string      `json:"bedrooms"`
	Price         string      `json:"price"`
	Pathway       interface{} `json:"pathway_json"`
}


type Response struct {
	Data struct {
		Prs []MyHomeData `json:"Prs"`
	} `json:"Data"`
}

type DataNormalized interface {
	MyHomeData | interface{}
}
