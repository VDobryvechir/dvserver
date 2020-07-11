package stock

type PrimaryCategory struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	OrderNumber   int                `json:"orderNumber"`
	SubCategories []*PrimaryCategory `json:"subCategories"`
}

type TopOffering struct {
	DefaultCategoryId    string             `json:"defaultCategoryId"`
	DefaultSubCategoryId string             `json:"defaultSubCategoryId"`
	TopCategories        []*PrimaryCategory `json:"topCategories"`
}

type BaseDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type BasePrice struct {
	Value             float64  `json:"value"`
	Currency          string   `json:"currency"`
	ChargePeriodicity *BaseDto `json:"changePeriodicity"`
}

type Offering struct {
	Id                     string       `json:"id"`
	Name                   string       `json:"name"`
	Description            string       `json:"description"`
	BasePrices             []*BasePrice `json:"basePrices"`
	AddressUnitEligibility bool         `json:"addressUnitEligibility"`
}

type InOutConfig struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type CommandPattern struct {
	SessionId   string `json:"sessionId"`
	CommandName string `json:"commandName"`
	OfferingId  string `json:"offeringId"`
	Quantity    int    `json:"quantity"`
}

type OfferingCard struct {
	Offering *Offering       `json:"offering"`
	Command  *CommandPattern `json:"command"`
}
