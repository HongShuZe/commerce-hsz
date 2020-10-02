package repositories

import (
	"commerce-hsz/datamodels"
	"database/sql"
	"commerce-hsz/common"
	"strconv"
)

// 定义接口
type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectOne(int64)(*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	// 查询和订单关联的消息
	SelectAllWithInfo() (map[int]map[string]string, error)
}

// 实现类
type OrderManager struct {
	table string
	mysqlConn *sql.DB
}

// 初始化函数
func NewOrderManager(table string, db *sql.DB) IOrderRepository {
	return &OrderManager{table: table, mysqlConn: db}
}

// 连接数据库
func (o *OrderManager)Conn() (err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}

	if o.table == "" {
		o.table = "tbl_order"
	}
	return
}

// 新建订单
func (o *OrderManager)Insert(order *datamodels.Order) (int64, error) {
	err := o.Conn()
	if err != nil {
		return 0, err
	}

	sql := "insert ignore into tbl_order (`product_id`,`user_id`,`order_num`,`total_price`,`status`) value (?,?,?,?,0)"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(order.ProductID, order.UserID, order.OrderNum, order.TotalPrice)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// 删除订单
// TODO 改为标记删除
func (o *OrderManager)Delete(id int64) bool {
	err := o.Conn()
	if err != nil {
		return false
	}

	sql := "delete from tbl_order where id=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err !=nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}

	return true
}

// 更新订单
func (o *OrderManager)Update(order *datamodels.Order) error {
	err := o.Conn()
	if err != nil {
		return err
	}

	sql := "update tbl_order set product_id=?, user_id=?, order_num=?, total_price=?, status=? where id=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ProductID, order.UserID, order.OrderNum, order.TotalPrice, order.Status, order.ID)
	if err != nil {
		return err
	}

	return err
}

// 查询单个订单
func (o *OrderManager)SelectOne(id int64)(*datamodels.Order, error) {
	err := o.Conn()
	if err != nil {
		return &datamodels.Order{}, err
	}

	sql := "select * from tbl_order where id="+strconv.FormatInt(id, 10)
	stmt, err := o.mysqlConn.Query(sql)
	if err !=nil {
		return &datamodels.Order{}, err
	}
	defer stmt.Close()

	result := common.GetResultRow(stmt)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	orderResult := &datamodels.Order{}
	common.DataToStructByTagSql(result, orderResult)
	return orderResult, err
}

// 查询全部订单
func (o *OrderManager)SelectAll() ([]*datamodels.Order, error) {
	err := o.Conn()
	if err != nil {
		return nil, err
	}

	sql := "select * from tbl_order"
	stmt, err := o.mysqlConn.Query(sql)
	if err !=nil {
		return nil, err
	}
	defer stmt.Close()

	result := common.GetResultRows(stmt)
	if len(result) == 0 {
		return nil, err
	}

	orderArray := make([]*datamodels.Order, 0)
	for _, v := range result{
		orderResult := &datamodels.Order{}
		common.DataToStructByTagSql(v, orderResult)
		orderArray = append(orderArray, orderResult)
	}

	return orderArray, err
}

// 查询全部和订单有关的信息
func (o *OrderManager)SelectAllWithInfo() (map[int]map[string]string, error)  {
	err := o.Conn()
	if err != nil {
		return nil, err
	}

	sql := "Select o.id,p.product_name,o.status From tbl_order as o left join tbl_product as p on o.product_id=p.id"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return common.GetResultRows(rows), err
}
