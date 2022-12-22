package controllers

import (
	"fmt"
	"html/template"
	"myapp_v1/app/models"
	"myapp_v1/config"
	"net/http"
	"regexp"
	"strconv"
)

// テンプレート作成
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// セッション呼び出し
func user_session(w http.ResponseWriter, r *http.Request) (sess models.User_Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.User_Session{UUID: cookie.Value}
		if ok, _ := sess.Checkuser_session(); !ok {
			err = fmt.Errorf("Invalid user_session")
		}
	}
	return sess, err
}


var validPath = regexp.MustCompile("^/topics/(edit|update|delete|comment|comment_save|comment_delete)/([0-9]+)$")

// IDパス作成
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.Redirect(w, r, "/404", 302)
			return
		}
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.Redirect(w, r, "/404", 302)
			return
		}

		fn(w, r, qi)
	}
}

// サーバー
func StartMainServer() error {

	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/404", NotFound)
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/topics", index)
	http.HandleFunc("/browse", browse)
	http.HandleFunc("/browse/increment", Increment)
	http.HandleFunc("/ranking", ranking)
	http.HandleFunc("/topics/new", topicNew)
	http.HandleFunc("/topics/topic_save", topicSave)
	http.HandleFunc("/topics/edit/", parseURL(topicEdit))
	http.HandleFunc("/topics/update/", parseURL(topicUpdate))
	http.HandleFunc("/topics/delete/", parseURL(topicDelete))
	http.HandleFunc("/topics/comment/", parseURL(commentCreate))
	http.HandleFunc("/topics/comment_save/", parseURL(commentSave))
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
