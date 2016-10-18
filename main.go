package main

import (
	"assnat/controller"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
)

func main() {
	router := httprouter.New()

	setupController := controller.NewSetupController()
	setupController.Init()
	session := setupController.GetDatabaseConnection()
	dataController := controller.NewDataController(session)
	setupController.UpdateLocalData(dataController)

	deputyController := controller.NewDeputyController(session)

	router.GET("/deputy/:name", deputyController.GetDeputyByName)
	router.GET("/deputy", deputyController.GetAllDeputies)

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
