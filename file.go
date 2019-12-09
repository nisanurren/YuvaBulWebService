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

func main() {

	var err error
	AppDB, err := gorm.Open("mysql", "root:**************************************")
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
		AppDB.Where(map[string]interface{}{"mail": mail, "password": password}).Find(&user)
		if user.Password != "" {
			context.JSON(200, user)
			log.Print("kullanıcı bulundu...")
		} else {
			log.Println("kullanıcı bulunamadı")
			context.JSON(404, "kullanıcı yok")
		}

	})

	r.Run(":8080")

}
