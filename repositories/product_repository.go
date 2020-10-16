package repositories

import (
	"commerce-hsz/datamodels"
	"database/sql"
	"commerce-hsz/common"
	"strconv"
)

// 定义接口
type IProduct interface {
	Conn()(error)
	Insert(*datamodels.Product)(int64, error)
	Delete(int64)(bool)
	Update(*datamodels.Product)(error)
	SelectOne(int64)(*datamodels.Product, error)
	SelectAll()([]*datamodels.Product, error)
	SubProductNum(int64) error
}

// 实现类
type ProductManager struct {
	table string
	mysqlConn *sql.DB
}

// 初始化函数
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table:table, mysqlConn:db}
}

// 连接数据库
func (p *ProductManager)Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}

	if p.table == "" {
		p.table = "tbl_product"
	}
	return
}

// 新增商品
func (p *ProductManager)Insert(product *datamodels.Product)(int64,error) {
	// 判断连接是否存在
	err := p.Conn()
	if err != nil {
		return 0, err
	}
	// sql预编译
	sql := "insert ignore into tbl_product (`product_name`,`product_num`,`product_image`,`product_url`,`product_price`,`product_info`,`status`) values(?,?,?,?,?,?,0) "
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// 传入参数
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPrice, product.ProductInfo)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// 删除商品
func (p *ProductManager)Delete(id int64)(bool) {
	err := p.Conn()
	if err != nil {
		return false
	}

	sql := "update tbl_product set status=1 where id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(strconv.FormatInt(id, 10))
	if err != nil {
		return false
	}

	return true
}

// 更新商品信息
func (p *ProductManager)Update(product *datamodels.Product)(error) {
	err := p.Conn()
	if err != nil {
		return err
	}

	sql := "update tbl_product set product_name=?, product_num=?, product_image=?, product_url=?, product_price=?, product_info=? where id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPrice, product.ProductInfo,product.ID)
	if err != nil {
		return err
	}

	return nil
}

// 根据id查询商品
func (p *ProductManager)SelectOne(id int64)(*datamodels.Product ,error) {
	err := p.Conn()
	if err != nil {
		return &datamodels.Product{}, err
	}

	sql := "select * from tbl_product where status=0 and id=" + strconv.FormatInt(id, 10)
	stmt, err := p.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Product{}, err
	}
	defer stmt.Close()

	result := common.GetResultRow(stmt)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	productResult := &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return productResult, nil
}

// 查询全部商品
func (p *ProductManager)SelectAll()([]*datamodels.Product ,error) {
	err := p.Conn()
	if err != nil {
		return nil, err
	}

	sql := "select * from tbl_product where status=0"
	stmt, err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result := common.GetResultRows(stmt)
	if len(result) == 0 {
		return nil, nil
	}

	productArray := make([]*datamodels.Product, 0)
	for _, v := range result {
		productResult := &datamodels.Product{}
		common.DataToStructByTagSql(v, productResult)
		productArray = append(productArray, productResult)
	}

	return productArray, nil
}

func (p *ProductManager)SubProductNum(productID int64) error  {
	err := p.Conn()
	if err != nil {
		return err
	}

	sql := "update tbl_product set product_num=product_num-1 where id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(productID)
	if err != nil {
		return err
	}

	return nil
}




