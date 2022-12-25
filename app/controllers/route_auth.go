package controllers

import (
	"log"
	"myapp_v1/app/models"
	"net/http"
)

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := user_session(w, r)
		if err != nil {
			// 初期バリデーション格納
			var emptyUser models.User
			data := make(map[string]interface{})
			data["validation"] = emptyUser

			tmpda := models.TemplateData{
				Form: *models.New(nil),
				Data: data,
			}

			// アカウント登録ページ作成
			user := models.User{
				Name:         r.PostFormValue("name"),
				Password:     r.PostFormValue("password"),
				TemplateData: tmpda,
			}
			generateHTML(w, user, "layout", "top_navbar", "signup")
		} else {
			http.Redirect(w, r, "/topics", 302)
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}

		// バリデーションチェック
		form := models.New(r.PostForm)
		form.Required("name", "password")
		form.MinLength("password", 4, r)
		users, err := models.GetUsers()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		for _, user := range users {
			form.NotSameName("name", user)
		}

		validation := models.User{
			Name:     r.PostFormValue("name"),
			Password: r.PostFormValue("password"),
		}

		if !form.Valid() {
			data := make(map[string]interface{})
			data["validation"] = validation
			tmpda := models.TemplateData{
				Form: *form,
				Data: data,
			}
			user := models.User{
				Name:         r.PostFormValue("name"),
				Password:     r.PostFormValue("password"),
				TemplateData: tmpda,
			}
			generateHTML(w, user, "layout", "top_navbar", "signup")

		} else {
			// ユーザー作成
			user := models.User{
				Name:     r.PostFormValue("name"),
				Password: r.PostFormValue("password"),
			}
			if err := user.CreateUser(); err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}
			http.Redirect(w, r, "/", 302)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// セッションチェック
	_, err := user_session(w, r)
	if err != nil {
		// 初期バリデーション格納
		var emptyUser models.User
		data := make(map[string]interface{})
		data["validation"] = emptyUser

		tmpda := models.TemplateData{
			Form: *models.New(nil),
			Data: data,
		}

		// ログインページ作成
		user := models.User{
			Name:     r.PostFormValue("name"),
			Password:     r.PostFormValue("password"),
			TemplateData: tmpda,
		}
		generateHTML(w, user, "layout", "top_navbar", "login")
	} else {
		http.Redirect(w, r, "/topics", 302)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/login", 302)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/404", 302)
		}
		// バリデーションチェック
		user, _ := models.GetUserByName(r.PostFormValue("name"))
		form := models.New(r.PostForm)
		form.Required("name", "password")
		form.MinLength("password", 4, r)
		form.SameNameAndPassword("name", "password", user)

		validation := models.User{
			Name: r.PostFormValue("name"),
			Password: r.PostFormValue("password"),
		}

		if !form.Valid() {
			data := make(map[string]interface{})
			data["validation"] = validation
			tmpda := models.TemplateData{
				Form: *form,
				Data: data,
			}
			user := models.User{
				Password:     r.PostFormValue("password"),
				TemplateData: tmpda,
			}
			generateHTML(w, user, "layout", "top_navbar", "login")
		}

		if user.Password  == models.Encrypt(r.PostFormValue("password")) {
			// セッション作成
			user_session, err := user.CreateUser_session()
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/404", 302)
			}

			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    user_session.UUID,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/404", 302)
	}

	// セッション削除
	if err != http.ErrNoCookie {
		user_session := models.User_Session{UUID: cookie.Value}
		user_session.Deleteuser_sessionByUUID()
	}
	http.Redirect(w, r, "/", 302)

}
