package rest

import (
	"be/auth"
	"be/service"
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Rest struct {
	svc  service.Service
	auth auth.Auth
}

func New(_svc service.Service, _a auth.Auth) *Rest {
	return &Rest{svc: _svc, auth: _a}
}

func (r *Rest) Register(re *gin.Engine) {
	re.GET("/ping", r.HandlePing)
	re.GET("/user/:username", r.GetUser)
	re.GET("/users", r.GetUsers)
	re.POST("/user", r.PostUser)
	re.POST("/validateuser", r.CheckUsernamePassword)
	re.GET("/auth/google/login", r.GoogleLogin)
	re.GET("/auth/google/callback", r.GoogleCallback)
}

func (r *Rest) HandlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (r *Rest) GoogleLogin(c *gin.Context) {
	url := r.auth.OauthGoogleLogin(c.Writer, c.Request)
	c.JSON(http.StatusOK, url)
	return
}

func (r *Rest) GoogleCallback(c *gin.Context) {
	token, err := r.auth.OauthGoogleCallback(c.Writer, c.Request)
	if err != nil {
		c.Error(err)
	}

	// set access token to cookie
	var expiration = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	cookie := http.Cookie{Name: "at", Value: token, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, "OK")
	return
}

func (r *Rest) GetUser(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, errors.New("Bad Request"))
	}

	u, err := r.svc.GetUser(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, u)
}

func (r *Rest) GetUsers(c *gin.Context) {
	users, err := r.svc.GetUsers(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, users)
}

func (r *Rest) PostUser(c *gin.Context) {
	var req service.UserRequest
	c.BindJSON(&req)
	err := r.svc.InsertUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "OK")
}

func (r *Rest) CheckUsernamePassword(c *gin.Context) {
	var req service.UserPasswordCheckRequest
	c.BindJSON(&req)
	res, err := r.svc.CheckUsernamePassword(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, res)
}
