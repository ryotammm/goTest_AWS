package controllers

import (
	"goTest/app/models"
	"log"
	"net/http"
	"unicode/utf8"
)

func signup(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		_, err := sessionU(w, r)
		if err != nil {
			GenerateHTML(w, nil, "layout", "signup", "public_navbar", "public_navbarMobile")
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		name := r.PostFormValue("name")

		nameC := utf8.RuneCountInString(name)

		if nameC > 10 {
			errorMessage := ErrorMessage{
				Err3: "お名前は10文字以下で入力してください",
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "signup", "public_navbar", "public_navbarMobile")
			return
		}

		result := MatchEmailString(r.PostFormValue("email"))

		if !result {
			errorMessage := ErrorMessage{
				Err3: "Emailの形式ではありません",
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "signup", "public_navbar", "public_navbarMobile")
			return

		}

		_, err = models.ExistEmail(r.PostFormValue("email"))

		password := r.PostFormValue("password")

		match := MatchPasswordString(password)

		if err == nil {
			errorMessage := ErrorMessage{
				Err3: "入力されたEmailは既に使用されています",
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "signup", "public_navbar", "public_navbarMobile")
			return
		} else if !match {
			errorMessage := ErrorMessage{
				Err3: "<br>※パスワードは半角英数字の<br>大文字１文字以上と数字を含めてください。<br>また、8文字以上15文字以下で入力してください。" +
					"<br> 半角英数字で入力してください"}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "signup", "public_navbar", "public_navbarMobile")
			return
		}

		user := models.User{
			Name:     name,
			Email:    r.PostFormValue("email"),
			PassWord: password,
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}

		user2, _ := models.GetUserByEmail(r.PostFormValue("email"))

		session, err := user2.CreateSession()
		if err != nil {
			log.Println(err)
		}

		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	_, err := sessionU(w, r)

	if err != nil {
		GenerateHTML(w, nil, "layout", "public_navbar", "login", "public_navbarMobile")

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	_, err = models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		errorMessage := ErrorMessage{
			Err5: "Emailが間違っています。"}

		pageData := map[string]interface{}{

			"error": errorMessage,
		}
		GenerateHTML(w, pageData, "layout", "public_navbar", "login", "public_navbarMobile")
		return

	}

	user, err := models.GetUserByEmailAndPassWord(r.PostFormValue("email"), models.Encrypt(r.PostFormValue("password")))
	//存在しない場合
	if err != nil {
		errorMessage := ErrorMessage{
			Err6: "Passwordが間違っています。"}

		pageData := map[string]interface{}{

			"error": errorMessage,
		}
		GenerateHTML(w, pageData, "layout", "public_navbar", "login", "public_navbarMobile")
		return
	}

	session, err := user.CreateSession()
	if err != nil {
		log.Println(err)
	}

	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.UUID,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

//ログアウト
func logout(writer http.ResponseWriter, request *http.Request) {
	_, err := sessionU(writer, request)

	if err != nil {
		GenerateHTML(writer, nil, "layout", "public_navbar", "login", "public_navbarMobile")
	} else {
		cookie, err := request.Cookie("_cookie")
		if err != nil {
			// log.Println(err)
			GenerateHTML(writer, nil, "layout", "public_navbar", "login", "public_navbarMobile")
		}

		if err != http.ErrNoCookie {
			session := models.Session{UUID: cookie.Value}
			session.DeleteSessionByUUID()
			cookie.MaxAge = -1             // 格納した変数cのMaxAgeフィールドに-1を指定
			http.SetCookie(writer, cookie) // 変更を反映するためにcをCookieにセット
		}

		GenerateHTML(writer, nil, "layout", "public_navbar", "login", "public_navbarMobile")
	}

}
