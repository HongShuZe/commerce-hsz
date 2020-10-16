package services

import (
	"commerce-hsz/datamodels"
	"commerce-hsz/repositories"
)

type IOrderService interface {
	GetOrder(int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	DeleteOrder(int64) bool
	InsertOrder(*datamodels.Order)(int64, error)
	UpdateOrder(*datamodels.Order) error
	// 查询订单相关信息
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

// 初始化函数
func NewOrderService(orderRepository repositories.IOrderRepository) IOrderService {
	return &OrderService{orderRepository:orderRepository}
}

func (o *OrderService)GetOrder(id int64) (*datamodels.Order, error) {
	return o.orderRepository.SelectOne(id)
}

func (o *OrderService)GetAllOrder() ([]*datamodels.Order, error){
	return o.orderRepository.SelectAll()
}

func (o *OrderService)DeleteOrder(id int64) bool{
	return o.orderRepository.Delete(id)
}

func (o *OrderService)InsertOrder(order *datamodels.Order)(int64, error){
	return o.orderRepository.Insert(order)
}

func (o *OrderService)UpdateOrder(order *datamodels.Order) error{
	return o.orderRepository.Update(order)
}

func (o *OrderService)GetAllOrderInfo() (map[int]map[string]string, error)  {
	return o.orderRepository.SelectAllWithInfo()
}

// 根据消息创建订单
func (o *OrderService)InsertOrderByMessage(message *datamodels.Message) (int64, error)  {
	order := &datamodels.Order{
		UserID: int(message.UserID),
		ProductID: int(message.ProductID),
		OrderNum: 1,
	}
	return o.InsertOrder(order)
}
