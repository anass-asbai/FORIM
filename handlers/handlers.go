package handlers

import (
	"fmt"
	"net/http"

	"forim/bcryptp"
	"forim/database"
)

var limit = 0

func GetHome(w http.ResponseWriter, r *http.Request) {
	catigorie := r.FormValue("category")
	action := r.FormValue("Next")
	if action != "" && database.CountPost(limit+1) {
		limit += 5
	}
	
	action = r.FormValue("Back")
	if action != "" && limit != 0 {
		limit -= 5
	}
	ch := 0
	cookie, err := r.Cookie("session")
		if err != nil {
			ch = 1
		}
	posts, _ := database.GetPosts(catigorie, limit, ch)
	id_post := r.FormValue("id-post")
	comment := r.FormValue("comment")
	if len(comment) > 200 {
		http.Error(w, "comment is too long", http.StatusBadRequest)
		return
	}

	if comment != "" {
		if err := database.Createcomment(comment, id_post, cookie.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Println(posts)
	RenderTemplate(w, "./assets/templates/post.html", posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if title != "" && content != "" && category != "" {
			if len(title) < 5 || len(title) > 50 {
				http.Error(w, "title is too long or too short", http.StatusBadRequest)
				return
			}
			if len(content) < 10 || len(content) > 500 {
				http.Error(w, "content is too long or too short", http.StatusBadRequest)
				return
			}
			if err := database.InsertPost(title, content, cookie.Value, category); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println(comments)
	RenderTemplate(w, "./assets/templates/comment.html", comments)
}

func NewComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id_post := r.FormValue("id-post")
	comment := r.FormValue("comment")
	if len(comment) > 200 {
		http.Error(w, "comment is too long", http.StatusBadRequest)
		return
	}
	
	if comment != "" {
		cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
		if err := database.Createcomment(comment, id_post, cookie.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Like_post(w http.ResponseWriter, r *http.Request) {
	post_like := r.FormValue("like_post")
	post_deslike := r.FormValue("deslike_post")
	comment_like := r.FormValue("like_comment")
	comment_deslike := r.FormValue("deslike_comment")


	if post_like != "" {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err = database.InsertLike(post_like, cookie.Value, true,true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if post_deslike != "" {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err = database.InsertLike(post_deslike, cookie.Value, false,true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}else if comment_like != "" {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err = database.InsertLike(comment_like, cookie.Value, true,false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		/*
			hna 5aliwa trja3 lhome 7ta ndiro link fih query ikon fih id dyal post
			bax wa9tama darna like irja3 lblasto fnafs lpost
		*/
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err = database.InsertLike(comment_deslike, cookie.Value, false,false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("=",password,"=",email,"=",doz)

	if doz  {
		cookie := http.Cookie{
			Name:  "session",
			Value: email,
		}
		http.SetCookie(w, &cookie)
		catigorie := r.FormValue("category")
		posts, err := database.GetPosts(catigorie, 0,0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		fmt.Println("----",cookie,"===",posts)
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		RenderTemplate(w, "./assets/templates/register.html", nil)
	}
}