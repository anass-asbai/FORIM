package handlers

import (
	"forim/database"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, "./assets/templates/post.html", posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title != "" && content != "" {
			if err := database.InsertPost(title, content); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/post", http.StatusSeeOther)
			return
		}
	}

	RenderTemplate(w, "./assets/templates/post.create.page.html", nil)
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
		posts, err := database.GetPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

		if email == "" || password == "" || name == "" {
			RenderTemplate(w, "./assets/templates/register.html", nil)
			return
		}
		if err := database.CreateAcount(name, email, password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		RenderTemplate(w, "./assets/templates/register.html", nil)
	}
}
