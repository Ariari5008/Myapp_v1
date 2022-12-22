package models

import (
	"log"
	"time"
)

type Topic struct {
	ID              int
	Name            string
	Content         string
	UserID          int
	Likes           int
	Del_flag        int
	CreatedAt       time.Time
	Comments        []Comment
	User            User
	TemplateData    TemplateData
	Sort            string
	CommentDel_flag bool
}

// トピック作成
func (u *User) CreateTopic(content string, likes int, flag int) (err error) {
	cmd := `insert into topics (
		name,
		content,
		user_id, 
		likes,
		del_flag,
		created_at) values (?, ?, ?, ?, ?, ?)`

	_, err = Db.Exec(cmd, u.Name, content, u.ID, likes, flag, time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// トピック取得
func GetTopic(id int) (topic Topic, err error) {
	cmd := `select id, name, content, user_id, likes, del_flag, created_at from topics where id = ?`
	topic = Topic{}

	err = Db.QueryRow(cmd, id).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Content,
		&topic.UserID,
		&topic.Likes,
		&topic.Del_flag,
		&topic.CreatedAt)

	return topic, err
}

// 全トピック取得
func GetTopics() (topics []Topic, err error) {
	cmd := `select id, name, content, user_id, likes, del_flag, created_at from topics
	order by created_at asc`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic Topic
		err = rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Content,
			&topic.UserID,
			&topic.Likes,
			&topic.Del_flag,
			&topic.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		topics = append(topics, topic)
	}
	rows.Close()

	return topics, err
}

// ユーザー一致のトピック取得
func (u *User) GetTopicsByUser() (topics []Topic, err error) {
	cmd := `select id, name, content, user_id, likes, del_flag, created_at from topics where user_id = ? order by created_at asc`

	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic Topic
		err = rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Content,
			&topic.UserID,
			&topic.Likes,
			&topic.Del_flag,
			&topic.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		topics = append(topics, topic)
	}
	rows.Close()

	return topics, err
}

// Del_flag一致のトピック取得
func GetTopicsByDel_flag() (topics []Topic, err error) {
	cmd := `select id, name, content, user_id, likes, del_flag, created_at from topics where del_flag = ? order by created_at desc`
	rows, err := Db.Query(cmd, 1)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var topic Topic
		err = rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Content,
			&topic.UserID,
			&topic.Likes,
			&topic.Del_flag,
			&topic.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		topics = append(topics, topic)
	}
	rows.Close()

	return topics, err
}

// トピック更新
func (t *Topic) UpdateTopic() error {
	cmd := `update topics set name = ?, content = ?, user_id = ?, likes = ?, del_flag = ? where id = ?`
	_, err := Db.Exec(cmd, t.Name, t.Content, t.UserID, t.Likes, t.Del_flag, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// トピック消去
func (t *Topic) DeleteTopic() error {
	cmd := `delete from topics where id = ?`
	_, err = Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// いいね数+1
func (t *Topic) IncrementTopicLikes() error {
	cmd := `update topics set name = ?, content = ?, user_id = ?, likes = ? + 1, del_flag = ? where id = ?`
	_, err := Db.Exec(cmd, t.Name, t.Content, t.UserID, t.Likes, t.Del_flag, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
