package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string `json:"email" gorm:"unique; not null" validate:"email,required"`
	Password  string `json:"password" gorm:"not null" validate:"min=8,required"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	CreatedAt int    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int    `json:"updated_at" gorm:"autoUpdateTime"`
}

type Task struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null" validate:"required"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed" gorm:"not null"`
	UserID      string `json:"user_id"`
	User        User   `json:"user"`
}

func JSON() func(g *gin.Context) {
	return func(g *gin.Context) {
		g.Writer.Header().Set("Content-Type", "application/json")

		g.Next()
	}
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := binding.JSON.Bind(c.Request, obj); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return err
	}
	return nil
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func BindWith(obj interface{}, b binding.Binding, c *gin.Context) error {
	if err := b.Bind(c.Request, obj); err != nil {
		//c.AbortWithError(400, err).SetType(ErrorTypeBind)
		return err
	}
	return nil
}

func Validator() *validator.Validate {
	return validator.New()
}

func SignUp(g *gin.Context) {
	var data User

	validate := Validator()

	if err := g.BindJSON(&data); err != nil {
		g.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if err := validate.Struct(g.BindJSON(&data)); err != nil {
		g.AbortWithStatus(http.StatusBadRequest)

		return
	}

	g.JSON(http.StatusCreated, "POl")
}

func main() {
	g := gin.New()

	g.Use(JSON())
	g.Use(gin.Logger())
	g.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
	}))

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=postgres password=root port=5432 dbname=tasker",
	}))

	if err != nil {
		log.Fatal(err)
	}

	defer database.AutoMigrate(&User{}, &Task{})

	router := g.Group("/authentication")

	router.POST("/sign-up", SignUp)

	g.Run(":8000")
}
