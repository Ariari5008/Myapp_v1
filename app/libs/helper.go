package libs

import (
	"myapp_v1/app/models"
)

// トピックまとめ
func TopicAll(id int, name string, content string, user_id int, likes int, del_flag int) (m *models.Topic) {
	t := &models.Topic{
		ID:        id,
		Name:      name,
		Content:   content,
		UserID:    user_id,
		Likes:     likes,
		Del_flag:  del_flag,
	}
	m = t
	return m
}
