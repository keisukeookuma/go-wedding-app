package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type User struct {
	gorm.Model
	Name          string `gorm:"unique"`
	OrganizerName string
	NameKanji     string
	Password      string
	Message       string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("env読み込み失敗")
	}
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	// indexページ
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{
			"error": "",
		})
	})

	// loginページ
	router.POST("/message", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()
		name := ctx.PostForm("name")
		password := ctx.PostForm("password")
		var user User
		err := db.Where("name = ? AND password = ?", name, password).Find(&user).Error
		if err != nil {
			ctx.HTML(200, "login.html", gin.H{
				"error": "false",
			})
			return
		}
		ctx.HTML(200, "message.html", gin.H{
			"name_kanji":     user.NameKanji,
			"message":        user.Message,
			"organizer_name": user.OrganizerName,
		})
	})

	// 編集ページ
	router.GET("/edit_3160ktsia", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()
		var users []User
		db.Order("created_at asc").Find(&users)

		ctx.HTML(200, "edit.html", gin.H{
			"result": users,
		})
	})

	// 編集ページ
	router.POST("/update", func(ctx *gin.Context) {
		fmt.Println("hello")
		db := sqlConnect()
		defer db.Close()

		id := ctx.PostForm("id")
		name := ctx.PostForm("name")
		name_kanji := ctx.PostForm("name_kanji")
		password := ctx.PostForm("password")
		message := ctx.PostForm("message")

		err := db.Model(&User{}).Where("id = ?", id).Update(&User{Name: name, NameKanji: name_kanji, Password: password, Message: message}).Error
		if err != nil {
			ctx.String(
				200, "ユーザー編集失敗：%v", err,
			)
		}
		ctx.Redirect(302, "/edit_3160ktsia")
	})

	router.POST("/new", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()

		name := ctx.PostForm("name")
		name_kanji := ctx.PostForm("name_kanji")
		password := ctx.PostForm("password")
		message := ctx.PostForm("message")
		organizer_name := ctx.PostForm("organizer_name")
		err := db.Create(&User{Name: name, NameKanji: name_kanji, Password: password, Message: message, OrganizerName: organizer_name}).Error
		if err != nil {
			ctx.String(
				200, "ユーザー作成失敗：%v", err,
			)
		}
		ctx.Redirect(302, "/edit_3160ktsia")
	})

	router.POST("/delete_kassi", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.PostForm("id")
		p := ctx.PostForm("pass")
		if p != "kkdai" {
			ctx.String(
				200, "ユーザー削除失敗：%v", err,
			)
		}
		id, err := strconv.Atoi(n)
		if err != nil {
			ctx.String(
				200, "ユーザー削除失敗：%v", err,
			)
		}
		var user User
		err = db.First(&user, id).Error
		fmt.Println(user)
		if err != nil {
			log.Print(http.StatusBadRequest, "取得失敗：%s", err)
		}
		err = db.Unscoped().Delete(&user).Error
		if err != nil {
			log.Print(http.StatusBadRequest, "削除失敗：%s", err)
		}
		defer db.Close()

		ctx.Redirect(302, "/edit_3160ktsia")
	})

	router.Run()
}

func sqlConnect() (database *gorm.DB) {
	DBMS := os.Getenv("DB_MS")
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	PROTOCOL := "tcp(db:3306)"
	DBNAME := os.Getenv("DB_NAME")

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
