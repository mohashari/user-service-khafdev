package main

import (
	"fmt"
	"khaf-dev/helper"
	"net/http"

	model "khaf-dev/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var db *gorm.DB

type users model.Users

type resultModel model.Result

type loginRequest model.LoginRequest

func generateUUID() string {
	u1, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	return u1.String()
}

func createUser(c echo.Context) error {
	u := new(users)
	if err := c.Bind(u); err != nil {
		return err
	}

	db.Create(&users{
		SecureId:  generateUUID(),
		Name:      u.Name,
		Email:     u.Email,
		NoTelepon: u.NoTelepon,
		Password:  u.Password,
		Role:      u.Role})
	return c.JSON(http.StatusOK, true)
}

func getUser(c echo.Context) error {

	var model []users
	db.Find(&model)
	if len(model) <= 0 {
		c.JSON(http.StatusNotFound, "Data not fount")
	}
	return c.JSON(http.StatusOK, model)
}

func registerMember(c echo.Context) error {
	u := new(users)

	if err := c.Bind(u); err != nil {
		return err
	}

	emailEmpty := resultModel{StatusCode: http.StatusBadRequest, Data: "Email tidak boleh kosong"}
	if u.Email == "" {
		return c.JSON(http.StatusBadRequest, emailEmpty)
	}

	nameEmpty := resultModel{StatusCode: http.StatusBadRequest, Data: "Nama tidak boleh kosong"}
	if u.Name == "" {
		return c.JSON(http.StatusBadRequest, nameEmpty)
	}

	noTeleponEmpty := resultModel{StatusCode: http.StatusBadRequest, Data: "No telepon tidak boleh kosong"}
	if u.NoTelepon == "" {
		return c.JSON(http.StatusBadRequest, noTeleponEmpty)
	}

	passwordEmpty := resultModel{StatusCode: http.StatusBadRequest, Data: "Password tidak boleh kosong"}
	if u.Password == "" {
		return c.JSON(http.StatusBadRequest, passwordEmpty)
	}

	password, _ := helper.HashPassword(u.Password)
	db.Create(&users{
		SecureId:  generateUUID(),
		Name:      u.Name,
		Email:     u.Email,
		NoTelepon: u.NoTelepon,
		Password:  password,
		Role:      "MEMBER"})

	result := resultModel{StatusCode: http.StatusOK, Data: true}
	return c.JSON(http.StatusOK, result)
}

func login(c echo.Context) error {
	var user users
	u := new(loginRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	db.Find(&user, u.Email)
	if helper.CheckPasswordHash(u.Password, user.Password) {
		result := resultModel{StatusCode: http.StatusOK, Data: "ok"}
		return c.JSON(http.StatusOK, result)
	}
	result := resultModel{StatusCode: http.StatusBadRequest, Data: "Password not Corect"}
	return c.JSON(http.StatusBadRequest, result)

}

func main() {

	e := echo.New()

	e.POST("/user", createUser)
	e.POST("/register/member", registerMember)
	e.GET("/user", getUser)
	e.POST("/login", login)

	e.Logger.Info(e.Start(":8089"))

}
