package controllers

import (
	"log"
	"myapp_v1/app/libs"
	"myapp_v1/app/models"
	"net/http"
	"strconv"
)


func top(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	_, err := user_session(w, r)
	if err != nil {

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		// トップページ作成
		flag := r.PostFormValue("flag")
		topic_id := r.PostFormValue("topic_id")
		id, _ := strconv.Atoi(topic_id)
		topic, _ := models.GetTopic(id)

		users, err := models.GetUsers()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		for _, v := range users {
			if flag == "公開" {
				t := libs.TopicAll(id, v.Name, topic.Content, v.ID, topic.Likes, 1)
				if topic.Name == v.Name {
					if err := t.UpdateTopic(); err != nil {
						log.Println(err)
						http.Redirect(w, r, "/404", 302)
					}
				}
				} else if flag == "非公開" {
					t := libs.TopicAll(id, v.Name, topic.Content, v.ID, topic.Likes, 0)
					if topic.Name == v.Name {
						if err := t.UpdateTopic(); err != nil {
							log.Println(err)
							http.Redirect(w, r, "/404", 302)
						}
				}
			}
		} 

		topics, _ := models.GetTopicsByDel_flag()
		generateHTML(w, topics, "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/topics", 302)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
		} else {
			user, err := sess.GetUserBySession()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			// ホームページ作成
			topics, _ := user.GetTopicsByUser()
			user.Topics = topics
			generateHTML(w, user, "layout", "private_navbar", "index")
		}
	}
	
	func NotFound(w http.ResponseWriter, r *http.Request) {
		// 404ページ作成
		generateHTML(w, nil, "layout", "404", "top_navbar")
	}




