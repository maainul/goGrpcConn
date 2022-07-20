package handler

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	blgk "goGrpcConn/api/gunk/v1/admin/blog"
	"goGrpcConn/svcUtils/logging"
	"goGrpcConn/svcUtils/mw"

	"github.com/Masterminds/sprig"
	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Server struct {
	templates *template.Template
	env       string
	config    *viper.Viper
	logger    *logrus.Entry
	cookies   *sessions.CookieStore
	decoder   *schema.Decoder
	assetFS   *hashfs.FS
	assets    fs.FS
	blog      blgk.BlogServiceClient
}

func NewServer(
	env string,
	config *viper.Viper,
	logger *logrus.Entry,
	cookies *sessions.CookieStore,
	decoder *schema.Decoder,
	assets fs.FS,
	api *grpc.ClientConn,
) (*mux.Router, error) {
	s := Server{
		env:     env,
		config:  config,
		logger:  logger,
		cookies: cookies,
		decoder: decoder,
		assets:  assets,
		blog:    struct{ blgk.BlogServiceClient }{BlogServiceClient: blgk.NewBlogServiceClient(api)},
	}
	if err := s.parseTemplates(); err != nil {
		logging.NewLogger().Print("error while parse template")
		return nil, err
	}
	csrfSecure := config.GetBool("csrf.secure")
	csrfSecret := config.GetString("csrf.secret")
	if csrfSecret == "" {
		return nil, errors.New("CSRF secret must not be empty")
	}
	r := mux.NewRouter()
	mw.ChainHTTPMiddleware(r, logger, mw.CSRF([]byte(csrfSecret), csrf.Secure(csrfSecure), csrf.Path("/")))
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", cacheStaticFiles(http.FileServer(http.FS(s.assetFS)))))
	handlers(r, s)
	return r, nil
}

func cacheStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if asset is hashed extend cache to 180 days
		e := `"4FROTHS24N"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age=15552000")
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) lookupTemplate(name string) *template.Template {
	if s.env == "development" {
		if err := s.parseTemplates(); err != nil {
			s.logger.WithError(err).Error("template reload")
			return nil
		}
	}
	return s.templates.Lookup(name)
}

func (s *Server) parseTemplates() error {
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/")))
		},
	}).Funcs(sprig.FuncMap())
	tmpl, err := templates.ParseFS(s.assets, "templates/*/*/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}
