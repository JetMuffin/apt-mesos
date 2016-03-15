package server

import (
    "net/http"

	"github.com/go-martini/martini"
    "github.com/martini-contrib/cors"
	"github.com/icsnju/apt-mesos/api"
	"github.com/icsnju/apt-mesos/registry"
	"github.com/icsnju/apt-mesos/core"
)

func logger() martini.Handler {
    return func(res http.ResponseWriter, req *http.Request, ctx martini.Context) {
        ctx.Next()
    }
}

func recovery() martini.Handler {
	return func(w http.ResponseWriter, ctx martini.Context) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}

func createRouter(core *core.Core, apis *api.API) martini.Router {
	router := martini.NewRouter()

	// create user endpoints
    router.Get("/api/handshake", apis.Handshake())
    router.Get("/api/tasks", apis.ListTasks())
    router.Post("/api/tasks", apis.AddTask())
    router.Delete("/api/tasks/:id", apis.DeleteTask())  

    router.Get("/api/jobs", apis.ListJobs())
    router.Post("/api/jobs", apis.AddJob())
    router.Delete("/api/jobs/:id", apis.DeleteJob())  
    router.Get("/api/job/:id", apis.GetJob())  

    // create monitor endpoints
    router.Get("/api/system/metrics", apis.SystemMetrics())
    router.Get("/api/slave/metrics", apis.SlaveMetrics())

    // create mesos endpoints
    for method, routes := range core.Endpoints {
    	for route, function := range routes {
    		switch method {
    			case "POST":
    				router.Post(route, function)
    			case "GET":
    				router.Get(route, function)
    			case "DELETE":
    				router.Delete(route, function)
    			case "PUT":
    				router.Put(route, function)
    		}
    	}
    }

    return router
}

func ListenAndServe(addr string, registry *registry.Registry, core *core.Core) {
	apis := api.NewAPI(core, registry)
	r := createRouter(core, apis)

	m := martini.New()
    m.Use(cors.Allow(&cors.Options{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "x-requested-with", "Content-Type", "Content-Range", "Content-Disposition", "Content-Description"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: false,
    }))    
    m.Use(logger())
    m.Use(recovery())
    m.Use(martini.Static("static"))
	m.Action(r.Handle)
    go m.RunOnAddr(addr)
}
