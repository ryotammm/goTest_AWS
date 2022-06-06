package controllers

import (
	"fmt"
	"goTest/app/models"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ErrorMessage struct {
	Err1 string
	Err2 string
	Err3 string
	Err4 string
	Err5 string
	Err6 string
	Err7 string
}
type SuccessMessage struct {
	Succ1 string
	Succ2 string

	//画像
}

func nameEdit(w http.ResponseWriter, r *http.Request) {

	sess, err := sessionU(w, r)
	if err != nil {
		log.Fatalln(err)
	}
	switch r.Method {
	case http.MethodGet:

		pageData := pageData(sess.Name, nil)

		GenerateHTML(w, pageData, "layout", "name_edit", "private_navbar", "private_navbarMobile")

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

			pageData := pageData(sess.Name, nil)
			pageData["error"] = errorMessage

			GenerateHTML(w, pageData, "layout", "name_edit", "private_navbar", "private_navbarMobile")
			return
		}

		result := models.NameEdit(sess.UserID, name)

		if result != nil {
			fmt.Println("変更できませんでした")
		}

		err3 := models.NameEditSession(sess.UserID, name)
		if err3 != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	}

}

func emailEdit(w http.ResponseWriter, r *http.Request) {

	sess, err := sessionU(w, r)
	if err != nil {
		log.Fatalln(err)
	}

	switch r.Method {
	case http.MethodGet:
		profile := Profile{
			Name:  sess.Name,
			Email: sess.Email,
		}

		profileData := map[string]interface{}{
			"profile": profile,
			"pracs":   nil,
		}

		if err != nil {
			log.Fatalln(err)
		}

		GenerateHTML(w, profileData, "layout", "email_edit", "private_navbar", "private_navbarMobile")

	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		email := r.PostFormValue("email")

		result2 := strings.Compare(sess.Email, email)
		if result2 == 0 {

			errorMessage := ErrorMessage{
				Err4: "現在のEmailと同じです",
			}

			profile := Profile{
				Name:  sess.Name,
				Email: sess.Email,
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "email_edit", "private_navbar", "private_navbarMobile")
			return

		}
		result3 := MatchEmailString(email)

		if !result3 {

			profile := Profile{
				Name:  sess.Name,
				Email: sess.Email,
			}

			errorMessage := ErrorMessage{
				Err7: "Emailの形式ではありません",
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}
			GenerateHTML(w, profileData, "layout", "email_edit", "private_navbar", "private_navbarMobile")
			return
		}

		_, err = models.ExistEmail(email)

		if err == nil {
			errorMessage := ErrorMessage{
				Err3: "入力されたEmailは既に使用されています",
			}

			profile := Profile{
				Name:  sess.Name,
				Email: sess.Email,
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "email_edit", "private_navbar", "private_navbarMobile")
			return
		}

		result := models.EmailEdit(sess.UserID, email)
		if result != nil {
			fmt.Println("変更できませんでした")
		}

		sess.EmailEditSession(sess.ID, email)

		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func passwordEdit(w http.ResponseWriter, r *http.Request) {

	sess, err := sessionU(w, r)
	if err != nil {
		log.Fatalln(err)
	}

	switch r.Method {
	case http.MethodGet:
		profile := Profile{
			Name:  sess.Name,
			Email: sess.Email,
		}

		profileData := map[string]interface{}{
			"profile": profile,
			"pracs":   nil,
		}

		GenerateHTML(w, profileData, "layout", "password_edit", "private_navbar", "private_navbarMobile")

	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}

		Current_Password := r.PostFormValue("Current_Password")

		crypto_pass := models.Encrypt(Current_Password)

		_, err = models.CheckPassWord(sess.UserID, crypto_pass)

		New_Password := r.PostFormValue("New_Password")
		New_Password_Confirm := r.PostFormValue("New_Password_Confirm")

		match := MatchPasswordString(New_Password)

		//入力比較
		result := strings.Compare(New_Password, New_Password_Confirm)

		profile := Profile{
			Name:  sess.Name,
			Email: sess.Email,
		}

		if err != nil {

			errorMessage := ErrorMessage{
				Err1: "現在のパスワードが違います",
			}

			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}

			GenerateHTML(w, profileData, "layout", "password_edit", "private_navbar", "private_navbarMobile")

			log.Println("現在のパスワードが違います")
		} else if !match {
			errorMessage := ErrorMessage{
				Err2: "<br>※パスワードは半角英数字の<br>大文字１文字以上と数字を含めてください。<br>また、8文字以上15文字以下で入力してください。" +
					"<br> 半角英数字で入力してください",
			}
			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}
			GenerateHTML(w, profileData, "layout", "password_edit", "private_navbar", "private_navbarMobile")
			return
		} else if result != 0 {
			errorMessage := ErrorMessage{
				Err2: "新しいパスワードと新しいパスワード確認が一致しません",
			}
			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   nil,
				"error":   errorMessage,
			}
			GenerateHTML(w, profileData, "layout", "password_edit", "private_navbar", "private_navbarMobile")

			log.Println("新しいパスワードと新しいパスワード確認が一致しません")
		} else {

			err = models.PasswordEdit(sess.UserID, models.Encrypt(New_Password))
			if err != nil {
				log.Fatalln(err)
			}

			SuccessMessage := SuccessMessage{
				Succ1: "パスワードを変更しました",
			}
			pracs, err := models.GetPracticecontentByUserID(sess.UserID)
			if err != nil {
				log.Fatalln(err)
			}
			profileData := map[string]interface{}{
				"profile": profile,
				"pracs":   pracs,
				"succ":    SuccessMessage,
			}

			GenerateHTML(w, profileData, "layout", "profile", "private_navbar", "private_navbarMobile")

		}

	}
}

//パスワードのチェック
func MatchPasswordString(password string) bool {

	if len(password) < 8 || len(password) > 16 {
		return false
	}
	// must contain at least each one uppercase, lowercase and number
	reg := []*regexp.Regexp{
		regexp.MustCompile(`[a-zA-Z\d]`), regexp.MustCompile(`[A-Z]`),
		regexp.MustCompile(`\d`),
	}
	for _, r := range reg {
		if r.FindString(password) == "" {
			return false
		}
	}
	return true
}

//Emailドのチェック
func MatchEmailString(email string) bool {

	reg := regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)

	var result = true

	if reg.FindString(email) == "" {
		result = false
		return result
	}
	return result
}
