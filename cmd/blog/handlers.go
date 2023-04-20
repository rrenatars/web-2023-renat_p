package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type indexPageData struct {
	FeaturedPosts []*featuredPostData
	MostRecent    []*mostRecentData
}

type featuredPostData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	PublishDate  string `db:"publish_date"`
	Author       string `db:"author"`
	AuthorAvatar string `db:"author_url"`
	Image        string `db:"image_url_modifier"`
	PostID       string `db:"post_id"`
	PostURL      string
}

type mostRecentData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	PublishDate  string `db:"publish_date"`
	Author       string `db:"author"`
	AuthorAvatar string `db:"author_url"`
	Image        string `db:"image_url"`
	PostID       string `db:"post_id"`
	PostURL      string
}

type postData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	Content      string `db:"content"`
	ImagePostURL string `db:"image_url"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		featuredPosts, err := featuredPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		recentPosts, err := mostRecent(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := indexPageData{
			FeaturedPosts: featuredPosts,
			MostRecent:    recentPosts,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"]

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post id", 403)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		err = ts.Execute(w, post)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func featuredPosts(db *sqlx.DB) ([]*featuredPostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url_modifier,
			post_id
		FROM
		    post
		WHERE
		    featured = 1
	`

	var featuredPosts []*featuredPostData

	err := db.Select(&featuredPosts, query)
	if err != nil { // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, post := range featuredPosts {
		post.PostURL = "/post/" + post.PostID
	}

	return featuredPosts, nil
}

func mostRecent(db *sqlx.DB) ([]*mostRecentData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			post_id
		FROM
		    post
		WHERE
		    featured = 0
	`

	var mostRecent []*mostRecentData

	err := db.Select(&mostRecent, query)
	if err != nil { // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, post := range mostRecent {
		post.PostURL = "/post/" + post.PostID
	}

	return mostRecent, nil
}

func postByID(db *sqlx.DB, postID int) (postData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			content,
			image_url
		FROM
			post
		WHERE
			post_id = ?
	`

	var post postData

	err := db.Get(&post, query, postID)
	if err != nil {
		return postData{}, err
	}

	return post, nil
}
