package datamodels

type Product struct {
	ID           int     `json:"id" sql:"id" goods:"id"`
	SellerID     int     `json:"seller_id" sql:"seller_id" goods:"seller_id"`
	ProductName  string  `json:"product_name" sql:"product_name" goods:"product_name"`
	ProductNum   int     `json:"product_num" sql:"product_num" goods:"product_num"`
	ProductImage string  `json:"product_image" sql:"product_image" goods:"product_image"`
	ProductUrl   string  `json:"product_url" sql:"product_url" goods:"product_url"`
	ProductPrice float64 `json:"product_price" sql:"product_price" goods:"product_price"`
	ProductInfo  string  `json:"product_info" sql:"product_info" goods:"product_info"`
}
