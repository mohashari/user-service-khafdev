package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

type Users struct {
	gorm.Model
	SecureId  string `json:"secure_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type Result struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"emai"`
	Password string `json:"password"`
}

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=ec2-54-75-230-41.eu-west-1.compute.amazonaws.com port=5432 user=goknnyynzexcbm password=4891288138b8ea55471458366e14ed8e7442d740b4d554b8c40afac546bf5e37 dbname=d306bneg8bcb87 sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Db Connected....")
	}
	db.AutoMigrate(&Users{})
}

func generateUUID() string {
	u1, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	return u1.String()
}

func createUser(c echo.Context) error {
	u := new(Users)
	if err := c.Bind(u); err != nil {
		return err
	}

	db.Create(&Users{
		SecureId:  generateUUID(),
		Name:      u.Name,
		Email:     u.Email,
		NoTelepon: u.NoTelepon,
		Password:  u.Password,
		Role:      u.Role})
	return c.JSON(http.StatusOK, true)
}

func getUser(c echo.Context) error {

	var model []Users
	db.Find(&model)
	if len(model) <= 0 {
		c.JSON(http.StatusNotFound, "Data not fount")
	}
	return c.JSON(http.StatusOK, model)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func registerMember(c echo.Context) error {
	u := new(Users)

	if err := c.Bind(u); err != nil {
		return err
	}

	emailEmpty := Result{StatusCode: http.StatusBadRequest, Data: "Email tidak boleh kosong"}
	if u.Email == "" {
		return c.JSON(http.StatusBadRequest, emailEmpty)
	}

	nameEmpty := Result{StatusCode: http.StatusBadRequest, Data: "Nama tidak boleh kosong"}
	if u.Name == "" {
		return c.JSON(http.StatusBadRequest, nameEmpty)
	}

	noTeleponEmpty := Result{StatusCode: http.StatusBadRequest, Data: "No telepon tidak boleh kosong"}
	if u.NoTelepon == "" {
		return c.JSON(http.StatusBadRequest, noTeleponEmpty)
	}

	passwordEmpty := Result{StatusCode: http.StatusBadRequest, Data: "Password tidak boleh kosong"}
	if u.Password == "" {
		return c.JSON(http.StatusBadRequest, passwordEmpty)
	}

	password, _ := HashPassword(u.Password)
	db.Create(&Users{
		SecureId:  generateUUID(),
		Name:      u.Name,
		Email:     u.Email,
		NoTelepon: u.NoTelepon,
		Password:  password,
		Role:      "MEMBER"})

	result := Result{StatusCode: http.StatusOK, Data: true}
	return c.JSON(http.StatusOK, result)
}

func login(c echo.Context) error {
	var user Users
	u := new(LoginRequest)
	if err := c.Bind(u); err != nil {
		return err
	}
	db.Find(&user, u.Email)
	if CheckPasswordHash(u.Password, user.Password) {
		result := Result{StatusCode: http.StatusOK, Data: "ok"}
		return c.JSON(http.StatusOK, result)
	} else {
		result := Result{StatusCode: http.StatusBadRequest, Data: "Password not Corect"}
		return c.JSON(http.StatusBadRequest, result)
	}

}

func main() {

	e := echo.New()

	e.POST("/user", createUser)
	e.POST("/register/member", registerMember)
	e.GET("/user", getUser)
	e.POST("/login", login)

	e.Logger.Info(e.Start(":"))

}
