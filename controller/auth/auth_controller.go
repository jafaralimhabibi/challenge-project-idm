package controller

import (
	"context"
	"net/http"
	"time"

	config "challenge-project/config"
	lib "challenge-project/lib"
	userMdl "challenge-project/model/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Params!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := config.GetDatabase()
	count, err := db.Collection("users").CountDocuments(ctx, bson.M{"username": req.Username})
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Internal DB Error!")
	}
	if count > 0 {
		return lib.JSONError(c, http.StatusConflict, "Already exists!")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Hashed Error!")
	}

	user := userMdl.User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
		Password: string(hashed),
	}

	_, err = db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Insert Error!")
	}

	return lib.JSONCreated(c, map[string]string{
		"id":       user.ID.Hex(),
		"username": user.Username,
	})
}

func Login(c echo.Context) error {
	var req loginRequest
	var user userMdl.User

	if err := c.Bind(&req); err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Params!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := config.GetDatabase()

	err := db.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return lib.JSONError(c, http.StatusUnauthorized, "Not Found!")
	} else if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "db error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return lib.JSONError(c, http.StatusUnauthorized, "Invalid Credentials!")
	}

	cfg := config.GetConfig()

	claims := jwt.MapClaims{}
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Duration(cfg.JWT.Expire) * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tknStr, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Signing Token Error!")
	}

	return lib.JSONSuccess(c, map[string]interface{}{
		"token": tknStr,
		"user": map[string]interface{}{
			"id":       user.ID.Hex(),
			"username": user.Username,
		},
	})
}
