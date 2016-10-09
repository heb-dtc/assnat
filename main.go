package main

import (
    "assnat/controller"
    "net/http"
    "gopkg.in/mgo.v2"
    "github.com/julienschmidt/httprouter"
    "log"
)

func getSession() *mgo.Session {
    // Connect to our local mongo
    s, err := mgo.Dial("mongodb://localhost")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    return s
}


func main() {
    
    router := httprouter.New()
    
    deputyController := controller.NewDeputyController(getSession())
    //deputyController.UpdateList()
    
    router.GET("/deputy/:name", deputyController.GetDeputyByName)

    err := http.ListenAndServe(":3000", router)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}