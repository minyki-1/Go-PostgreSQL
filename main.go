package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	middleware "psql/Tools/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	middleware.Setting(app)
	var (
		host     = os.Getenv("HOST")
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DBNAME")
	)
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbname)
	// db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {

		// * 테이블 생성
		// db.Migrator().CreateTable(&User{})

		// * 데이터 생성
		// user := User{Name: "CWIN77"}
		// result := db.Create(&user)
		// if result.Error != nil {
		// 	panic(result.Error)
		// }
		// db.Select("Name", "ID", "CreatedAt","UpdatedAt","DeletedAt").Create(&user)

		// * 데이터 가져오기
		// var result User
		// db.Model(User{}).First(&result)
		// fmt.Println(result)

		getData(db)
		return c.Status(200).JSON("*result")
	})

	app.Listen(":5000")
}

func createTable(db *sql.DB) {
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS test (
			id   SERIAL NOT NULL PRIMARY KEY,
			name TEXT
		)`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func dropTable(db *sql.DB) {
	statement, err := db.Prepare(`DROP TABLE users`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func getData(db *sql.DB) {
	rows, err := db.Query(`SELECT tablename FROM pg_catalog.pg_tables`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		// var name string
		// var id int

		var data interface{}
		err = rows.Scan(&data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s", data)
	}

	if err != nil {
		panic(err)
	}
}

func addData(db *sql.DB) {
	insertDynStmt := `insert into "users"("name") values($1)`
	// * insert into users ("name" ,created_at,updated_at) values ('cwub7777',NOW(),NOW())
	_, err := db.Exec(insertDynStmt, "CWIN")
	if err != nil {
		panic(err)
	}
}
