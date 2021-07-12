package dao

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"webOrder/dao/config"
	"webOrder/model"

	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	suite.Suite              //工具
	order       *model.Order //订单模型
	orders      []*model.Order
	dao         DaoOrder //接口
}

func (s *OrderTestSuite) SetupSuite() {
	s.T().Log("SetupSuite")
	s.dao = NewOrderDAO(config.NewOrm())

	s.order = &model.Order{
		OrderNo:  fmt.Sprintf("%d", time.Now().UnixNano()),
		UserName: "pual",
		Amount:   1000,
		Status:   false,
		FileUrl:  "FileUrl",
	}
	s.orders = []*model.Order{
		{
			OrderNo:  fmt.Sprintf("%d", time.Now().UnixNano()),
			UserName: "kaer",
			Amount:   1000,
			Status:   false,
			FileUrl:  "FileUrl",
		},
		{
			OrderNo:  fmt.Sprintf("%d", time.Now().UnixNano()),
			UserName: "dare",
			Amount:   1000,
			Status:   false,
			FileUrl:  "FileUrl",
		},
	}
}

func (s *OrderTestSuite) SetupTest() {
	s.T().Log("SetupTest")
}

func (s *OrderTestSuite) TestNewOrder(t *testing.T) {
	s.T().Log("TestAddNewOrder")
	require.NoError(s.T(), s.dao.NewOrder(s.order))
}

func (s *OrderTestSuite) TestUpdateOrderByOrderNo(t *testing.T) {
	s.T().Log("TestUpdateOrderByOrderNo")
	require.NoError(s.T(), s.dao.UpdateOrderByOrderNo(s.order.OrderNo, map[string]interface{}{
		"Amount":  s.order.Amount, //amount
		"Status":  s.order.Status,
		"FileUrl": s.order.FileUrl,
	}))
}

func (s *OrderTestSuite) TestGetAllOrder(t *testing.T) {
	s.T().Log("TestGetAllOrder")
	require.NoError(s.T(), s.dao.GetAllOrder(s.orders))
}

func (s *OrderTestSuite) TestQueryOrderByNo(t *testing.T) {
	s.T().Log("TestQueryOrderByNo")
	result, err := s.dao.QueryOrderByNo(s.order.OrderNo)
	require.NoError(s.T(), err, "数据库订单号查询失败")
	require.NotNil(s.T(), result, "未获取到结果")
}

func (s *OrderTestSuite) TestQueryLikeOrder(t *testing.T) {
	s.T().Log("TestQueryLikeOrder")
	var keyWords = []string{
		"12321", "123", "nooguy",
	}
	for _, keyWord := range keyWords {
		result, err := s.dao.QueryLikeOrder(keyWord)
		require.NoError(s.T(), err, "数据库模糊查询用户名失败")
		require.NotNil(s.T(), result, "未获取到结果")

	}
}

func (s *OrderTestSuite) TestQueryOrderList(t *testing.T) {
	s.T().Log("TestQueryOrderList")
	up := true
	var keyWords = []string{
		"12321", "123", "nooguy",
	}
	for _, keyWord := range keyWords {
		result, err := s.dao.QueryOrderList(keyWord, up)
		require.NoError(s.T(), err, "凭借关键词排序数据库失败")
		require.NotNil(s.T(), result, "未获取到结果")

	}
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}