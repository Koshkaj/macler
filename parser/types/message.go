package types

type EventMessageFilter struct {
	PriceMin         int    `json:"price_min" validate:"required,gte=9,lte=999999"`
	PriceMax         int    `json:"price_max" validate:"required,gte=10,lte=1000000"`
	SquareMin        int    `json:"square_min" validate:"required,gte=9,lte=443"`
	SquareMax        int    `json:"square_max" validate:"required,gte=10,lte=444"`
	ADType           string `json:"ad_type" validate:"required,oneof=იყიდება ქირავდება"`
	PropertyType     string `json:"property_type" validate:"required,oneof=ბინა სახლი"`
	Keyword          string `json:"keyword" validate:"required"`
	City             string `json:"city" validate:"required,oneof=თბილისი გუდაური ბათუმი"`
	EmergencyListing bool   `json:"emergency_listing"`
}

type IncommingEventMessage struct {
	Phone          string             `json:"phone" validate:"required"`
	Amount         int                `json:"amount" validate:"required,gte=5,lte=50"`
	Interval       string             `json:"interval" validate:"required,interval"`
	TimeRangeWeeks int                `json:"time_range_weeks" validate:"required,gte=1,lte=12"`
	Filters        EventMessageFilter `json:"filters" validate:"required"`
}
