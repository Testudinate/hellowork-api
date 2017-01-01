package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/goengine"
	"github.com/hellofresh/goengine/inmemory"
	"github.com/hellofresh/goengine/mongodb"
	"github.com/italolelis/hellowork-api/config"
	"github.com/italolelis/hellowork-api/errors"
	"github.com/italolelis/hellowork-api/middleware"
	"github.com/italolelis/hellowork-api/response"
	"github.com/italolelis/hellowork-api/router"
	"github.com/italolelis/hellowork-api/status"
	"github.com/italolelis/hellowork-api/user"
	"gopkg.in/mgo.v2"
)

var (
	err          error
	globalConfig *config.Specification
)

// initializes the global configuration
func init() {
	globalConfig, err = config.LoadEnv()
	if nil != err {
		log.Panic(err.Error())
	}
}

// initializes the basic configuration for the log wrapper
func init() {
	level, err := log.ParseLevel(strings.ToLower(globalConfig.LogLevel))
	if err != nil {
		log.Error("Error getting level", err)
	}

	log.SetLevel(level)
}

func main() {
	log.Infof("Connecting to write database %s", globalConfig.Database.WriteDSN)
	writeSession, err := mgo.Dial(globalConfig.Database.WriteDSN)
	if err != nil {
		log.Panic(err)
	}
	defer writeSession.Close()

	log.Infof("Connecting to read database %s", globalConfig.Database.ReadDSN)
	readSession, err := mgo.Dial(globalConfig.Database.ReadDSN)
	if err != nil {
		log.Panic(err)
	}
	defer readSession.Close()

	// Optional. Switch the session to a monotonic behavior.
	writeSession.SetMode(mgo.Monotonic, true)
	readSession.SetMode(mgo.Monotonic, true)

	log.Info("Setting up the registry")
	registry := goengine.NewInMemmoryTypeRegistry()
	registry.RegisterType(&user.UserCreated{})

	log.Info("Setting up the event bus")
	bus := inmemory.NewInMemoryEventBus()

	log.Info("Setting up the event store")
	es := mongodb.NewEventStore(writeSession, registry)

	publisherRepo := goengine.NewPublisherRepository(es, bus)

	statusWriteRepo, err := status.NewMongoDBWriteRepository(publisherRepo)
	if err != nil {
		log.Panic(err)
	}

	statusReadRepo, err := status.NewMongoDBReadRepository(readSession.DB(""))
	if err != nil {
		log.Panic(err)
	}

	eventDispatcher := goengine.NewVersionedEventDispatchManager(bus, registry)

	initStatusesProjection(readSession, eventDispatcher)
	stopChannel := make(chan bool)
	go eventDispatcher.Listen(stopChannel, false)

	r := router.NewHttpTreeMuxRouter()
	r.Use(middleware.NewRecovery(RecoveryHandler).Handler)

	r.GET("/", Home())
	loadUsersEndpoints(r, publisherRepo)
	loadStatusesEndpoints(r, statusWriteRepo, statusReadRepo)

	log.Debugf("Listening on :%v", globalConfig.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.Port), r))
}

func loadUsersEndpoints(router router.Router, repo goengine.AggregateRepository) {
	log.Debug("Loading Users Endpoints")

	usersController := user.NewUsersController(repo)
	group := router.Group("/users")
	{
		group.GET("", usersController.Get)
		group.POST("", usersController.Post)
		group.GET("/:id", usersController.GetBy)
		group.PUT("/:id", usersController.PutBy)
	}
}

func loadStatusesEndpoints(router router.Router, write status.WriteRepository, read status.ReadStatusRepository) {
	log.Debug("Loading Statuses Endpoints")

	c := status.NewController(write, read)
	group := router.Group("/statuses")
	{
		group.GET("", c.Get)
		group.POST("", c.Post)
		group.GET("/:id", c.GetBy)
		group.PUT("/:id", c.PutBy)
		group.DELETE("/:id", c.DeleteBy)
	}
}

func initStatusesProjection(session *mgo.Session, dispatcher *goengine.VersionedEventDispatchManager) {
	log.Debug("Loading statuses projection listeners")
	p := status.NewProjection(session.DB(""))

	dispatcher.RegisterEventHandler(&status.StatusCreated{}, p.StatusCreated)
	dispatcher.RegisterEventHandler(&status.StatusRemoved{}, p.StatusRemoved)
	dispatcher.RegisterEventHandler(&status.StatusReasonChanged{}, p.StatusReasonChanged)
	dispatcher.RegisterEventHandler(&status.StatusTime{}, p.StatusTime)
	dispatcher.RegisterEventHandler(&status.AllDay{}, p.StatusAllDay)
	dispatcher.RegisterEventHandler(&status.StatusTimePeriodChanged{}, p.StatusTimePeriodChanged)
}

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, fmt.Sprint("Welcome to the hellowork api"))
	}
}

// RecoveryHandler handler for the apis
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch internalErr := err.(type) {
	case *errors.Error:
		log.Error(internalErr.Error())
		response.JSON(w, internalErr.Code, internalErr.Error())
	default:
		response.JSON(w, http.StatusInternalServerError, err)
	}
}
