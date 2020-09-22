package app

import (
	"encoding/base64"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
	"net/http"
	"regexp"
	"task_ws_et/config"
	"task_ws_et/storage"
	"time"
)

func newAuth(storage *storage.Storage, conf config.Config) (*authboss.Authboss, error) {
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(conf.CookieKey)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(conf.SessionKey)

	cookieStore := abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false

	sessionStore := abclientstate.NewSessionStorer("auth_token", sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	ab := authboss.New()

	ab.Config.Paths.RootURL = fmt.Sprintf("http://localhost:%s", conf.HttpPort)

	ab.Config.Modules.LogoutMethod = "GET"

	ab.Config.Storage.Server = storage
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "views")

	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	defaults.SetCore(&ab.Config, false, false)

	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	roleRule := defaults.Rules{
		FieldName: "role", Required: true,
		MatchError: "Role must be admin or user",
		MustMatch:  regexp.MustCompile(`^(admin|user)$`),
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: false,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, roleRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": {"email", "role", "password"},
		},
	}

	if err := ab.Init(); err != nil {
		return nil, errors.Wrap(err, "authboss init err")
	}

	return ab, nil
}

func (a App) onlyAdmin(h http.Handler) http.Handler {
	return a.handler(func(w http.ResponseWriter, r *http.Request) error {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		user, err := a.GetCurrentUser(r)
		if err != nil {
			return &AuthErr{http.StatusForbidden,"forbidden, can`t get user role"}
		}
		if user.Role != "admin" {
			return &AuthErr{http.StatusForbidden,"forbidden! need to be admin"}
		}

		h.ServeHTTP(w, r)
		return nil
	})
}