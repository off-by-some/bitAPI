package user

import (
	json "encoding/json"
	"io"
	"io/ioutil"
	http "net/http"

	db "github.com/Pholey/bitAPI/db"
	pass "github.com/Pholey/bitAPI/resources/lib"

	"github.com/labstack/echo"
)

// User - Struct for dealing with users
type User struct {
	Salt       []byte `db:"password_salt"`
	Iterations int    `db:"password_iterations"`
	Hash       []byte `db:"password_hash"`
	Password   string `json:"password"`
	UserName   string `db:"username" json:"userName"`
	Name       string `db:"name"     json:"name"`
	Email      string `db:"email"    json:"email"`
}

// Users - Our array of users
type Users []User

// Create - Create a user
func Create(c echo.Context) error {
	// TODO(pholey): Abstract this out or find a better lib
	body, err := ioutil.ReadAll(io.LimitReader(c.Request().Body(), 1048576))

	if err != nil {
		// TODO(pholey): Proper error handing
		panic(err)
	}

	// Grab our user data
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		panic(err)
	}

	// TODO(pholey): Password validation, refactoring
	// Hash our password
	salt, iterations, hash := pass.HashPass(user.Password)

	// Create our statement
	sql, args, _ := db.Sq.
		Insert("\"user\"").
		Columns("username", "name", "email", "password_hash",
			"password_iterations", "password_salt").
		Values(user.UserName, user.Name, user.Email, hash, iterations, salt).
		ToSql()

	rows, err := db.Client.Query(sql, args...)

	_ = (rows)

	if err != nil {
		panic(err)
	}

	c.Response().Header().Set("Content-Type", "application/json;charset=UTF-8")
	// Since we were sent JSON, we should be nice and return an empty JSON object
	c.JSON(http.StatusCreated, struct{}{})

	return nil;
}
