package main

import (
  "fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/gorilla/context"

// 	_ "github.com/go-sql-driver/mysql"
//	"database/sql"

  "gomukluk/stores/nodes"
  "gomukluk/stores/nodesredis"
  //"gomukluk/stores/nodesmysql"
  "gomukluk/stores/nodesdiscovered"
  "gomukluk/stores/nodesdiscoveredredis"
  //"gomukluk/stores/nodesdiscoveredmysql"
)


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

/*
func helloHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
*/

/*
type appContext struct {
  db *sql.DB
	config *config
  redispool *redis.Pool
}
*/

type appContext struct {
  nodestore nodes.NodeStore
  nodesdiscoveredstore nodesdiscovered.NodesDiscoveredStore
}

func main() {

  // Closing channel
  // quitting := make(chan bool)

	// Options Parse

	// Config Stage
	config, configerr := loadConfig("config.toml")
	if configerr != nil {
		log.Fatal(configerr)
	}

/*
  // MySQL Stage
  log.Println("mysql: opening mysql connection")
  mysqlpool, mysqlerr := mysqlStart(config.Mysqlconfig) // this is a db *sql.DB
  if mysqlerr != nil {
    log.Fatal(mysqlerr)
  }
  defer mysqlpool.Close()
  log.Println("mysql: open")
*/


  // Redis Stage
  log.Println("redis: opening redis connection")
  redispool, rediserr := redisStart(config.Redisconfig) // this is a // this is a redispool *redis.Pool
  if rediserr != nil {
    log.Fatal(rediserr)
  }
  // defer db.Close()
  log.Println("redis: open")



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
  /*
  // Redis Stage
  log.Println("redis pool opening")
  redispool := newRedisPool(config.Redisconfig.Host, config.Redisconfig.Password) // this is a redispool *redis.Pool
  log.Println("redis open")
  defer log.Println("redis closing")
  defer redispool.Close()
  */

  // make the redis NodeStoreDB
  log.Println("opening redis NodeStoreDB")
  local_nodesredis := nodesredis.NewNodesRedis(redispool)

  // make the mysql NodeStoreDB
  //log.Println("opening mysql NodeStoreDB")
  //local_nodesmysql := nodesmysql.NewNodesMysql(mysqlpool)

  // make the redis NodeStore from all the NodeStoreDBs
  log.Println("opening NodeStore")
  local_nodestore := nodes.NewNodeStore(local_nodesredis)
  // local_nodestore := nodes.NewNodeStore(local_nodesmysql)


  // make the redis NodeDiscoveredStoreDB
  log.Println("opening redis NodeDiscoveredStoreDB")
  local_nodesdiscoveredredis := nodesdiscoveredredis.NewNodesDiscoveredRedis(redispool)

  // make the mysql NodeDiscoveredStoreDB
  // log.Println("opening mysql NodeDiscoveredStoreDB")
  // local_nodesdiscoveredmysql := nodesdiscoveredmysql.NewNodesDiscoveredMysql(mysqlpool)

  // make the redis NodesDiscoveredStore from all the NodesDiscoveredStoreDB's
  log.Println("opening NodesDiscoveredStore")
  local_nodesdiscoveredstore := nodesdiscovered.NewNodesDiscoveredStore(local_nodesdiscoveredredis)
  // local_nodesdiscoveredstore := nodesdiscovered.NewNodesDiscoveredStore(local_nodesdiscoveredmysql)


	// app context
  //db: db,
	//appC := appContext{ config: config, redispool: redispool }
  appC := appContext{ nodestore: local_nodestore, nodesdiscoveredstore: local_nodesdiscoveredstore }
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
