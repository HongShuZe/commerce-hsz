package repositories

import (
	"commerce-hsz/datamodels"
	"database/sql"
	"commerce-hsz/common"
	"errors"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (*datamodels.User, error)
	Insert(user *datamodels.User) (int64, error)
}

type UserManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func NewUserRepositry(table string, db *sql.DB) IUserRepository{
	return &UserManagerRepository{table:table, mysqlConn:db}
}

func (u *UserManagerRepository)Conn() error {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}

	if u.table == "" {
		u.table = "tbl_user"
	}
	return nil
}

// // 根据username查询用户
func (u *UserManagerRepository)Select(userName string) (*datamodels.User, error){
	if userName == "" {
		return &datamodels.User{}, errors.New("用户名为空")
	}

	err := u.Conn()
	if err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from tbl_user where user_name=? and status=0"
	rows, err := u.mysqlConn.Query(sql, userName)
	if err != nil {
		return &datamodels.User{}, err
	}
	defer rows.Close()

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}

	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

// 新增用户
func (u *UserManagerRepository)Insert(user *datamodels.User) (int64, error){
	err := u.Conn()
	if err != nil {
		return 0, err
	}

	sql := "insert into tbl_user (`nick_name`,`user_name`,`password`,`user_balance`,`status`) values (?,?,?,?,0)"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.NickName, user.UserName, user.Password, 1000)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// 根据id查询用户
func (u *UserManagerRepository)SelectByID(id int) (*datamodels.User, error) {
	err := u.Conn()
	if err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from tbl_user where id=? and status=0"
	row, err := u.mysqlConn.Query(sql, id)
	if err != nil {
		return &datamodels.User{}, err
	}
	defer row.Close()

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, err
	}

	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

