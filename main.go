package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	md "github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	"github.com/olexandr-harmash/token.service/rest"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var (
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
	host      string
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "222222", "The client id being passed in")
	flag.StringVar(&secretvar, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
	flag.StringVar(&host, "h", "http://localhost:8080", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {

	pgxConn, err := pgx.Connect(context.TODO(), "postgres://test:test@localhost:5432/test?sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("ok")
	adapter := pgx4adapter.NewConn(pgxConn)

	tokenStore, err := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		log.Fatalf("cannot init token server reason: %s", err.Error())
	}

	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		log.Fatalf("cannot init token server reason: %s", err.Error())
	}
	//TODO added client services
	clientStore.Create(&md.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	})

	//TODO save user data for gen/regenerate token
	// userStore, err := user.NewStore(adapter)
	// if err != nil {
	// 	log.Fatalf("cannot init token server reason: %s", err.Error())
	// }

	manager := manage.NewDefaultManager()

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetUserAuthorizationHandler(nil)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	srv.SetUserAuthorizationHandler(rest.UserAuthorizeHandler(dumpvar))

	http.HandleFunc("/login", rest.LoginRequest(dumpvar, host))
	http.HandleFunc("/auth", rest.AuthRequest(dumpvar, host))
	http.HandleFunc("/oauth/authorize", rest.AuthorizeRequest(dumpvar, srv))
	http.HandleFunc("/oauth/token", rest.TokenRequest(dumpvar, srv))
	http.HandleFunc("/test", rest.TestRequest(dumpvar, srv))

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}
