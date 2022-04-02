package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"unique"`
	NameKanji string
	Password  string
	Message   string
}

func main() {
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("img", "./img")

	// indexページ
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	})

	// loginページ
	router.POST("/login", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()
		name := ctx.PostForm("name")
		password := ctx.PostForm("password")
		var user User
		err := db.Where("name = ? AND password = ?", name, password).Find(&user).Error
		if err != nil {
			ctx.String(http.StatusBadRequest, "名前もしくはパスワードが間違っています。%s", err)
			return
		}
		ctx.HTML(200, "message.html", gin.H{
			"name_kanji": user.NameKanji,
			"message":    user.Message,
		})
	})

	// 編集ページ
	router.GET("/edit", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()
		var users []User
		db.Order("created_at asc").Find(&users)

		ctx.HTML(200, "edit.html", gin.H{
			"result": users,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()

		name := ctx.PostForm("name")
		name_kanji := ctx.PostForm("name_kanji")
		password := ctx.PostForm("password")
		message := ctx.PostForm("message")
		db.Create(&User{Name: name, NameKanji: name_kanji, Password: password, Message: message})

		ctx.Redirect(302, "/edit")
	})

	router.GET("/delete/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("id is not a number")
		}
		var user User
		err = db.First(&user, id).Error
		if err != nil {
			ctx.String(http.StatusBadRequest, "取得失敗：%s", err)
		}
		err = db.Delete(&user).Error
		if err != nil {
			ctx.String(http.StatusBadRequest, "削除失敗：%s", err)
		}
		defer db.Close()

		ctx.Redirect(302, "/edit")
	})

	router.Run()
}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				panic(err)
			}
			db, err = gorm.Open(DBMS, CONNECT)
		}
	}

	return db
}
