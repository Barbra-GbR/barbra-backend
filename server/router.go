package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/Barbra-GbR/barbra-backend/config"
	"github.com/Barbra-GbR/barbra-backend/controllers"
	"github.com/Barbra-GbR/barbra-backend/middlewares"
)

//Creates a new router and sets all paths, middlewares etc
func NewRouter() *gin.Engine {
	c := config.GetConfig();
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//Login-Flow cookie store
	store := sessions.NewCookieStore([]byte(c.GetString("server.cookie_store_secret")))
	router.Use(sessions.Sessions("login_session", store))

	//Controllers
	authenticationController := new(controllers.AuthenticationController)
	userController := new(controllers.UserController)
	suggestionController := new(controllers.SuggestionController)

	//----------------------------------------- Public routes -----------------------------------------
	public := router.Group("/api/v1/")

	//Authentication
	public.Handle(http.MethodGet, "/login/:provider/callback", authenticationController.AuthenticationCallback)
	public.Handle(http.MethodGet, "/login/:provider", authenticationController.Authenticate)

	//----------------------------------------- Private routes -----------------------------------------
	private := router.Group("/api/v1")
	private.Use(middlewares.AuthorizationMiddleware(false))

	//User profile
	private.Handle(http.MethodGet, "/user/me", userController.GetAccount)
	private.Handle(http.MethodPatch, "/user/me", userController.UpdateProfile)

	//----------------------------- Private routes (enrolled accounts only) -----------------------------
	enrolled := router.Group("/api/v1")
	enrolled.Use(middlewares.AuthorizationMiddleware(true))

	//Suggestions
	enrolled.Handle(http.MethodGet, "/suggestion", suggestionController.GetSuggestions)
	enrolled.Handle(http.MethodGet, "/suggestion/:id", suggestionController.GetSuggestion)

	//User bookmarks
	enrolled.Handle(http.MethodPost, "/user/me/bookmark", userController.AddBookmark)
	enrolled.Handle(http.MethodDelete, "/user/me/bookmark", userController.RemoveBookmark)

	return router
}
