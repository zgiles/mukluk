package main

import (
  "fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/gorilla/context"

  "gomukluk/stores/nodes"
  "gomukluk/stores/nodesredis"
  "gomukluk/stores/nodesmysql"
  "gomukluk/stores/nodesdiscovered"
  "gomukluk/stores/nodesdiscoveredredis"
  "gomukluk/stores/nodesdiscoveredmysql"

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

  var nodestoredb nodes.NodeStoreDB
  var nodestore nodes.NodeStore
  var nodesdiscoveredstoredb nodesdiscovered.NodesDiscoveredStoreDB
  var nodesdiscoveredstore nodesdiscovered.NodesDiscoveredStore

  switch config.Serverconfig.Maindb {
    case "mysql":
      if config.Mysqlconfig.Enabled == false { log.Fatal("mysql selected, but not enabled") }

      log.Println("mysql: opening mysql connection")
      mysqlpool, mysqlerr := mysqlStart(config.Mysqlconfig) // this is a db *sql.DB
      if mysqlerr != nil {
        log.Fatal(mysqlerr)
      }
      defer mysqlpool.Close()
      log.Println("mysql: open")

      log.Println("mysql: opening NodeStoreDB")
      nodestoredb = nodesmysql.NewNodesMysql(mysqlpool)
      log.Println("mysql: opening NodeStore")
      nodestore = nodes.NewNodeStore(nodestoredb)

      log.Println("mysql: opening NodeDiscoveredStoreDB")
      nodesdiscoveredstoredb = nodesdiscoveredmysql.NewNodesDiscoveredMysql(mysqlpool)
      log.Println("mysql: opening NodesDiscoveredStore")
      nodesdiscoveredstore = nodesdiscovered.NewNodesDiscoveredStore(nodesdiscoveredstoredb)

    case "redis":
      if config.Redisconfig.Enabled == false { log.Fatal("redis selected, but not enabled") }

      log.Println("redis: opening redis connection")
      redispool, rediserr := redisStart(config.Redisconfig) // this is a redispool *redis.Pool
      if rediserr != nil {
        log.Fatal(rediserr)
      }
      // defer db.Close()
      log.Println("redis: open")

      log.Println("redis: opening NodeStoreDB")
      nodestoredb = nodesredis.NewNodesRedis(redispool)
      log.Println("redis: opening NodeStore")
      nodestore = nodes.NewNodeStore(nodestoredb)

      log.Println("redis: opening NodeDiscoveredStoreDB")
      nodesdiscoveredstoredb = nodesdiscoveredredis.NewNodesDiscoveredRedis(redispool)
      log.Println("redis: opening NodesDiscoveredStore")
      nodesdiscoveredstore = nodesdiscovered.NewNodesDiscoveredStore(nodesdiscoveredstoredb)

    default:
      log.Fatal("no valid db selected as primary")

  }

	// app context
  appC := appContext{ nodestore: nodestore, nodesdiscoveredstore: nodesdiscoveredstore }
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
