package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func main() {
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(ctx *gin.Context) {
		db := sqlConnect()
		var users []User
		db.Order("created_at asc").Find(&users)
		defer db.Close()

		ctx.HTML(200, "index.html", gin.H{
			"users": users,
		})
	})

	router.POST("/create", func(ctx *gin.Context) {
		db := sqlConnect()
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")
		fmt.Println("create user " + name + " with email " + email)
		db.Create(&User{Name: name, Email: email})
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	router.POST("/update/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")

		id, err := strconv.Atoi(n)
		if err != nil {
			panic("id is not a number")
		}

		var user User
		db.First(&user, id)
		user.Name = name
		user.Email = email
		db.Save(&user)
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("id is not a number")
		}
		var user User
		db.First(&user, id)
		db.Delete(&user)
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	router.Run()
}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"

	// CONNECT => USER:PASS@tcp(db:3306)/DBNAME?OPTIONS
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		fmt.Println("DB接続失敗")
		panic(err)
	}

	fmt.Println("DB接続成功")

	return db
}
