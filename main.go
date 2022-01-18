package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DevHostURL                         = "localhost:8080"
	ARTICLES_FILE                      = "articles.json"
	MAX_ARTICLES_OF_TAGNAME_DATE_QUERY = 10
)

// R4(Requirement#4)'s implementation
// article represents data about a record article.
//!+
type Article struct {
	Id        string
	Title     string
	Date      string // supposed to be publishing date
	Body      string
	EntryTime string
	Tags      []string
}

// R6's implementation
type TagNameDateSummary struct {
	Tag         string
	Count       int
	Articles    []string
	RelatedTags []string
}

// R9's implementation
type IdEntryTime struct {
	EntryTime time.Time
	Id        string
}

type Articles struct {
	Articles []Article `json:"articles"`
}

// R5's implementation
var articles Articles

// R4's implementation
// getArticles responds with the list of all articles as JSON.
func getArticles(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, articles)
}

// R2's implementation
// getArticleByID locates the article whose ID value matches the id
// parameter sent by the client, then returns that article as a response.
func getArticleByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of articles, looking for
	// an article whose ID value matches the parameter.
	for _, a := range articles.Articles {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "article not found"})
}

// R3's implementation
// getArticleByTagDate locates the articles whose Tag and Date values matche the those
// sent by the client, then returns that articles as response.
func getArticleByTagDate(c *gin.Context) {
	tagName := c.Param("tagName")
	dt := c.Param("date")

	fmt.Println("received parameter tagName=", tagName)
	fmt.Println("received parameter date=", dt)

	// Loop over the list of articles, looking for
	// articles whose tagname and date matched the parameters.
	var foundMatched bool = false
	var tagNameDateStats TagNameDateSummary // R6's implemenation

	tagNameDateStats.Tag = tagName
	tagNameDateStats.Count = 0
	var related_tags = make(map[string]int)
	var idEntryTime []IdEntryTime

	for _, a := range articles.Articles {
		if a.Date == dt {
			// collect EntryTime of underlying article
			iet := new(IdEntryTime)
			t, _ := time.Parse(time.RFC3339, a.EntryTime)
			iet.EntryTime = t
			iet.Id = a.Id
			idEntryTime = append(idEntryTime, *iet)
			fmt.Println("a.Id=", a.Id, " a.EntryTime=", a.EntryTime, " coverted time=", t)
			//idEntryTime = a.Id

			// ok, found matching date ...
			fmt.Println("matched Id=", a.Id, " and Date=", a.Date)
			for _, t := range a.Tags {
				// how about tag?
				if t == tagName {
					fmt.Println("matched tag=", t)
					foundMatched = true
				} else {
					// R7's implementation
					// handle the rest of tags associated with the article
					// note: counting underlying related_tag for future use, just in case
					related_tags[t]++
				}
				// R8's implementation
				tagNameDateStats.Count = len(related_tags) + 1 // the additional 1 is for the inquiring tagName itself
			}
		}
	}
	if !foundMatched {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "article not found"})
	} else {
		// don't forget the rest of tags associated with underlying article.
		for tag, _ := range related_tags {
			tagNameDateStats.RelatedTags = append(tagNameDateStats.RelatedTags, tag)
		}
		//c.IndentedJSON(http.StatusOK, tagNameDateStats)
	}

	// R10's implementation
	fmt.Println("idEntryTime size=", len(idEntryTime))
	sort.Slice(idEntryTime, func(i, j int) bool {
		return idEntryTime[i].EntryTime.After(idEntryTime[j].EntryTime)
	})
	for _, i := range idEntryTime {
		fmt.Println("after sorting, EntryTime=", i.EntryTime, " Id=", i.Id)
	}

	limit := 0
	if len(idEntryTime) >= MAX_ARTICLES_OF_TAGNAME_DATE_QUERY {
		limit = MAX_ARTICLES_OF_TAGNAME_DATE_QUERY
	} else {
		limit = len(idEntryTime)
	}
	// R9's implementation
	// transfer these article Ids to tagNameDateStats.Articles within limit
	for i := limit - 1; i >= 0; i-- {
		fmt.Println("idEntryTime[", i, "]=", idEntryTime[i].EntryTime, idEntryTime[i].Id)
		tagNameDateStats.Articles = append(tagNameDateStats.Articles, idEntryTime[i].Id)
	}

	c.IndentedJSON(http.StatusOK, tagNameDateStats)
}

// R1's implementation
// postAlbums adds an album from JSON received in the request body.
func postArticles(c *gin.Context) {
	var article Article

	if err := c.BindJSON(&article); err != nil {
		return
	}

	// Add the new album to the slice.
	articles.Articles = append(articles.Articles, article)
	c.IndentedJSON(http.StatusCreated, article)
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open(ARTICLES_FILE)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Successfully Opened " + ARTICLES_FILE)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &articles)

	router := gin.Default()
	router.POST("/articles", postArticles)                  // R1's implementation
	router.GET("/articles/:id", getArticleByID)             // R2's implementation
	router.GET("/tags/:tagName/:date", getArticleByTagDate) // R3's implementation
	router.GET("/articles/all", getArticles)                // R4's implementation, an extra added by SK
	router.Run(DevHostURL)
}
