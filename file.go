package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	_ "net/http"
	_ "strconv"
)

var AppDB *gorm.DB

type User struct {
	UserID   int    `json:"user_id" `
	Mail     string `json:"mail" `
	Name     string `json:"name" `
	Surname  string `json:"surname" `
	Password string `json:"password" `
}

type Post struct {
	PostId            int    `json:"post_id" `
	Kind              string `json:"kind" `
	PostDescription   string `json:"post_description" `
	CreatorMail       string `json:"creator_mail" `
	City              string `json:"city" `
	Base64ImageString string `json:"base64_image_string" `
}

func main() {

	var err error
	AppDB, err := gorm.Open("mysql", "root:Nisanur77.@/yuvadb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Failed to connect database" + err.Error())
	}
	defer AppDB.Close()

	r := gin.Default()

	r.GET("/getAllUsers", func(context *gin.Context) {
		var users []User
		AppDB.Find(&users)
		context.JSON(200, users)

	})

	r.GET("/getUser/:mail", func(context *gin.Context) {
		var users User
		mail := context.Params.ByName("mail")

		AppDB.Where("mail = ?", mail).First(&users)
		context.JSON(200, users)
	})

	r.POST("/SignUpUser", func(context *gin.Context) {
		user := User{}
		user.Mail = context.PostForm("mail")
		user.Name = context.PostForm("name")
		user.Surname = context.PostForm("surname")
		user.Password = context.PostForm("password")

		AppDB.Create(user)
		context.JSON(200, user)

	})

	r.GET("SignInControl/:mail/:password", func(context *gin.Context) {
		mail := context.Params.ByName("mail")
		password := context.Params.ByName("password")
		user := User{}
		if mail == "" {
			context.JSON(400, "400 : Bad Request")
		} else if password == "" {
			context.JSON(400, "400 : Bad Request")
		}
		AppDB.Where(map[string]interface{}{"mail": mail, "password": password}).Find(&user)
		if user.Password != "" {
			context.JSON(200, user)
			log.Print("kullanıcı bulundu...")
		} else {
			log.Println("kullanıcı bulunamadı")
			context.JSON(404, "404 : Not Found")
		}

	})

	r.POST("/CreatePost", func(context *gin.Context) {
		post := Post{}

		post.CreatorMail = context.PostForm("creator_mail")
		post.Kind = context.PostForm("kind")
		post.City = context.PostForm("city")
		post.Base64ImageString = context.PostForm("base64_image_string")
		post.PostDescription = context.PostForm("post_description")

		if post.PostDescription == "" {
			context.JSON(400, "Bad Request.")
		} else {

			AppDB.Create(post)
			context.JSON(201, post)

		}

	})

	r.GET("/GetAllPosts", func(context *gin.Context) {
		var posts []Post
		AppDB.Find(&posts)
		context.JSON(200, posts)
	})

	r.Run(":8080")

}
