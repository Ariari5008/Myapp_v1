package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"myapp_v1/config"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

var err error

const (
	tabalNameUser         = "users"
	tableNameUser_session = "user_sessions"
	tableNameTopic        = "topics"
	tableNameLikeUser     = "like_user"
	tableNameComment      = "comments"
)

func init() {
	Db, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE 
		IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		password STRING,
		created_at DATETIME)`, tabalNameUser)

	Db.Exec(cmdU)

	cmdT := fmt.Sprintf(`CREATE TABLE 
		IF NOT EXISTS %s(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name STRING,
			content TEXT,
			user_id INTEGER,
			likes INTEGER,
			del_flag INTEGER NOT NULL,
			created_at DATETIME)`, tableNameTopic)

	Db.Exec(cmdT)

	cmdS := fmt.Sprintf(`CREATE TABLE 
		IF NOT EXISTS %s(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid STRING NOT NULL UNIQUE,
			user_id INTEGER,
			created_at DATETIME)`, tableNameUser_session)

	Db.Exec(cmdS)

	cmdC := fmt.Sprintf(`CREATE TABLE 
		IF NOT EXISTS %s(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name STRING,
			content STRING,
			user_id INTEGER,
			topic_id INTEGER,
			created_at DATETIME)`, tableNameComment)

	Db.Exec(cmdC)

}

// UUID作成
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

// ハッシュ値パスワード作成
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
