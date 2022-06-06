package controllers

import (
	"errors"
	"fmt"
	"goTest/app/config"
	"goTest/app/models"
	"text/template"

	"net/http"
)

func StartMainServer() error {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/recruitment", recruitment)
	http.HandleFunc("/profile", profile)
	http.HandleFunc("/search", search)
	http.HandleFunc("/detail", detail)
	http.HandleFunc("/recruitment/delete", recruitmentDelete)

	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/profile/edit/name", nameEdit)
	http.HandleFunc("/profile/edit/email", emailEdit)
	http.HandleFunc("/profile/edit/password", passwordEdit)

	// port := os.Getenv("PORT")
	// return http.ListenAndServe(":"+port, nil)
	return http.ListenAndServe(":"+config.Config.Port, nil)

}

func GenerateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)

}

func sessionU(writer http.ResponseWriter, request *http.Request) (sess models.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}
