package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/italolelis/hellowork-api/config"
	"github.com/italolelis/hellowork-api/errors"
	"github.com/italolelis/hellowork-api/middleware"
	"github.com/italolelis/hellowork-api/response"
	"github.com/italolelis/hellowork-api/router"
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
	r := router.NewHttpTreeMuxRouter()
	r.Use(middleware.NewRecovery(RecoveryHandler).Handler)

	r.GET("/", Home())

	log.Debugf("Listening on :%v", globalConfig.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.Port), r))
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
