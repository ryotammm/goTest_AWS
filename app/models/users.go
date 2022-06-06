package models

import (
	"time"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
}

type Session struct {
	ID       int
	UUID     string
	Name     string
	Email    string
	UserID   int
	CreateAt time.Time
}

func (u *User) CreateUser() (err error) {
	// user = User{}

	//データの挿入
	stmt, err := Db.Prepare("insert users SET uuid=?,name=?,email=?,password=?,created_at=?")
	checkErr(err)
	_, err = stmt.Exec(createUUID(), u.Name, u.Email, Encrypt(u.PassWord), time.Now())
	checkErr(err)

	return err
}

func GetUser(id int) (user User, err error) {
	user = User{}
	stmt, err := Db.Prepare("select id, uuid,name,email,password,created_at from users where id = ?")

	checkErr(err)
	row, err := stmt.Query(id)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.PassWord, &user.CreatedAt)
	checkErr(err)

	return user, err
}

func (u *User) UpdateUser() (err error) {
	stmt, err := Db.Prepare("update users set name = ?, email = ? where id =?")
	checkErr(err)

	_, err = stmt.Exec(u.Name, u.Email, u.ID)
	checkErr(err)

	return err
}

func (u *User) DeleteUser() (err error) {

	stmt, err := Db.Prepare("delete  from users where id = ?")
	checkErr(err)
	_, err = stmt.Exec(u.ID)

	checkErr(err)

	return err
}

func GetUserByEmail(email string) (user User, err error) {

	user = User{}
	stmt, err := Db.Prepare("select id ,uuid,name,email,password,created_at from users where email = ?")

	checkErr(err)
	row, err := stmt.Query(email)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.PassWord, &user.CreatedAt)
	checkErr(err)

	return user, err
}

//session作成
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	stmt, err := Db.Prepare("insert  sessions SET uuid=?,name=?,email=?,user_id=?,created_at=?")
	checkErr(err)
	uuid := createUUID()
	_, err = stmt.Exec(uuid, u.Name, u.Email, u.ID, time.Now())
	checkErr(err)

	stmt, err = Db.Prepare("select id, uuid, name,email, user_id, created_at from sessions where uuid = ? and user_id = ? and email = ?")

	checkErr(err)
	row, err := stmt.Query(uuid, u.ID, u.Email)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&session.ID, &session.UUID, &session.Name, &session.Email, &session.UserID, &session.CreateAt)
	checkErr(err)

	return session, err

}

func (sess *Session) CheckSession() (valid bool, err error) {

	stmt, err := Db.Prepare("select id, uuid,name, email,user_id, created_at from sessions where uuid = ? ")

	checkErr(err)
	row, err := stmt.Query(sess.UUID)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&sess.ID, &sess.UUID, &sess.Name, &sess.Email, &sess.UserID, &sess.CreateAt)

	if err != nil {
		valid = false
		return
	}
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

func (sess *Session) DeleteSessionByUUID() (err error) {
	stmt, err := Db.Prepare("delete from sessions where uuid = ?")
	checkErr(err)
	_, err = stmt.Exec(sess.UUID)

	checkErr(err)

	return err
}

func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	stmt, err := Db.Prepare("select id, uuid,name,email,created_at FROM users where id = ? ")

	checkErr(err)
	row, err := stmt.Query(sess.UserID)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt)

	return user, err

}

/**/
func GetUserByEmailAndPassWord(email, password string) (user User, err error) {

	user = User{}
	stmt, err := Db.Prepare("select id , uuid , name , email, password, created_at from users where email = ?  and password = ?")

	checkErr(err)
	row, err := stmt.Query(email, password)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.PassWord, &user.CreatedAt)

	return user, err

}

func NameEdit(id int, name string) error {
	stmt, err := Db.Prepare("update users set name = ? where id = ?")
	checkErr(err)

	_, err = stmt.Exec(name, id)
	checkErr(err)

	return err
}

// func (sess *Session) NameEditSession(id int, name string) (err error) {

// 	cmd := `update sessions set name = $1  where id = $2`
// 	_, err = Db.Exec(cmd, name, id)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return err
// }

func NameEditSession(id int, name string) (err error) {
	stmt, err := Db.Prepare("update sessions set name = ?  where user_id = ?")
	checkErr(err)

	_, err = stmt.Exec(name, id)
	checkErr(err)

	return err

}

func EmailEdit(id int, email string) error {
	stmt, err := Db.Prepare("update users set email = ? where id = ?")
	checkErr(err)

	_, err = stmt.Exec(email, id)
	checkErr(err)

	return err
}

func (sess *Session) EmailEditSession(id int, email string) (err error) {
	stmt, err := Db.Prepare("update sessions set email = ?  where id = ?")
	checkErr(err)

	_, err = stmt.Exec(email, id)
	checkErr(err)
	return err
}

func CheckPassWord(id int, pass string) (user User, err error) {
	user = User{}
	stmt, err := Db.Prepare("select id, password  from users where id = ? and password = ?")

	checkErr(err)
	row, err := stmt.Query(id, pass)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.ID, &user.PassWord)

	return user, err

}

func PasswordEdit(id int, pass string) (err error) {
	stmt, err := Db.Prepare("update users set password = ? where id = ?")
	checkErr(err)

	_, err = stmt.Exec(pass, id)
	checkErr(err)

	return err
}

//Email存在チェック
func ExistEmail(email string) (user User, err error) {
	stmt, err := Db.Prepare("select  email  from users where email = ?")

	checkErr(err)
	row, err := stmt.Query(email)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&user.Email)

	return user, err
}
