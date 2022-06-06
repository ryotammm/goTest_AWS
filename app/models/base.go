package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"goTest/app/config"
	"log"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
)

var Db *sql.DB

var err error

const (
	REGION = "ap-northeast-1"
)

func init() {
	// parameterNameをパラメータストアから取得
	parameterName := "/tabletennis/app/MYSQL_DATABASE"
	DB_NAME, err := fetchParameterStore(parameterName)
	if err != nil {
		os.Exit(0)
	}
	parameterName = "/tabletennis/app/USER_NAME"
	USER_NAME, err := fetchParameterStore(parameterName)
	if err != nil {
		os.Exit(0)
	}
	parameterName = "/tabletennis/app/MYSQL_PASSWORD"
	MYSQL_PASSWORD, err := fetchParameterStore(parameterName)
	if err != nil {
		os.Exit(0)
	}

	parameterName = "/tabletennis/app/MYSQL_HOST"
	MYSQL_HOST, err := fetchParameterStore(parameterName)
	if err != nil {
		os.Exit(0)
	}

	parameterName = "/tabletennis/app/MYSQL_PORT"
	MYSQL_PORT, err := fetchParameterStore(parameterName)
	if err != nil {
		os.Exit(0)
	}

	opt := "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	//データベース接続
	Db, err = sql.Open(config.Config.SQLDriver, USER_NAME+":"+MYSQL_PASSWORD+"@tcp("+MYSQL_HOST+":"+MYSQL_PORT+")/"+DB_NAME+opt)
	if err != nil {
		log.Fatalln(err)
	}
}

func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext

}

// パラメータストアから設定値取得
func fetchParameterStore(param string) (string, error) {

	sess := session.Must(session.NewSession())
	svc := ssm.New(
		sess,
		aws.NewConfig().WithRegion(REGION),
	)

	res, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(param),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "Fetch Error", err
	}

	value := *res.Parameter.Value
	return value, nil
}
