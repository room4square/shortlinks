package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"encoding/json"
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"strings"
)

type Links struct {
	UrlOriginal string `json:"urloriginal"`
	UrlShort string `json:"urlshort"`
}

type LinksAuto struct {
	UrlOriginal string `json:"urloriginal"`
	UrlShort string `json:"urlshort"`
}

var db sql.DB

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}




func reqShortURL(c echo.Context) error {
	links := Links{}

	// db conn
	db, err := sql.Open("sqlite3", "./links.db")
	checkErr(err)


	//fmt.Println(generateRandom())

	defer c.Request().Body.Close()
	errorJSON := json.NewDecoder(c.Request().Body).Decode(&links)
	if err != nil {
		log.Printf("Failed Processing Request Short URL: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorJSON.Error())
	}


	// check generate url by user
	var exists bool = false
	newURL := "mthr.ml/"+ strings.Replace(links.UrlShort, ".", "-", -1)

	err, exists = urlByUserCheck(links.UrlOriginal, newURL)
	if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errorJSON.Error())
	}else{
		if exists == true {
			return echo.NewHTTPError(http.StatusInternalServerError, "URL Input already exists")
		}
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO urlmap(url, code, created) values(?,?,?)")
	checkErr(err)
fmt.Println(links.UrlOriginal)
	res, err := stmt.Exec(string(links.UrlOriginal), string(newURL), time.Now())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	// query
	rows, err := db.Query("SELECT * FROM urlmap")
	checkErr(err)
	var url string
	var code string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&url, &code, &created)
		checkErr(err)
	}

	rows.Close()
	defer db.Close()

	return c.String(http.StatusOK, "Short URL: "+newURL)

}

func reqAutoShortURL(c echo.Context) error {
	linksAuto := LinksAuto{}

	// db conn
	db, err := sql.Open("sqlite3", "./links.db")
	checkErr(err)


	//fmt.Println(generateRandom())

	defer c.Request().Body.Close()
	errorJSON := json.NewDecoder(c.Request().Body).Decode(&linksAuto)
	if err != nil {
		log.Printf("Failed Processing Request Short URL: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorJSON.Error())
	}


	// check generate url by user
	var exists bool = false
	newURL := "mthr.ml/"+ generateRandom()
	fmt.Println(linksAuto.UrlOriginal)
	err, exists = urlByUserCheck(linksAuto.UrlOriginal, newURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errorJSON.Error())
	}else{
		if exists == true {
			return echo.NewHTTPError(http.StatusInternalServerError, "URL Input already exists "+linksAuto.UrlOriginal)
		}
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO urlmap(url, code, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(string(linksAuto.UrlOriginal), string(newURL), time.Now())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	// query
	rows, err := db.Query("SELECT * FROM urlmap")
	checkErr(err)
	var url string
	var code string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&url, &code, &created)
		checkErr(err)
	}

	rows.Close()
	defer db.Close()

	return c.String(http.StatusOK, "Short URL: "+newURL)

}

// check url created by user
func urlByUserCheck(urlOriginal, inputURL string) (e error, exists bool) {
	db, err := sql.Open("sqlite3", "./links.db")
	errQuery := db.QueryRow("SELECT EXISTS(SELECT url FROM urlmap WHERE url = ? or code= ?)", urlOriginal, inputURL).Scan(&exists)
	if err != nil {
		return errQuery, false
	}

	db.Close()
	return nil, exists
}

// check url created by auto
func urlByAutoCheck(urlOriginal string) (e error, exists bool) {

	db, err := sql.Open("sqlite3", "./links.db")
	errQuery := db.QueryRow("SELECT EXISTS(SELECT * FROM urlmap WHERE url = ? )", urlOriginal).Scan(&exists)
	if err != nil {
		return errQuery, false
	}

	db.Close()
	return nil, exists
}

func generateRandom() string {

	var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	s := make([]rune, 6)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}

	return string(s)
}



func main() {


	fmt.Printf("Server Start...")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Shortlinks API")
	})

	// short url with user defined short url
	e.POST("/doshort", reqShortURL)

	// short url with user defined short url
	e.POST("/doautoshort", reqAutoShortURL)

	// start servers
	e.Start(":8000")


}