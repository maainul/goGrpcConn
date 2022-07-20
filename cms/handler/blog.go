package handler

import (
	"html/template"
	"net/http"
	"goGrpcConn/svcUtils/logging"
	"time"

	gk "goGrpcConn/api/gunk/v1/admin/blog"

	"github.com/gorilla/csrf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Blog struct {
	ID        string
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt time.Time
	DeletedBy string
}

type BlogTempData struct {
	CSRFField   template.HTML
	Form        Blog
	FormAction  string
	FormErrors  map[string]string
	FormMessage map[string]string
	FormName    string
}

func (s *Server) blogFormHandler(w http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).WithField("method", "BlogFormHander")
	data := BlogTempData{
		CSRFField:  csrf.TemplateField(r),
		Form:       Blog{},
		FormAction: blogCreate,
		FormName:   "blog-form.html",
	}
	s.loadBlogTemplate(w, r, data)
}

func (s *Server) blogCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "blogCreateHandler")
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var blog Blog
	if err := s.decoder.Decode(&blog, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		return
	}
	_, err := s.blog.CreateBlog(ctx, &gk.CreateBlogRequest{
		Blog: &gk.Blog{
			Title:     blog.Name,
			CreatedAt: timestamppb.Now(),
			CreatedBy: "12345",
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("create Blog failed")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		return
	}

	data := BlogTempData{
		CSRFField:  csrf.TemplateField(r),
		FormAction: blogCreate,
		FormName:   "blog-form.html",
	}
	s.loadBlogTemplate(w, r, data)
}

func (s *Server) loadBlogTemplate(w http.ResponseWriter, r *http.Request, data BlogTempData) {
	log := logging.FromContext(r.Context()).WithField("method", "loadBlogTemplate")
	tmpl := s.lookupTemplate(data.FormName)
	if tmpl == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
		return
	}
}
