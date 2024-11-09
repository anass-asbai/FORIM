package handlers

import (
	"net/http"

	"forim/bcryptp"
	"forim/database"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	catigorie := r.FormValue("category")
	posts, err := database.GetPosts(catigorie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id_post := r.FormValue("id-post")
	comment := r.FormValue("comment")
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if cookie == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if comment != "" {
		if err := database.Createcomment(comment, id_post, cookie.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	RenderTemplate(w, "./assets/templates/post.html", posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if title != "" && content != "" && category != "" {
			if err := database.InsertPost(title, content, cookie.Value, category); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/post", http.StatusSeeOther)
			return
		}
	}

	RenderTemplate(w, "./assets/templates/post.create.page.html", nil)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	id_post := r.FormValue("id-post")
	comments, err := database.GetComment(id_post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RenderTemplate(w, "./assets/templates/comment.html", comments)
}

func Like_post(w http.ResponseWriter, r *http.Request) {
	like := r.FormValue("like_post")
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = database.InsertLike(like, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")
	doz, err := database.GetLogin(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if doz == true {
		catigorie := r.FormValue("category")
		posts, err := database.GetPosts(catigorie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:  "session",
			Value: email,
		}
		http.SetCookie(w, &cookie)
		RenderTemplate(w, "./assets/templates/post.html", posts)
	} else {
		errorMessage := ""
		if email != "" || password != "" {
			errorMessage = "Password or email not working"
		}
		RenderTemplate(w, "./assets/templates/index.html", errorMessage)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		name := r.FormValue("username")
		p, err := bcryptp.HashPassword(password)
		if email == "" || p == "" || name == "" {
			RenderTemplate(w, "./assets/templates/register.html", nil)
			return
		}
		if err != nil {
			RenderTemplate(w, "./assets/templates/register.html", nil)
			return
		}
		if err := database.CreateAcount(name, email, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		RenderTemplate(w, "./assets/templates/register.html", nil)
	}
}
