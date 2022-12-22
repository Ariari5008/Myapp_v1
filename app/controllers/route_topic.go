package controllers

import (
	"fmt"
	"log"
	"myapp_v1/app/libs"
	"myapp_v1/app/models"
	"net/http"
	"strconv"
)


func topicNew(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	_, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// 初期バリデーション格納
		var emptyTopic models.Topic
		data := make(map[string]interface{})
		data["validation"] = emptyTopic

		tmpda := models.TemplateData{
			Form: *models.New(nil),
			Data: data,
		}

		// トピックページ作成
		topic := models.Topic{
			Content:      r.PostFormValue("content"),
			TemplateData: tmpda,
		}
		generateHTML(w, topic, "layout", "private_navbar", "topic_new")
	}
}

func topicSave(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		if r.Method == "GET" {
			http.Redirect(w, r, "/topics/new", 302)
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}

			// バリデーションチェック
			form := models.New(r.PostForm)
			form.Required("content")
			form.MaxLength("content", 100, r)

			validation := models.Topic{
				Content: r.PostFormValue("content"),
			}

			if !form.Valid() {
				data := make(map[string]interface{})
				data["validation"] = validation
				tmpda := models.TemplateData{
					Form: *form,
					Data: data,
				}
				topic := models.Topic{
					Content:      r.PostFormValue("content"),
					TemplateData: tmpda,
				}
				generateHTML(w, topic, "layout", "private_navbar", "topic_new")
			} else {
				user, err := sess.GetUserBySession()
				if err != nil {
					log.Println(err)
					http.Redirect(w, r, "/404", 302)
				}
				// トピック作成
				content := r.PostFormValue("content")
				if err := user.CreateTopic(content, 0, 0); err != nil {
					log.Println(err)
					http.Redirect(w, r, "/404", 302)
				}

				http.Redirect(w, r, "/topics", 302)
			}
		}
	}
}

func topicEdit(w http.ResponseWriter, r *http.Request, id int) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}	
		
		topic, err := models.GetTopic(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		// 初期バリデーション格納
		data := make(map[string]interface{})
		data["validation"] = topic
		
		tmpda := models.TemplateData{
			Form: *models.New(nil),
			Data: data,
		}
		
		topic.TemplateData = tmpda
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		// トピック更新ページ作成
		generateHTML(w, topic, "layout", "private_navbar", "topic_edit")
	}
}

func topicUpdate(w http.ResponseWriter, r *http.Request, id int) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		if r.Method == "GET" {
			i := strconv.Itoa(id)
			http.Redirect(w, r, "/topics/edit/"+i, 302)
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			// バリデーションチェック
			form := models.New(r.PostForm)
			form.Required("content")
			form.MaxLength("content", 100, r)

			topic, err := models.GetTopic(id)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			validation := models.Topic{
				Content:   r.PostFormValue("content"),
			}

			if !form.Valid() {
				data := make(map[string]interface{})
				data["validation"] = validation
				tmpda := models.TemplateData{
					Form: *form,
					Data: data,
				}
				topic.TemplateData = tmpda
				generateHTML(w, topic, "layout", "private_navbar", "topic_edit")
			} else {
				user, err := sess.GetUserBySession()
				if err != nil {
					log.Println(err)
					http.Redirect(w, r, "/404", 302)
				}
				// トピック更新
				name := user.Name
				content := r.PostFormValue("content")
				t := libs.TopicAll(id, name, content, user.ID, topic.Likes, topic.Del_flag)
				if err := t.UpdateTopic(); err != nil {
					log.Println(err)
					http.Redirect(w, r, "/404", 302)
				}
				http.Redirect(w, r, "/topics", 302)
			}
		}
	}
}

func topicDelete(w http.ResponseWriter, r *http.Request, id int) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		t, err := models.GetTopic(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		
		// トピック削除
		if user.ID == t.UserID {
			comments, _, err := t.GetComments_OldByTopic()
			for _, comment := range comments {
				if err := comment.DeleteComment; err != nil {
					fmt.Println(err)
					http.Redirect(w, r, "/404", 302)
				}
			}
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			
			if err := t.DeleteTopic(); err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
		}
	}
	http.Redirect(w, r, "/topics", 302)

}
