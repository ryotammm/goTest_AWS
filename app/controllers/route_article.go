package controllers

import (
	"goTest/app/models"
	"log"
	"net/http"
	"strconv"
)

func detail(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	id := r.FormValue("id")
	idInt, _ := strconv.Atoi(id)

	prac := models.GetPracticecontentByID(idInt)

	var pracS []models.Practicecontent

	pracS = append(pracS, prac)
	sess, err := sessionU(w, r)

	pageData := pageData(sess.Name, pracS)
	if err != nil {
		GenerateHTML(w, pageData, "layout", "public_navbarMobile", "public_navbar", "detail")
	} else {
		GenerateHTML(w, pageData, "layout", "private_navbarMobile", "private_navbar", "detail")
	}

}
