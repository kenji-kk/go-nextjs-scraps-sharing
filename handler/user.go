package handler

import (
	"log"
	"net/http"

	"app/db"
	"app/jwt"

	"github.com/gin-gonic/gin"
)


func CheckErr(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}


func Signup(ctx *gin.Context) {
  user := new(db.User)
  if err := ctx.Bind(user); err != nil {
    ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
    return
  }

	AddedUser, err := user.AddUser()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

  ctx.JSON(http.StatusOK, gin.H{
    "msg": "Signed up successfully.",
    "jwt": jwt.GenerateJWT(&AddedUser),
		"user": AddedUser,
  })
}

func Signin(ctx *gin.Context) {
	user := new(db.User)
	if err := ctx.Bind(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	signinUser, err := user.Signin()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Signed in successfully.",
		"jwt": jwt.GenerateJWT(&signinUser),
		"user": signinUser,
	})
	
}

