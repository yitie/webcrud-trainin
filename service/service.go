package service

import (
	"webOrder/dao"
	"webOrder/model"
)

type OrderService struct {
	orderDao dao.DaoOrder
}

func NewOrderService(orderDao dao.DaoOrder) *OrderService {
	return &OrderService{orderDao: orderDao}
}

func (o *OrderService) NewOrder(order *model.Order) error {
	return o.orderDao.NewOrder(order)
}

func (o *OrderService) UpdateOrderByOrderNo(no string, data map[string]interface{}) error {
	return o.orderDao.UpdateOrderByOrderNo(no, data)
}

func (o *OrderService) QueryOrderByNo(no string) (req *model.Order, err error) {
	return o.orderDao.QueryOrderByNo(no)
}

func (o *OrderService) QueryLikeOrder(keyWord string) (orders []*model.Order, err error) {
	return o.orderDao.QueryLikeOrder(keyWord)

}

func (o *OrderService) QueryOrderList(keyWord string, up bool) (orders []*model.Order, err error) {
	return o.orderDao.QueryOrderList(keyWord, up)
}

func (o *OrderService) GetAllOrder(orders []*model.Order) error {
	return o.orderDao.GetAllOrder(orders)
}