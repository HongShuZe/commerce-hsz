package datamodels

type User struct {
	ID       int    `json:"id" sql:"id" form:"id"`
	NickName string `json:"nick_name" sql:"nick_name" form:"nick_name"`
	UserName string `json:"user_name" sql:"user_name" form:"user_name"`
	Password string `json:"password" sql:"password" form:"password"`
	Balance  string `json:"balance" sql:"user_balance" form:"user_balance"`
}
