package main

import (
	"log"
	"strconv"
	"time"

	"github.com/tylerb/graceful" // "gopkg.in/tylerb/graceful.v1"
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/zgiles/mukluk/stores/nodes"
	"github.com/zgiles/mukluk/stores/nodesdiscovered"
	"github.com/zgiles/mukluk/stores/nodesdiscoveredmysql"
	"github.com/zgiles/mukluk/stores/nodesdiscoveredredis"
	"github.com/zgiles/mukluk/stores/nodesmysql"
	"github.com/zgiles/mukluk/stores/nodesredis"
	"github.com/zgiles/mukluk/stores/oses"
	"github.com/zgiles/mukluk/stores/osesmysql"
)

type appContext struct {
	nodestore            	nodes.NodeStore
	nodesdiscoveredstore 	nodesdiscovered.NodesDiscoveredStore
	osstore             	oses.OsStore
	ipxeconfig            ipxeconfig
}

func main() {
	// Closing channel
	// Currently handled by graceful, but just in case.
	// quitting := make(chan os.Signal)
	// signal.Notify(quitting, syscall.SIGINT, syscall.SIGTERM)

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
	var osstoredb oses.OsStoreDB
	var osstore oses.OsStore

	switch config.Serverconfig.Maindb {
	case "mysql":
		if config.Mysqlconfig.Enabled == false {
			log.Fatal("mysql selected, but not enabled")
		}

		log.Println("mysql: opening mysql connection")
		mysqlpool, mysqlerr := mysqlStart(config.Mysqlconfig) // this is a db *sql.DB
		if mysqlerr != nil {
			log.Fatal(mysqlerr)
		}
		defer log.Println("mysql: closing")
		defer mysqlpool.Close()
		log.Println("mysql: open")

		log.Println("mysql: opening NodeStoreDB")
		nodestoredb = nodesmysql.NewNodesMysql(mysqlpool)
		log.Println("mysql: opening NodeStore")
		nodestore = nodes.NewNodeStore(nodestoredb)

		log.Println("mysql: opening NodesDiscoveredStoreDB")
		nodesdiscoveredstoredb = nodesdiscoveredmysql.NewNodesDiscoveredMysql(mysqlpool)
		log.Println("mysql: opening NodesDiscoveredStore")
		nodesdiscoveredstore = nodesdiscovered.NewNodesDiscoveredStore(nodesdiscoveredstoredb)

		log.Println("mysql: opening OsStoreDB")
		osstoredb = osesmysql.NewOsesMysql(mysqlpool)
		log.Println("mysql: opening OsStore")
		osstore = oses.NewOsStore(osstoredb)

	case "redis":
		if config.Redisconfig.Enabled == false {
			log.Fatal("redis selected, but not enabled")
		}

		log.Println("redis: opening redis connection")
		redispool, rediserr := redisStart(config.Redisconfig) // this is a redispool *redis.Pool
		if rediserr != nil {
			log.Fatal(rediserr)
		}
		defer log.Println("redis: no close needed...")
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
	appC := appContext{ nodestore: nodestore, nodesdiscoveredstore: nodesdiscoveredstore, osstore: osstore, ipxeconfig: config.Ipxeconfig }
	log.Println("app ready")

	// common routes
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)

	// routers
	router := httprouter.New()
	router.GET("/", wrapHandler(commonHandlers.ThenFunc(indexHandler)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.httpGetNodeByFieldHandler)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetNodeByFieldHandler)))
	router.GET("/api/1/nodes/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.httpGetNodesByFieldHandler)))
	router.GET("/api/1/discoverednode/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.httpGetDiscoveredNodeByFieldHandler)))
	router.GET("/api/1/discoverednode/:nodekey/:nodekeyvalue/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetDiscoveredNodeByFieldHandler)))
	router.GET("/api/1/discoverednodes/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.httpGetDiscoveredNodesByFieldHandler)))
	router.GET("/api/1/me/node", wrapHandler(commonHandlers.ThenFunc(appC.httpGetNodeByMyIP)))
	router.GET("/api/1/me/node/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetNodeByMyIP)))
	router.GET("/api/1/me/discoverednode", wrapHandler(commonHandlers.ThenFunc(appC.httpGetDiscoveredNodeByMyIP)))
	router.GET("/api/1/me/discoverednode/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetDiscoveredNodeByMyIP)))
	router.GET("/api/1/os/:os_name/step/:os_step", wrapHandler(commonHandlers.ThenFunc(appC.httpGetOsByNameAndStepHandler)))
	router.GET("/api/1/os/:os_name/step/:os_step/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetOsByNameAndStepHandler)))
	router.GET("/api/1/me/os", wrapHandler(commonHandlers.ThenFunc(appC.httpOsNodeByMyIP)))
	router.GET("/api/1/me/os/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpOsNodeByMyIP)))
	router.GET("/api/1/ipxe/chain1", wrapHandler(commonHandlers.ThenFunc(appC.httpipxechain)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/ipxe", wrapHandler(commonHandlers.ThenFunc(appC.httpipxeNode)))
	router.GET("/api/1/discover/uuid/:uuid/ipv4address/:ipv4address/macaddress/:macaddress", wrapHandler(commonHandlers.ThenFunc(appC.httpipxediscover)))

  router.NotFound = commonHandlers.ThenFunc(errorHandler)

	httpsrv := &graceful.Server{
		Timeout: time.Duration(config.Serverconfig.Closetimeout) * time.Second,
		Server: &http.Server{
			Addr:    config.Serverconfig.Ip + ":" + strconv.FormatInt(config.Serverconfig.Port, 10),
			Handler: router,
		},
	}
	httperr := httpsrv.ListenAndServe()
	if httperr != nil {
		log.Fatal(httperr)
	}

	log.Println("main: end of main")
}
