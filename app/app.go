package app

import (
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/logout"
	_ "github.com/volatiletech/authboss/v3/register"
	"html/template"
	"log"
	"net/http"
	"task_ws_et/config"
	"task_ws_et/models"
	"task_ws_et/storage"
)

type App struct {
	Port     string
	Router   *chi.Mux
	Auth 	 *authboss.Authboss
}

func New(storage *storage.Storage, conf config.Config) (*App, error) {
	a := &App{
		Port:	conf.HttpPort,
	}

	auth, err := newAuth(storage, conf)
	if err != nil {
		return nil, errors.Wrap(err, "creating new auth err")
	}

	a.Auth = auth
	a.initRouter()

	return a, nil
}

func (a *App) initRouter() {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger, a.Auth.LoadClientStateMiddleware, middleware.Recoverer)

	//only auth users
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(a.Auth, authboss.RequireNone, authboss.RespondUnauthorized))
		mux.Group(func(mux chi.Router) {
			mux.Use(a.onlyAdmin)
			mux.Method("GET", "/sigma", a.handler(a.index))
		})
		mux.Method("GET","/foo", a.handler(a.index))
		mux.Method("GET", "/bar", a.handler(a.index))
	})

	//authboss auth pages
	mux.Group(func(mux chi.Router) {
		mux.Mount("/auth", http.StripPrefix("/auth", a.Auth.Config.Core.Router))
	})

	mux.Method("GET","/", a.handler(a.main))

	a.Router = mux
}

func (a *App) handler(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if aerr, ok := err.(*AuthErr); ok {
				log.Printf("auth err: %s", aerr)
				http.Error(w, aerr.Message, aerr.Code)
			} else {
				log.Printf("internal server err: %v\n", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}

func (a App) GetCurrentUser(r *http.Request) (*models.User, error) {
	iuser, err := a.Auth.CurrentUser(r)
	if err != nil {
		return nil, fmt.Errorf("can`t get current user")
	}
	user, ok := iuser.(*models.User)
	if !ok {
		return nil, fmt.Errorf("can`t cast to user")
	}
	return user, nil
}

func (a *App) index(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Hello from " + r.RequestURI))
	return nil
}

type Item struct {
	Href 	string
	Name 	string
}

func (a *App) main(w http.ResponseWriter, r *http.Request) error {
	user, _ := a.GetCurrentUser(r)

	data := struct{
		LogIn bool
		User  *models.User
		Port  string
		Items []Item
	}{
		LogIn: user != nil,
		User:	user,
		Port: 	a.Port,
		Items: 	[]Item{
			{"/foo","foo"},
			{"/bar","bar"},
			{"/sigma","sigma"},
		},
	}
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return errors.Wrap(err, "parsing file err")
	}
	if err := tmpl.Execute(w, data); err != nil {
		return errors.Wrap(err, "execute template err")
	}
	return nil
}

type AuthErr struct {
	Code 	int
	Message string
}

func (ae AuthErr) Error() string {
	return ae.Message
}