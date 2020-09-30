package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"id" goods:"id"`
	ProductName  string `json:"productName" sql:"productName" goods:"productName"`
	ProductNum   string `json:"productNum" sql:"productNum" goods:"productNum"`
	ProductImage string `json:"productImage" sql:"productImage" goods:"productImage"`
	ProductUrl   string `json:"productUrl" sql:"productUrl" goods:"productUrl"`
}
