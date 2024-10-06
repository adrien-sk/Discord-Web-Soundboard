package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	discordOAuth "github.com/adrien-sk/Discord-Web-Soundboard/pkg/discordOauth"
	"github.com/adrien-sk/Discord-Web-Soundboard/pkg/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

// ---------------------------------------------- Variables and Structs

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

// ---------------------------------------------- Main function

func main() {
	fmt.Println("Starting Go server")

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	pool := websocket.NewPool()
	go pool.Start()

	router.GET("/sounds", getSoundsHandler)
	router.GET("/auth", authHandler)
	router.GET("/auth/callback", authCallbackHandler)
	router.GET("/auth/isauthenticated", authIsAuthenticatedHandler)
	router.GET("/ws", func(c *gin.Context) {
		serveWs(pool, c)
	})

	router.Run("localhost:8080")
}

// ---------------------------------------------- API Handlers

func getSoundsHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sounds)
}

// Step 1: Redirect to the OAuth 2.0 Authorization page.
func authHandler(c *gin.Context) {
	http.Redirect(c.Writer, c.Request, conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// Step 2: After user authenticates their accounts this callback is fired.
func authCallbackHandler(c *gin.Context) {
	fmt.Println("-- auth Callback Handler")

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
	fmt.Println("-- auth IsAuthenticated Callback Handler")

	cookie, err := c.Cookie("discord.oauth2")

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	userSession, exists := sessions[cookie]
	if !exists {
		fmt.Printf("IsAuthenticated for : " + userSession.Username)
		// If the session token is not present in session map, return an unauthorized error
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Print(sessions)

	c.Writer.WriteHeader(http.StatusOK)
}

// ---------------------------------------------- Websockets Methods

func serveWs(pool *websocket.Pool, c *gin.Context) {
	fmt.Println("WebSocket Endpoint Hit")
	cookie, err := c.Cookie("discord.oauth2")
	userSession, exists := sessions[cookie]
	if err != nil || !exists {
		fmt.Printf("serveWs fail for : " + userSession.Username)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintf(c.Writer, "%+v\n", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

// ---------------------------------------------- Discord bot methods
