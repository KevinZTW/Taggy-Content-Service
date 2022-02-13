package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/option"
)

type RepositoryClient struct {
	FirestoreClient *firestore.Client
	MySQLDB         *sql.DB
}

func NewRepositoryClient() *RepositoryClient {
	c := &RepositoryClient{}
	c.FirestoreClient = NewFirestoreClient()
	c.MySQLDB = NewMySQLDB()
	return c
}

func NewFirestoreClient() *firestore.Client {
	fmt.Println(config.FireStore.CredintialPath)
	ctx := context.Background()

	sa := option.WithCredentialsFile(config.FireStore.CredintialPath)
	// sa := option.WithCredentialsFile("./knowledge-base-tw-c00d0f55b34a.json")
	// sa := option.WithCredentialsFile("../repository/knowledge-base-tw-c00d0f55b34a.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client

}

func NewMySQLDB() *sql.DB {
	con := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Db)
	fmt.Println(con)
	db, err := sql.Open("mysql", con)

	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
