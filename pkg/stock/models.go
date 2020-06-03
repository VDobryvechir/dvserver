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

type InOutConfig struct {
	In  string `json:"in"`
	Out string `json:"out"`
}
