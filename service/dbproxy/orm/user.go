package orm

import (
	"log"

	mydb "filestore-server/service/dbproxy/conn"
)

// UserSignup : 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		log.Println("Failed to insert, err:" + err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		log.Println("Failed to insert, err:" + err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		res.Suc = true
		return
	}

	res.Suc = false
	res.Msg = "无记录更新"
	return
}

// UserSignin : 判断密码是否一致
func UserSignin(username string, encpwd string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		log.Println("username not found: " + username)
		res.Suc = false
		res.Msg = "用户名未注册"
		return
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		res.Suc = true
		res.Data = true
		return
	}
	res.Suc = false
	res.Msg = "用户名/密码不匹配"
	return
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// GetUserInfo : 查询用户信息
func GetUserInfo(username string) (res ExecResult) {
	user := TableUser{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		// error不为nil, 返回时user应当置为nil
		//return user, err
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = user
	return
}

// UserExist : 查询用户是否存在
func UserExist(username string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}
