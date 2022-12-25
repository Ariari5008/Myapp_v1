package models

import (
	"log"
	"time"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Password  string
	CreatedAt time.Time
	Topics    []Topic
	TemplateData TemplateData
}

type User_Session struct {
	ID        int
	UUID      string
	UserID    int
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {
	cmd := `insert into users (
		uuid,
		name,
		password,
		created_at) values (?, ?, ?, ?)`

	_, err = Db.Exec(cmd,
		createUUID(),
		u.Name,
		Encrypt(u.Password),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetUsers() (users []User, err error) {
	cmd := `select id, uuid, name, password, created_at from users`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.UUID,
			&user.Name,
			&user.Password,
			&user.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		users = append(users, user)
	}
	rows.Close()

	return users, err
}

func GetUserByName(Name string) (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, password, created_at from users where name = ?`
	err = Db.QueryRow(cmd, Name).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
	)

	return user, err
}

func (u *User) CreateUser_session() (user_session User_Session, err error) {
	user_session = User_Session{}
	cmd1 := `insert into user_sessions (
		uuid, 
		user_id, 
		created_at) values (?, ?, ?)`

	_, err = Db.Exec(cmd1, createUUID(), u.ID, time.Now())

	if err != nil {
		log.Println(err)
	}

	cmd2 := `select id, uuid, user_id, created_at from user_sessions where user_id = ?`

	err = Db.QueryRow(cmd2, u.ID).Scan(
		&user_session.ID,
		&user_session.UUID,
		&user_session.UserID,
		&user_session.CreatedAt)

	return user_session, err
}

func (sess *User_Session) Checkuser_session() (valid bool, err error) {
	cmd := `select id, uuid, user_id, created_at from user_sessions where uuid = ?`

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.UserID,
		&sess.CreatedAt)

	if err != nil {
		valid = false
		return
	}
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

func (sess *User_Session) Deleteuser_sessionByUUID() (err error) {
	cmd := `delete from user_sessions where uuid = ?`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *User_Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, created_at from users where id = ?`
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.CreatedAt)

	return user, err
}
