package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/controllers"
	"github.com/michaeltintiuc/shackle-api/pkg/middleware"
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
}

type DbInfo struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func NewApp(port string, dbInfo DbInfo, jwtSecret string) (*app, error) {
	a := &app{}
	db, err := connectAndVerifyDb(dbInfo)

	if err != nil {
		return nil, err
	}

	a.db = db.Database(dbInfo.Name)
	a.router = mux.NewRouter().StrictSlash(true)
	a.apiRouter = a.router.PathPrefix("/api").Subrouter()
	a.authRouter = a.router.PathPrefix("/auth").Subrouter()

	a.router.Use(middleware.Log, middleware.Json)
	a.apiRouter.Use(middleware.Auth(jwtSecret))
	controllers.Init(a, jwtSecret)

	a.server = &http.Server{
		Handler:      a.router,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return a, err
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
	if err := a.server.ListenAndServe(); err != nil {
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
