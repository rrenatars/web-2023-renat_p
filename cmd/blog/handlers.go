package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	CardImage    string `db:"card_image"`
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
	CardImage    string `db:"card_image"`
	Image        string `db:"image_url"`
	PostID       string `db:"post_id"`
}

type postData struct {
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	Content      string `db:"content"`
	ImagePostURL string `db:"image_url"`
}

type createPostRequest struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	AuthorName            string `json:"author_name"`
	AuthorAvatar          string `json:"author_avatar"`
	AvatarFileName        string `json:"avatar_file_name"`
	PublishDate           string `json:"publish_date"`
	BigHeroimage          string `json:"big_heroimage"`
	BigHeroimageFileName  string `json:"big_heroimage_file_name"`
	CardHeroimage         string `json:"small_heroimage"`
	CardHeroimageFileName string `json:"small_heroimage_file_name"`
	PostContent           string `json:"content"`
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

func admin(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/admin.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func loginlogin(w http.ResponseWriter, r *http.Request) {
    ts, err := template.ParseFiles("pages/login.html")
    if err != nil {
        http.Error(w, "Internal Server Error", 500)
	log.Println(err. Error())
	return
    }

    err = ts.Execute(w, nil)
    if err != nil {
	    http.Error(w, "Internal Server Error", 500)
	    log.Println(err. Error())
	    return
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
			card_image,
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
			card_image,
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

func createPost(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		var req createPostRequest

		err = json.Unmarshal(reqData, &req)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		err = savePost(db, req)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func savePost(db *sqlx.DB, req createPostRequest) error {
	authorImageUrl, err := saveImage(req.AvatarFileName, req.AuthorAvatar)
	if err != nil {
		log.Println(err)
		return err
	}

	heroimageUrl, err := saveImage(req.BigHeroimageFileName, req.BigHeroimage)
	if err != nil {
		log.Println(err)
		return err
	}

	cardImageUrl, err := saveImage(req.CardHeroimageFileName, req.CardHeroimage)
	if err != nil {
		log.Println(err)
		return err
	}

	const query = `
		INSERT INTO post
		(
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			card_image,
			image_url,
			content
		)
		VALUES
		(
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
	`

	_, err = db.Exec(query, req.Title, req.Description, req.AuthorName, authorImageUrl, req.PublishDate, cardImageUrl, heroimageUrl, req.PostContent)

	return err
}

func saveImage(imageName string, imageDecoded string) (string, error) {
	imageEncoded, err := base64.StdEncoding.DecodeString(imageDecoded)
	if err != nil {
		log.Println(err)
		return "", err
	}

	imageFile, err := os.Create("static/images/" + imageName)
	if err != nil {
		log.Println(err)
		return "", err
	}

	_, err = imageFile.Write(imageEncoded)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return "/static/images/" + imageName, err
}
