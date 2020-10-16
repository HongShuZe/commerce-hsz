package services

import (
	"commerce-hsz/datamodels"
	"commerce-hsz/repositories"
)

type IProductService interface {
	GetProduct(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProduct(int64) bool
	InsertProduct(product *datamodels.Product)(int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(int64) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

// 初始化函数
func NewProductService(product repositories.IProduct) IProductService {
	return &ProductService{productRepository:product}
}

func (s *ProductService)GetProduct(id int64) (*datamodels.Product, error) {
	return s.productRepository.SelectOne(id)
}

func (s *ProductService)GetAllProduct() ([]*datamodels.Product, error){
	return s.productRepository.SelectAll()
}

func (s *ProductService)DeleteProduct(id int64) bool{
	return s.productRepository.Delete(id)
}

func (s *ProductService)InsertProduct(product *datamodels.Product)(int64, error){
	return s.productRepository.Insert(product)
}

func (s *ProductService)UpdateProduct(product *datamodels.Product) error{
	return s.productRepository.Update(product)
}

func (s *ProductService)SubNumberOne(productID int64) error  {
	return s.productRepository.SubProductNum(productID)
}
