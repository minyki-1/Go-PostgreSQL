package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	middleware "psql/Tools/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
		port     = 5432
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DBNAME")
	)
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

		return c.Status(200).JSON("Hello World!")
	})

	app.Listen(":5000")
}

func createTable(db *sql.DB) {
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS test (
			id        INT PRIMARY KEY,
			name      TEXT
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
	statement, err := db.Prepare(`DROP TABLE test`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func getData(db *sql.DB) {
	rows, err := db.Query(`SELECT "name", "id" FROM "test"`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var id int

		err = rows.Scan(&name, &id)
		if err != nil {
			panic(err)
		}

		fmt.Println(name, id)
	}

	if err != nil {
		panic(err)
	}
}
