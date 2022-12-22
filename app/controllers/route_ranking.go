package controllers

import (
	"log"
	"myapp_v1/app/models"
	"net/http"
)

func ranking(w http.ResponseWriter, r *http.Request) {
	// セッション確認
	_, err := user_session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		topics, err := models.GetTopics()
		if err != nil {
			log.Fatalln(err)
			http.Redirect(w, r, "/404", 302)
		}
		// ランキング作成
		r := &models.Ranking{}
		for _, v := range topics {
			r.AddRanking(v)
		}
		
		// ランキングページ作成
		generateHTML(w, r, "layout", "ranking_navbar", "ranking")
	}

}
