package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		Id     uint64 `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Delete bool   `json:"-"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	PORT = ":8080"
)

var (
	DataStore = make(map[uint64]User)
)

func main() {

	DataStore = map[uint64]User{
		1: {
			Id:     1,
			Name:   "John Doe",
			Email:  "jhondoe@gmail.com",
			Delete: false,
		},
		2: {
			Id:     2,
			Name:   "Jane Doe",
			Email:  "janedow@gmail.com",
			Delete: false,
		},
	}

	//init server
	ginServer := gin.Default()
	//init middleware
	ginServer.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	//init router
	//get all users
	ginServer.GET("/users", func(ctx *gin.Context) {
		//api ini akan menampilkan semua users dari data store

		//business logic
		var users []User
		for _, usr := range DataStore {
			if !usr.Delete {
				users = append(users, usr)
			}
		}

		//response
		ctx.JSON(http.StatusOK, Response{
			Message: "Success get users",
			Data:    users,
		})
	})

	//get user by id
	ginServer.GET("/user", func(ctx *gin.Context) {
		id := ctx.Query("id")

		if id == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid query params",
			})
			return
		}

		//transform id string to int
		idUint, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid id params",
			})
			return
		}

		//business logic
		//get from data store
		user, ok := DataStore[uint64(idUint)]

		if !ok || user.Delete {
			ctx.JSON(http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		//response
		ctx.JSON(http.StatusOK, Response{
			Message: "Success get user",
			Data:    user,
		})
	})

	//create user
	ginServer.POST("/user", func(ctx *gin.Context) {
		//get body
		var usrIn User
		if err := ctx.Bind(&usrIn); err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid payload",
			})
			return
		}
		//validate
		if usrIn.Email == "" || usrIn.Name == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid payload",
			})
			return
		}
		//business logic
		//email is unique
		for _, usr := range DataStore {
			if strings.EqualFold(usr.Email, usrIn.Email) {
				ctx.JSON(http.StatusBadRequest, Response{
					Message: "Email already registered ",
				})
				return
			}
		}

		// success
		idInsert := len(DataStore) + 1
		usrIn.Id = uint64(idInsert)
		DataStore[uint64(idInsert)] = usrIn
		ctx.JSON(http.StatusOK, Response{
			Message: "Success create user",
			Data:    usrIn,
		})
	})

	//update user
	ginServer.PUT("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid query params",
			})
			return
		}

		//transform id string to int
		idUint, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid id params",
			})
			return
		}

		//binding payload
		var usrIn User
		if err := ctx.Bind(&usrIn); err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid payload",
			})
			return
		}
		//b logic
		//get user data if exist
		user, ok := DataStore[uint64(idUint)]
		if !ok || user.Delete {
			ctx.JSON(http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		//validate name
		if usrIn.Name == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid payload",
			})
			return
		}

		user.Name = usrIn.Name

		//update data store
		DataStore[uint64(idUint)] = user

		//response
		ctx.JSON(http.StatusAccepted, Response{
			Message: "Success update user",
			// Data:    user,
		})

	})

	//delete user
	ginServer.DELETE("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid query params",
			})
			return
		}

		//transform id string to int
		idUint, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Invalid id params",
			})
			return
		}

		//get user data if exist
		user, ok := DataStore[uint64(idUint)]
		if !ok || user.Delete {
			ctx.JSON(http.StatusNotFound, Response{
				Message: "User not found",
			})
			return
		}

		//delete user
		user.Delete = true
		DataStore[uint64(idUint)] = user

		//response
		ctx.JSON(http.StatusAccepted, Response{
			Message: "Success delete user",
		})
	})

	ginServer.Run(PORT)

}
