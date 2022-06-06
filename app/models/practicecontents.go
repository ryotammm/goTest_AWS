package models

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Practicecontent struct {
	ID         int
	UserID     int
	Prefecture string
	Place      string
	Strat_time string
	End_time   string
	Scale      string
	Tags       string
	UUID       string
	Describe   string
	CreatedAt  time.Time
}

func (p *Practicecontent) CreatePracticecontent(userID int) (Practicecontent_uuid uuid.UUID, err error) {
	Practicecontent_uuid = createUUID()

	//データの挿入

	stmt, err := Db.Prepare("insert practicecontents SET user_id=?, prefecture=?,place=?,strat_time=?,end_time=?,scale=?,tags=?,describes=?,uuid=?,created_at=?")
	checkErr(err)
	_, err = stmt.Exec(userID, p.Prefecture, p.Place, p.Strat_time, p.End_time, p.Scale, p.Tags, p.Describe, Practicecontent_uuid, time.Now())
	checkErr(err)

	return Practicecontent_uuid, err
}

func GetPracticecontent() (pracs []Practicecontent, err error) {
	stmt, err := Db.Prepare("select id, user_id, prefecture, place,strat_time,end_time,scale,tags,describes,uuid,created_at FROM practicecontents")
	checkErr(err)
	rows, err := stmt.Query()
	checkErr(err)

	for rows.Next() {
		var prac Practicecontent
		err = rows.Scan(
			&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)

		checkErr(err)

		pracs = append(pracs, prac)

	}
	rows.Close()
	return pracs, err
}

func GetPracticecontentByUUID(uid uuid.UUID) (prac Practicecontent) {
	prac = Practicecontent{}
	stmt, err := Db.Prepare("select id, user_id, prefecture, place,strat_time,end_time,scale,tags,describes,uuid,created_at FROM practicecontents where  uuid = ?")
	checkErr(err)

	row, err := stmt.Query(uid)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&prac.ID, &prac.UserID, &prac.Prefecture, &prac.Place, &prac.Strat_time, &prac.End_time, &prac.Scale, &prac.Tags, &prac.UUID, &prac.Describe, &prac.CreatedAt)
	checkErr(err)
	return prac
}

//idで一つの記事を取得する
func GetPracticecontentByID(id int) (prac Practicecontent) {
	prac = Practicecontent{}
	stmt, err := Db.Prepare("select id, user_id, prefecture, place,strat_time,end_time,scale,tags,describes,uuid,created_at FROM practicecontents where  id = ?")
	checkErr(err)

	row, err := stmt.Query(id)
	checkErr(err)
	defer row.Close()
	row.Next()
	err = row.Scan(&prac.ID, &prac.UserID, &prac.Prefecture, &prac.Place, &prac.Strat_time, &prac.End_time, &prac.Scale, &prac.Tags, &prac.Describe, &prac.UUID, &prac.CreatedAt)
	checkErr(err)
	return prac
}

/**/
func SearchPrefectures(pref string) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale, tags,describes,uuid, created_at FROM practicecontents where prefecture = ?")
	checkErr(err)
	row, err := stmt.Query(pref)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err
}

func SearchPrefecturesAndTags(pref string, tags []string) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale,describes, uuid, created_at FROM practicecontents where place = ?")
	checkErr(err)
	row, err := stmt.Query(pref)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err
}

//都道府県のみ
func SearchPrefecturesX(pref string) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale, tags,describes, uuid, created_at FROM practicecontents where prefecture = ? order by created_at DESC; ")
	checkErr(err)
	row, err := stmt.Query(pref)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err
}

//検索フォームのみ
func SearchTagsX(tags string) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale,tags,describes, uuid, created_at FROM practicecontents where tags like ? or prefecture  like ? or   place  like ?  order by created_at DESC;")
	checkErr(err)

	like := "%"

	tagsS1 := fmt.Sprintf("%s%s%s", like, tags, like)
	tagsS2 := fmt.Sprintf("%s%s%s", like, tags, like)
	tagsS3 := fmt.Sprintf("%s%s%s", like, tags, like)

	row, err := stmt.Query(tagsS1, tagsS2, tagsS3)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err

}

//検索フォーム && 都道府県
func SearchPrefecturesAndTagsX(pref string, tags string) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale,tags,describes, uuid,created_at FROM practicecontents where prefecture = ? and  (tags like ?  or  place  like ? ) order by created_at DESC;")
	checkErr(err)

	like := "%"

	tagsS1 := fmt.Sprintf("%s%s%s", like, tags, like)
	tagsS2 := fmt.Sprintf("%s%s%s", like, tags, like)

	row, err := stmt.Query(pref, tagsS1, tagsS2)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err

}

//Start 画面
func StartPrefectures() (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id,prefecture, place, strat_time,end_time,scale,tags,describes,uuid,created_at FROM practicecontents order by  created_at DESC LIMIT 50; ")
	checkErr(err)

	row, err := stmt.Query()
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)
	}
	return pracs, err
}

//ユーザの投稿取得 複数
func GetPracticecontentByUserID(id int) (pracs []Practicecontent, err error) {

	stmt, err := Db.Prepare("select id, user_id, prefecture, place,strat_time,end_time,scale,tags,describes,uuid,created_at FROM practicecontents where user_id = ? order by created_at DESC ")
	checkErr(err)

	row, err := stmt.Query(id)
	checkErr(err)

	for row.Next() {
		var prac Practicecontent
		err = row.Scan(
			&prac.ID,
			&prac.UserID,
			&prac.Prefecture,
			&prac.Place,
			&prac.Strat_time,
			&prac.End_time,
			&prac.Scale,
			&prac.Tags,
			&prac.Describe,
			&prac.UUID,
			&prac.CreatedAt)

		if err != nil {
			log.Fatalln(err)
		}
		pracs = append(pracs, prac)

	}
	row.Close()
	return pracs, err

}

func Deleterecruitment(id int) error {
	stmt, err := Db.Prepare("delete from practicecontents where id = ?")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)
	return err
}

func DeleterecruitmentByUUID(uuid string) error {
	stmt, err := Db.Prepare("delete from practicecontents where uuid = ?")
	checkErr(err)
	_, err = stmt.Exec(uuid)

	checkErr(err)
	return err
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error")
	}
}
