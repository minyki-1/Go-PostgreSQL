package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_component "psql/Tools/Router/component"
	_component_like "psql/Tools/Router/component/like"
	_project "psql/Tools/Router/project"
	_project_component "psql/Tools/Router/project/component"
	_style "psql/Tools/Router/style"
	_team "psql/Tools/Router/team"
	_team_member "psql/Tools/Router/team/member"
	_test "psql/Tools/Router/test"
	_user "psql/Tools/Router/user"
	middleware "psql/Tools/middleware"
	"psql/Tools/mongodb"
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

	if err := mongodb.ConnectDB(os.Getenv("MONGODB_URI")); err != nil {
		log.Fatal("Error mongoDB connect")
	}

	app.Get("/component", _component.Get)
	app.Post("/component", _component.Post)
	app.Put("/component", _component.Put)
	app.Delete("/component", _component.Delete)
	app.Put("/component/like", _component_like.Put)

	app.Post("/project", _project.Post)
	app.Put("/project", _project.Put)
	app.Delete("/project", _project.Delete)
	app.Get("/project/:params", _project.Get)
	app.Put("/project/component", _project_component.Put)
	app.Delete("/project/component", _project_component.Delete)

	app.Get("/team/:params", _team.Get)
	app.Post("/team", _team.Post)
	app.Delete("/team", _team.Delete)
	app.Put("/team", _team.Put)
	app.Put("/team/member", _team_member.Put)

	app.Get("/user/:params", _user.Get)
	app.Post("/user", _user.Post)

	app.Get("/style/:params", _style.Get)

	if os.Getenv("ENV_MODE") == "development" {
		app.Get("/test", _test.Test)
	}

	app.Listen(":5000")
}

func ConnectDB() (*sql.DB, error) {
	var (
		host     = os.Getenv("HOST")
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DBNAME")
	)
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
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
	rows, err := db.Query(`SELECT * FROM users`)
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
