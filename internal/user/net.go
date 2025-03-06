package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/util"

	"golang.org/x/crypto/bcrypt"
)

// CanProceed returns `true` if the user in `rq` has enough rights to access `route`.
func CanProceed(rq *http.Request, route string) bool {
	return FromRequest(rq).CanProceed(route)
}

// FromRequest returns user from `rq`. If there is no user, an anon user is returned instead.
func FromRequest(rq *http.Request) *User {
	username, ok := rq.Header["X-Webauth-User"]
	if !ok || len(username) < 1 {
		return EmptyUser()
	}
	return ByName(username[0])
}

// LogoutFromRequest logs the user in `rq` out and rewrites the cookie in `w`.
func LogoutFromRequest(w http.ResponseWriter, rq *http.Request) {
	cookieFromUser, err := rq.Cookie("mycorrhiza_token")
	if err == nil {
		http.SetCookie(w, cookie("token", "", time.Unix(0, 0)))
		terminateSession(cookieFromUser.Value)
	}
}

// Register registers the given user. If it fails, a non-nil error is returned.
func Register(username, password, group, source string, force bool) error {
	if !IsValidUsername(username) {
		return fmt.Errorf("illegal username ‘%s’", username)
	}
	username = util.CanonicalName(username)

	switch {
	case !IsValidUsername(username):
		return fmt.Errorf("illegal username ‘%s’", username)
	case !ValidGroup(group):
		return fmt.Errorf("invalid group ‘%s’", group)
	case !ValidSource(source):
		return fmt.Errorf("invalid source ‘%s’", source)
	case HasUsername(username):
		return fmt.Errorf("username ‘%s’ is already taken", username)
	case !force && cfg.RegistrationLimit > 0 && Count() >= cfg.RegistrationLimit:
		return fmt.Errorf("reached the limit of registered users (%d)", cfg.RegistrationLimit)
	case password == "" && source != "telegram":
		return fmt.Errorf("password must not be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := User{
		Name:         username,
		Group:        group,
		Source:       source,
		Password:     string(hash),
		RegisteredAt: time.Now(),
	}
	users.Store(username, &u)
	return SaveUserDatabase()
}

var (
	ErrUnknownUsername = errors.New("unknown username")
	ErrWrongPassword   = errors.New("wrong password")
)

// LoginDataHTTP logs such user in and returns string representation of an error if there is any.
//
// The HTTP parameters are used for setting header status (bad request, if it is bad) and saving a cookie.
func LoginDataHTTP(w http.ResponseWriter, username, password string) error {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if !HasUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("Unknown username entered", "username", username)
		return ErrUnknownUsername
	}
	if !CredentialsOK(username, password) {
		w.WriteHeader(http.StatusBadRequest)
		slog.Info("Wrong password entered", "username", username)
		return ErrWrongPassword
	}
	token, err := AddSession(username)
	if err != nil {
		slog.Error("Failed to add session", "username", username, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	http.SetCookie(w, cookie("token", token, time.Now().Add(365*24*time.Hour)))
	return nil
}

// AddSession saves a session for `username` and returns a token to use.
func AddSession(username string) (string, error) {
	token, err := util.RandomString(16)
	if err == nil {
		commenceSession(username, token)
		slog.Info("Added session", "username", username)
	}
	return token, err
}

// A handy cookie constructor
func cookie(nameSuffix, val string, t time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "mycorrhiza_" + nameSuffix,
		Value:   val,
		Expires: t,
		Path:    "/",
	}
}
