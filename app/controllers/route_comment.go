package controllers

import (
	"log"
	"myapp_v1/app/models"
	"net/http"
	"strconv"
)

func commentCreate(w http.ResponseWriter, r *http.Request, id int) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		// 初期バリデーション格納
		data := make(map[string]interface{})
		data["validation"] = ""

		tmpda := models.TemplateData{
			Form: *models.New(nil),
			Data: data,
		}

		// コメントページ作成
		topic, err := models.GetTopic(id)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		topic.TemplateData = tmpda
		topic.User = user

		if r.PostFormValue("new") == "new" {

			// 新しい順に並び変える
			comments, new, err := topic.GetComments_NewByTopic()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			topic.Comments = comments
			topic.Sort = new
			generateHTML(w, topic, "layout", "comment_navbar", "comment")
		} else {

			// 古い順に並び替える
			comments, old, err := topic.GetComments_OldByTopic()
			if err != nil {
				log.Panicln(err)
				http.Redirect(w, r, "/404", 302)
			}
			topic.Comments = comments
			topic.Sort = old
			generateHTML(w, topic, "layout", "comment_navbar", "comment")
		}
	}
}

func commentSave(w http.ResponseWriter, r *http.Request, id int) {
	// セッション確認
	sess, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		if r.Method == "GET" {

			i := strconv.Itoa(id)
			http.Redirect(w, r, "/topics/comment/"+i, 302)

		} else if r.Method == "POST" {

			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			// バリデーションチェック
			form := models.New(r.PostForm)
			form.Required("comment")
			form.MaxLength("comment", 50, r)

			topic, _ := models.GetTopic(id)
			comment := r.PostFormValue("comment")
			NowOrOld := r.PostFormValue("NewOrOld")

			if !form.Valid() {
				data := make(map[string]interface{})
				data["validation"] = comment
				tmpda := models.TemplateData{
					Form: *form,
					Data: data,
				}
				topic.TemplateData = tmpda
				if NowOrOld == "New" {
					comments, _, _ := topic.GetComments_NewByTopic()
					topic.Comments = comments
					generateHTML(w, topic, "layout", "comment_navbar", "comment")
				} else {
					comments, _, err := topic.GetComments_OldByTopic()
					if err != nil {
						log.Println(err)
						http.Redirect(w, r, "/404", 302)
					}
					topic.Comments = comments
					generateHTML(w, topic, "layout", "comment_navbar", "comment")

				}

			} else {
				user, err := sess.GetUserBySession()
				if err != nil {
					log.Println(err)
					http.Redirect(w, r, "/404", 302)
				}
				// コメント作成
				content := r.PostFormValue("comment")
				user.CreateComment(user.Name, content, id)

				i := strconv.Itoa(id)
				http.Redirect(w, r, "/topics/comment/"+i, 302)
			}
		}
	}
}
