package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	discordOAuth "github.com/adrien-sk/Discord-Web-Soundboard/pkg/discordOauth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

type Sound struct {
	Name string `json:"name"`
	Ext  string `json:"ext"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"global_name"`
	SessionID string
}

var sounds = []Sound{
	{Name: "Prout", Ext: "mp3"},
	{Name: "Ping", Ext: "mp3"},
}

var conf = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/callback",
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	Scopes:       []string{discordOAuth.ScopeIdentify, discordOAuth.ScopeEmail, discordOAuth.ScopeGuilds},
	Endpoint:     discordOAuth.Endpoint,
}

var sessions = map[string]User{}

var state = "random"

func main() {
	print("Starting Go server")

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	router.GET("/sounds", getSoundsHandler)
	router.GET("/auth", authHandler)
	router.GET("/auth/callback", authCallbackHandler)
	router.GET("/auth/isauthenticated", authIsAuthenticatedHandler)

	router.Run("localhost:8080")
}

func getSoundsHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sounds)
}

// Step 1: Redirect to the OAuth 2.0 Authorization page.
func authHandler(c *gin.Context) {
	http.Redirect(c.Writer, c.Request, conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// Step 2: After user authenticates their accounts this callback is fired.
func authCallbackHandler(c *gin.Context) {
	println("-- auth Callback Handler")

	// State verification before continuing
	if c.Request.FormValue("state") != state {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte("State does not match."))
		return
	}

	// Step 3: We exchange the code we got for an access token
	// Then we can use the access token to do actions, limited to scopes we requested
	token, err := conf.Exchange(context.Background(), c.Request.FormValue("code"))

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	// Step 4: Use the access token, here we use it to get the logged in user's info.
	res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
		} else {
			c.Writer.Write([]byte(res.Status))
		}
		return
	}

	userData := User{}
	json.NewDecoder(res.Body).Decode(&userData)
	userData.SessionID = uuid.NewString()

	if len(userData.ID) <= 0 || len(userData.Username) <= 0 {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessions[userData.SessionID] = userData

	c.SetCookie("discord.oauth2", userData.SessionID, 3600, "/", "http://localhost:5173/", false, true)

	http.Redirect(c.Writer, c.Request, "http://localhost:5173/", http.StatusTemporaryRedirect)

}

func authIsAuthenticatedHandler(c *gin.Context) {
	println("-- auth IsAuthenticated Callback Handler")

	cookie, err := c.Cookie("discord.oauth2")

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	userSession, exists := sessions[cookie]
	if !exists {
		log.Printf("IsAuthenticated for : " + userSession.Username)
		// If the session token is not present in session map, return an unauthorized error
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Print(sessions)

	c.Writer.WriteHeader(http.StatusOK)
}
