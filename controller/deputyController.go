package controller

import (
    "net/http"
    "fmt"
    "gopkg.in/mgo.v2"
    "assnat/model"
    "gopkg.in/mgo.v2/bson"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
)

type (
    DeputyController struct {
        session *mgo.Session
    }
)

func NewDeputyController(session *mgo.Session) *DeputyController {
    return &DeputyController{session}
}

func (deputyController DeputyController) GetDeputyByName(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    name := p.ByName("name")

    query := bson.M{
        "lastname": name,
    }
    deputy := model.Deputy {}

    err := deputyController.session.DB("ass_nat").C("deputy").Find(query).One(&deputy);
    if err != nil {
        fmt.Println("Ooops: ", err)
        w.WriteHeader(404)
        return
    }

    fmt.Println(deputy)
    deputyJson, _ := json.Marshal(deputy)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", deputyJson)
}
