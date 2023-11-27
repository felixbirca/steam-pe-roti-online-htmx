package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

type IdTokenRequst struct {
	IdToken   string
	CsrfToken string
}

type Location struct {
	Name string
	Date string
	Id   string
	Last bool
}

type FirestoreLocation struct {
	name string
	date string
}

func main() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "stem-pe-roti-online"}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		log.Printf("error initializing app: %v", err)
		return
	}

	// Initialize the Firebase Authentication client
	client, err := app.Auth(ctx)
	firestoreClient, err := app.Firestore(ctx)

	if err != nil {
		log.Panic(err.Error())
	}

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/**/*")

	router.StaticFile("/login/index.js", "./web/templates/login/index.js")

	router.GET("/posts/index", func(c *gin.Context) {
		tmpl := template.Must(template.New("main").Parse("layout/index.html"))
		tmpl = template.Must(tmpl.New("content").Parse("posts/index.html"))
		c.HTML(http.StatusOK, "main", gin.H{
			"title":    "Users",
			"user":     "Felix",
			"template": "posts/index.html",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/index.html", gin.H{
			"title": "Users",
		})
	})

	router.GET("/secured", func(c *gin.Context) {
		cookie, err := c.Request.Cookie("session")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		reqContext := c.Request.WithContext(c.Request.Context())
		decoded, err := client.VerifySessionCookieAndCheckRevoked(reqContext.Context(), cookie.Value)

		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		tmpl := template.Must(template.New("main").Parse("layout/index.html"))
		tmpl = template.Must(tmpl.Parse("secured/index.html"))

		c.HTML(http.StatusOK, "main", gin.H{
			"title": "Secured route",
			"user":  decoded,
		})
	})

	router.POST("/sessionLogin", func(c *gin.Context) {
		var req IdTokenRequst
		c.BindJSON(&req)

		expiresIn := time.Hour * 24 * 5

		if client == nil {
			log.Panicln("NILLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL")
		}

		reqContext := c.Request.WithContext(c.Request.Context())
		cookie, err := client.SessionCookie(reqContext.Context(), req.IdToken, expiresIn)

		if err != nil {
			log.Println(err.Error())
		}

		c.SetCookie("session", cookie, int(expiresIn.Seconds()), "", "", true, true)

		c.Redirect(http.StatusFound, "/secured")
	})

	router.POST("/create-location", func(c *gin.Context) {
		date, parseErr := time.Parse(time.DateOnly, c.PostForm("visitDate"))

		if parseErr != nil {
			c.Error(err)
			return
		}

		_, _, err := firestoreClient.Collection("locations").Add(ctx, map[string]interface{}{
			"name": c.PostForm("location"),
			"date": timestamppb.New(date),
		})

		if err != nil {
			c.Error(err)
			return
		}

		c.Redirect(http.StatusFound, "/visited-locations")
	})

	router.GET("/visited-locations", func(c *gin.Context) {
		page := c.Query("page")

		pageNumber, err := strconv.Atoi(page)

		if err != nil {
			pageNumber = 1
		}

		query := firestoreClient.Collection("locations").OrderBy("date", 2).Limit(5).Offset(5 * (pageNumber - 1))
		docs, err := query.Documents(ctx).GetAll()

		if err != nil {
			c.Error(err)
			return
		}

		locations := []Location{}

		for index, value := range docs {

			visitDate, _ := value.Data()["date"].(time.Time)

			log.Print(visitDate)

			loc := Location{
				Name: value.Data()["name"].(string),
				Date: visitDate.Local().Format(time.DateOnly),
				Id:   value.Ref.ID,
				Last: index == 4,
			}

			locations = append(locations, loc)
		}

		tmpl := template.Must(template.New("main").Parse("layout/index.html"))
		tmpl = template.Must(tmpl.Parse("visited-places/index.html"))

		if pageNumber > 1 {
			c.HTML(http.StatusOK, "visited-places/partial-places.html", gin.H{
				"locations": locations,
				"page":      pageNumber + 1,
			})
			return
		}

		c.HTML(http.StatusOK, "visited-places/index.html", gin.H{
			"locations": locations,
			"page":      pageNumber + 1,
		})
	})

	router.Run("localhost:8080")
}
