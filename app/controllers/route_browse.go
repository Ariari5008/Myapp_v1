package controllers

import (
	"log"
	"myapp_v1/app/libs"
	"myapp_v1/app/models"
	"net/http"
	"strconv"
)

func browse(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	_, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

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
		topics, err := models.GetTopicsByDel_flag()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		generateHTML(w, topics, "layout", "browse_navbar", "browse")

	}
}

func Increment(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		topic_id := r.PostFormValue("topic_id")
		id, err := strconv.Atoi(topic_id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		topic, err := models.GetTopic(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		// いいね数カウント
		err = topic.IncrementTopicLikes()
		if err != nil {
			log.Fatalln(err)
			http.Redirect(w, r, "/404", 302)
		}

		http.Redirect(w, r, "/browse", 302)
	}
}
