package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/bitphinix/barbra_backend/config"
	"github.com/bitphinix/barbra_backend/controller"
	"github.com/bitphinix/barbra_backend/middlewares"
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
	suggestionController := new(controller.SuggestionController)

	//Public routes
	public := router.Group("/api/v1/")
	public.Handle(http.MethodGet, "/login/:provider/callback", loginController.AuthCallback)
	public.Handle(http.MethodGet, "/login/:provider", loginController.Auth)

	//Private routes
	private := router.Group("/api/v1")
	private.Use(middlewares.AuthorizationMiddleware(false))

	private.Handle(http.MethodGet, "/user", userController.GetAccount)
	private.Handle(http.MethodPut, "/user", userController.UpdateProfile)

	//Private routes (enrolled accounts only)
	enrolled := router.Group("/api/v1")
	enrolled.Use(middlewares.AuthorizationMiddleware(true))

	enrolled.Handle(http.MethodGet, "/suggestions", suggestionController.GetSuggestions)

	return router
}
