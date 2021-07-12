package dao

import (
	"gorm.io/gorm"
	"webOrder/model"
)

type DaoOrder interface {
	NewOrder(order *model.Order) error
	UpdateOrderByOrderNo(no string, data map[string]interface{}) error
	GetAllOrder(orders []*model.Order) error
	QueryOrderByNo(no string) (req *model.Order, err error)
	QueryOrderList(keyWord string, up bool) (orders []*model.Order, err error)
	QueryLikeOrder(keyWord string) (orders []*model.Order, err error)
}
type orderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *orderDAO {
	return &orderDAO{db: db}
}

//新增一条订单记录
func (o *orderDAO) NewOrder(order *model.Order) error {
	return o.db.Create(order).Error
}

//凭借订单号更新细节
func (o *orderDAO) UpdateOrderByOrderNo(no string, data map[string]interface{}) error {
	return o.db.Model(&model.Order{}).Where("order_no=?", no).Updates(data).Error
}

//获取所有订单
func (o *orderDAO) GetAllOrder(orders []*model.Order) error {
	return o.db.Find(&orders).Error
}

//凭借订单号查询一条订单
func (o *orderDAO) QueryOrderByNo(no string) (req *model.Order, err error) {
	err = o.db.Where("order_no=?", no).Last(&req).Error
	return
}

//凭借username模糊查询
func (o *orderDAO) QueryLikeOrder(keyWord string) (orders []*model.Order, err error) {
	err = o.db.Where("user_name LIKE ?", "%"+keyWord+"%").Find(&orders).Error
	return
}

//凭借关键词对所有订单排序
func (o *orderDAO) QueryOrderList(keyWord string, up bool) (orders []*model.Order, err error) {
	if up == true {
		Order := "desc"
		err = o.db.Order(keyWord + " " + Order).Find(&orders).Error
		return
	} else {
		Order := "asc"
		err = o.db.Order(keyWord + " " + Order).Find(&orders).Error
		return
	}
}