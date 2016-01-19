package main

import (
  "fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/gorilla/context"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"

  "github.com/garyburd/redigo/redis"

)


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

/*
func helloHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
*/

type appContext struct {
  db *sql.DB
	config *config
  redispool *redis.Pool
}

func main() {

	// Options Parse

	// Config Stage
	config, err := loadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

  /*
  // DB Stage
	db, err := sql.Open("mysql", config.Mysqlconfig.Connectstring)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db open and pinged")
	defer log.Println("closingdb")
	defer db.Close()
  */

  // Redis Stage
  log.Println("redis pool opening")
  redispool := newRedisPool(config.Redisconfig.Host, config.Redisconfig.Password)
  log.Println("redis open")
  defer log.Println("redis closing")
  defer redispool.Close()

	// app context
  //db: db,
	appC := appContext{ config: config, redispool: redispool }
  log.Println("app ready")

	// common routes
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)

	// routers
	router := httprouter.New()
	router.GET("/", wrapHandler(commonHandlers.ThenFunc(indexHandler)))
	// router.GET("/hello/:name", helloHandler)
	router.GET("/api/1/node/:nodekey/:nodekeyvalue", appC.httpGetNodeByFieldHandler)
	router.GET("/api/1/nodes/:nodekey/:nodekeyvalue", appC.httpGetNodesByFieldHandler)
	router.GET("/api/1/discoverednode/:nodekey/:nodekeyvalue", appC.httpGetDiscoveredNodeByFieldHandler)
	router.GET("/api/1/discoverednodes/:nodekey/:nodekeyvalue", appC.httpGetDiscoveredNodesByFieldHandler)


  /*
	router.GET("/api/1/nodes/:nodekey/:nodekeyvalue", xx)
	router.GET("/api/1/node/:nodekey/:nodekeyvalue", xx)
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/:field", xx)
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/ipxe", xx)
	router.GET("/api/1/discover/uuid/:uuid/ipv4address/:ipv4address/macaddress/:macaddress", xx)
	router.GET("/api/1/me/node", xx)
	router.GET("/api/1/me/node/:field", xx)
	router.GET("/api/1/os/:os_name/step/:os_step", xx)
	router.GET("/api/1/os/:os_name/step/:os_step/:field", xx)
	router.GET("/api/1/ipxe/chain1", xx)
	*/

	http.ListenAndServe(":8080", router)

}
