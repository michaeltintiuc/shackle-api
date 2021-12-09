package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/controllers"
	"github.com/michaeltintiuc/shackle-api/pkg/middleware"
	"github.com/michaeltintiuc/shackle-api/pkg/session"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type app struct {
	server     *http.Server
	db         *mongo.Database
	router     *mux.Router
	apiRouter  *mux.Router
	authRouter *mux.Router
	session    *session.Session
}

type DbInfo struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type spaHandler struct {
	path string
}

func NewApp(port string, dbInfo DbInfo, sessionInfo session.SessionInfo) (*app, error) {
	a := &app{}
	db, err := connectAndVerifyDb(dbInfo)

	if err != nil {
		return nil, err
	}

	a.db = db.Database(dbInfo.Name)
	a.session = session.NewSession(sessionInfo)
	a.router = mux.NewRouter().StrictSlash(true)

	a.apiRouter = a.router.PathPrefix("/api").Subrouter()
	a.authRouter = a.router.PathPrefix("/auth").Subrouter()

	spa := spaHandler{"../client/build/web"}
	a.router.PathPrefix("/").Handler(spa)

	a.router.Use(middleware.Log, middleware.Csrf(), middleware.AddBaseHeaders)
	a.authRouter.Use(middleware.AddJsonHeaders)
	a.apiRouter.Use(middleware.AddJsonHeaders, middleware.Auth(a.session))

	// Preflight requests and CORS
	// a.apiRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// middleware.AddCorsHeaders(w, r)
	// w.WriteHeader(http.StatusOK)
	// })

	controllers.Init(a)

	a.server = &http.Server{
		Handler:      a.router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return a, err
}

func (a *app) Session() *session.Session {
	return a.session
}

func (a *app) Db() *mongo.Database {
	return a.db
}

func (a *app) ApiRouter() *mux.Router {
	return a.apiRouter
}

func (a *app) AuthRouter() *mux.Router {
	return a.authRouter
}

func (a *app) ListenAndServe() {
	log.Printf("Listening on %s\n", a.server.Addr)
	if err := a.server.ListenAndServeTLS("../certs/shackle.dev.pem", "../certs/shackle.dev-key.pem"); err != nil {
		log.Println(err)
	}
}

func (a *app) Shutdown() {
	log.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		a.disconnectDb()
		cancel()
	}()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	log.Println("Server gracefully stopped")
}

func (a *app) disconnectDb() {
	log.Println("Disconnecting DB")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := a.db.Client().Disconnect(ctx); err != nil {
		log.Println(err)
	}
}

func connectAndVerifyDb(dbInfo DbInfo) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Println("Connecting to DB")
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%s", dbInfo.User, dbInfo.Pass, dbInfo.Host, dbInfo.Port)),
	)
	if err == nil {
		err = db.Ping(ctx, readpref.Primary())
	}

	return db, err
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.path, path)
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/" {
		t, err := template.ParseFiles(filepath.Join(path, "index.html"))
		if utils.HasError(w, err, "Application error", http.StatusInternalServerError) {
			return
		}
		err = t.Execute(w, map[string]interface{}{"csrf": csrf.Token(r)})
		if utils.HasError(w, err, "Application error", http.StatusInternalServerError) {
			return
		}
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.path)).ServeHTTP(w, r)
}
