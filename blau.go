package main

import (
	"log"
	"os"

	"runtime/debug"
	"runtime/trace"
	"time"

	"fmt"
	"net/http"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/blauentag/blau/br"
	_ "github.com/blauentag/blau/controllerRegistry"
	"github.com/blauentag/blau/controllers"
	"github.com/blauentag/blau/cron"
	_ "github.com/blauentag/blau/gitWebHooks"
	_ "github.com/blauentag/blau/models/v1/model"
	"github.com/blauentag/blau/sessionFunctions"
	"github.com/blauentag/blau/settings"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("studio.go", "Panic Recovered at main():"+fmt.Sprintf("%+v", r))
			time.Sleep(time.Millisecond * 3000)
			main()
			return
		}
	}()

	err := app.Initialize("src/github.com/blauentag/blau", "webConfig.json")
	settings.Initialize()
	br.Schedules.UpdateLinuxToGMT()

	if err != nil {
		//lastError := err.Error()
		ginServer.Router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "%v", "An error occurred and the blau app cannot run (most likely due to mongo database services being down).\n\nError description: "+err.Error())
		})
		app.Run()
	} else {
		if settings.AppSettings.DeveloperGoTrace {
			f, err := os.Create("src/github.com/blauentag/blau/log/trace.log")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			err = trace.Start(f)
			if err != nil {
				panic(err)
			}
			mgo.SetDebug(true)

			file, _ := os.OpenFile("src/github.com/blauentag/blau/log/studioMongo.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

			var aLogger *log.Logger
			aLogger = log.New(file, "", log.LstdFlags)

			mgo.SetLogger(aLogger)
			mgo.SetStats(true)
		}

		controllers.Initialize()

		core.CronJobs.Start()
		cron.Start()

		go logger.GoRoutineLogger(func() {
			time.Sleep(time.Millisecond * 5000)
			br.Schedules.LoadDay(time.Now())
		}, "main->Loading Schedules")

		app.Run()
	}
}
