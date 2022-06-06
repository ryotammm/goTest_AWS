package models

import (
	"github.com/google/uuid"
)

type Tags struct {
	ID  uuid.UUID
	Tag string
}

func GetTags(tagWord []string) (tagsP []Tags, err error) {

	for _, v := range tagWord {

		stmt, err := Db.Prepare("select id, tag from tags where tag = ?")
		checkErr(err)
		row, err := stmt.Query(v)

		checkErr(err)
		defer row.Close()

		for row.Next() {
			var tag Tags
			err = row.Scan(&tag.ID, &tag.Tag)
			checkErr(err)
			tagsP = append(tagsP, tag)
		}
		row.Close()
	}
	return tagsP, err
}

func CreateTags(uid uuid.UUID, tags []string) (err error) {
	stmt, err := Db.Prepare("insert tags SET id=?,tag=?")
	checkErr(err)

	for _, tag := range tags {

		_, err = stmt.Exec(uid, tag)
		checkErr(err)

	}
	return err
}
