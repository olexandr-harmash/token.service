package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

func AuthorizeRequest(
	dumpvar bool,
	srv *server.Server,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", r)
		}

		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()
		//TODO redirect to client from vue

		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func TokenRequest(
	dumpvar bool,
	srv *server.Server,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
		}

		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func TestRequest(
	dumpvar bool,
	srv *server.Server,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "test", r) // Ignore the error
		}
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	}
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func UserAuthorizeHandler(
	dumpvar bool,
) func(http.ResponseWriter, *http.Request) (string, error) {
	return func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
		}

		store, err := session.Start(nil, w, r)
		if err != nil {
			return
		}

		uid, ok := store.Get("LoggedInUserID")
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			store.Set("ReturnUri", r.Form)
			store.Save()

			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}

		
		if _, ok = store.Get("Logged"); !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			userID = uid.(string)
		}
		
		store.Delete("LoggedInUserID")
		store.Delete("Logged")
		store.Delete("ConfirmCode")
		store.Save()
		return
	}
}

func LoginRequest(
	dumpvar bool,
	host string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "login", r) // Ignore the error
		}
		store, err := session.Start(nil, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "POST" {
			store.Set("LoggedInUserID", r.URL.Query().Get("username"))

			//TODO email confirm
			rand.Seed(time.Now().UnixNano())
			store.Set("ConfirmCode", strconv.FormatInt(int64(rand.Intn(10000)), 10))
			store.Save()
			fmt.Println(store.Get("ConfirmCode"))

			setCorsHeaders(w, host)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Location", host+"/")
		w.WriteHeader(http.StatusFound)
		return
	}
}

func AuthRequest(
	dumpvar bool,
	host string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "auth", r) // Ignore the error
		}
		store, err := session.Start(nil, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, l := store.Get("LoggedInUserID")
		code, c := store.Get("ConfirmCode")

		if !l || !c {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}

		//TODO email confirm
		if r.URL.Query().Get("code") != code {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		store.Set("Logged", true)
		store.Save()

		setCorsHeaders(w, host)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func setCorsHeaders(w http.ResponseWriter, host string) {
	w.Header().Add("Access-Control-Allow-Origin", host)
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}
