package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../controller"
	"../config"
	"../middlewares"
	"github.com/gin-contrib/sessions"
)

func NewRouter() *gin.Engine {
	c := config.GetConfig();

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//Login-Flow cookie store
	store := sessions.NewCookieStore([]byte(c.GetString("server.cookie_store_secret")))
	router.Use(sessions.Sessions("login_session", store))

	//Controller
	loginController := new(controller.LoginController)
	userController := new(controller.UserController)

	//Public routes
	public := router.Group("/api/v1/")
	public.Handle(http.MethodGet, "/login/:provider/callback", loginController.AuthCallback)
	public.Handle(http.MethodGet, "/login/:provider", loginController.Auth)

	//Private routes
	private := router.Group("/api/v1")
	private.Use(middlewares.AuthorizationMiddleware)

	private.Handle(http.MethodGet, "/user", userController.GetAccount)
	private.Handle(http.MethodPut, "/user", userController.UpdateProfile)

	return router
}
