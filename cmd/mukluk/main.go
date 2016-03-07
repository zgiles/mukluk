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

	"github.com/zgiles/mukluk/stores/nodestore"
	"github.com/zgiles/mukluk/stores/osstore"
	"github.com/zgiles/mukluk/stores/nodesdiscoveredstore"

	mysql_nodesdiscovereddb "github.com/zgiles/mukluk/db/mysql/nodesdiscovereddb"
	redis_nodesdiscovereddb "github.com/zgiles/mukluk/db/redis/nodesdiscovereddb"
	mysql_nodesdb "github.com/zgiles/mukluk/db/mysql/nodesdb"
	redis_nodesdb "github.com/zgiles/mukluk/db/redis/nodesdb"
	mysql_osdb "github.com/zgiles/mukluk/db/mysql/osdb"
	redis_osdb "github.com/zgiles/mukluk/db/redis/osdb"
)

type appContext struct {
	nodestore            	nodestore.StoreI
	nodesdiscoveredstore 	nodesdiscoveredstore.StoreI
	osstore             	osstore.StoreI
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

	var l_nodestoredb nodestore.StoreDBI
	var l_nodestore nodestore.StoreI
	var l_nodesdiscoveredstoredb nodesdiscoveredstore.StoreDBI
	var l_nodesdiscoveredstore nodesdiscoveredstore.StoreI
	var l_osstoredb osstore.StoreDBI
	var l_osstore osstore.StoreI

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
		l_nodestoredb = mysql_nodesdb.New(mysqlpool)
		log.Println("mysql: opening NodeStore")
		l_nodestore = nodestore.New(l_nodestoredb)

		log.Println("mysql: opening NodesDiscoveredStoreDB")
		l_nodesdiscoveredstoredb = mysql_nodesdiscovereddb.New(mysqlpool)
		log.Println("mysql: opening NodesDiscoveredStore")
		l_nodesdiscoveredstore = nodesdiscoveredstore.New(l_nodesdiscoveredstoredb)

		log.Println("mysql: opening OsStoreDB")
		l_osstoredb = mysql_osdb.New(mysqlpool)
		log.Println("mysql: opening OsStore")
		l_osstore = osstore.New(l_osstoredb)

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
		l_nodestoredb = redis_nodesdb.New(redispool)
		log.Println("redis: opening NodeStore")
		l_nodestore = nodestore.New(l_nodestoredb)

		log.Println("redis: opening NodeDiscoveredStoreDB")
		l_nodesdiscoveredstoredb = redis_nodesdiscovereddb.New(redispool)
		log.Println("redis: opening NodesDiscoveredStore")
		l_nodesdiscoveredstore = nodesdiscoveredstore.New(l_nodesdiscoveredstoredb)

		log.Println("redis: opening OsStoreDB")
		l_osstoredb = redis_osdb.New(redispool)
		log.Println("redis: opening OsStore")
		l_osstore = osstore.New(l_osstoredb)

	default:
		log.Fatal("no valid db selected as primary")

	}

	// app context
	appC := appContext{ nodestore: l_nodestore,
		nodesdiscoveredstore: l_nodesdiscoveredstore,
		osstore: l_osstore,
		ipxeconfig: config.Ipxeconfig }
	log.Println("app ready")

	// common routes
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)

	// routers
	router := httprouter.New()
	router.GET("/", wrapHandler(commonHandlers.ThenFunc(indexHandler)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.getnode)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.getnode)))
	router.GET("/api/1/node/:nodekey/:nodekeyvalue/ipxe", wrapHandler(commonHandlers.ThenFunc(appC.httpipxeNode)))

	router.GET("/api/1/nodes/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.getnodes)))

	router.GET("/api/1/discover/uuid/:uuid/ipv4address/:ipv4address/macaddress/:macaddress", wrapHandler(commonHandlers.ThenFunc(appC.httpipxediscover)))

	router.GET("/api/1/discoverednode/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.getnodediscovered)))
	router.GET("/api/1/discoverednode/:nodekey/:nodekeyvalue/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.getnodediscovered)))

	router.GET("/api/1/discoverednodes/:nodekey/:nodekeyvalue", wrapHandler(commonHandlers.ThenFunc(appC.getnodediscovereds)))

	router.GET("/api/1/me/node", wrapHandler(commonHandlers.ThenFunc(appC.getnodebyip)))
	router.GET("/api/1/me/node/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.getnodebyip)))
	router.GET("/api/1/me/discoverednode", wrapHandler(commonHandlers.ThenFunc(appC.getnodediscoveredbyip)))
	router.GET("/api/1/me/discoverednode/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.getnodediscoveredbyip)))
	router.GET("/api/1/me/os", wrapHandler(commonHandlers.ThenFunc(appC.httpOsNodeByMyIP)))
	router.GET("/api/1/me/os/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpOsNodeByMyIP)))

	router.GET("/api/1/os/:os_name/step/:os_step", wrapHandler(commonHandlers.ThenFunc(appC.httpGetOsByNameAndStepHandler)))
	router.GET("/api/1/os/:os_name/step/:os_step/field/:field", wrapHandler(commonHandlers.ThenFunc(appC.httpGetOsByNameAndStepHandler)))

	router.GET("/api/1/ipxe/chain1", wrapHandler(commonHandlers.ThenFunc(appC.httpipxechain)))

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
