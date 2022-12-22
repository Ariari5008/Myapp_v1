package models

import (
	"log"
	"time"
)

type Comment struct {
	ID        int
	Name      string
	Content   string
	UserID    int
	TopicID   int
	CreatedAt time.Time
}

// コメント作成
func (u *User) CreateComment(name string, content string, id int) (err error) {
	cmd := `insert into comments (
		name,
		content,
		user_id,
		topic_id, 
		created_at) values (?, ?, ?, ?, ?)`

	_, err = Db.Exec(cmd, name, content, u.ID, id, time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// 古い順に並べ替えてコメント取得
func (t *Topic) GetComments_OldByTopic() (comments []Comment, Sort string, err error) {
	cmd := `select id, name, content, user_id, topic_id, created_at from comments where topic_id = ? order by created_at asc`

	rows, err := Db.Query(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.ID,
			&comment.Name,
			&comment.Content,
			&comment.UserID,
			&comment.TopicID,
			&comment.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		comments = append(comments, comment)
	}
	rows.Close()

	return comments, "Old", err
}

// 新しい順に並べ替えてコメント取得
func (t *Topic) GetComments_NewByTopic() (comments []Comment, Sort string, err error) {
	cmd := `select id, name, content, user_id, topic_id, created_at from comments where topic_id = ? order by created_at desc `
	rows, err := Db.Query(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.ID,
			&comment.Name,
			&comment.Content,
			&comment.UserID,
			&comment.TopicID,
			&comment.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		comments = append(comments, comment)
	}
	rows.Close()

	return comments, "New", err
}

func (c *Comment) DeleteComment() error {
	cmd := `delete from comments where id = ?`
	_, err = Db.Exec(cmd, c.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
