package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shuSemester/model"
	"shuSemester/service/token"
	"time"
)

func getSemesterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dateString := r.URL.Query().Get("date")
	var date time.Time
	var err error = nil
	if dateString == "now" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateString)
	}
	if err != nil {
		w.WriteHeader(400)
		return
	}
	semester, err := model.GetByDate(date)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	data, _ := json.Marshal(semester)
	_, _ = w.Write(data)
}

func setSemesterHandler(w http.ResponseWriter, r *http.Request) {
	tokenInHeader := r.Header.Get("Authorization")
	if tokenInHeader == "" {
		w.WriteHeader(401)
		return
	}
	tokenString := tokenInHeader[7:]
	if !token.ValidateToken(tokenString) {
		w.WriteHeader(403)
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	var input model.Semester
	_ = json.Unmarshal(body, &input)
	model.Save(input)
}

func SemesterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getSemesterHandler(w, r)
	case "POST":
		setSemesterHandler(w, r)
	}
}

func PingPongHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
