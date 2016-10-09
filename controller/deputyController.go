package controller

import (
    "net/http"
    "os"
    "io"
    "archive/zip"
    "path/filepath"
    "fmt"
    "os/exec"
    "strings"
    "gopkg.in/mgo.v2"
    "assnat/model"
    "gopkg.in/mgo.v2/bson"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
)

const DEPUTY_LIST_URL = "http://data.assemblee-nationale.fr/static/openData/repository/AMO/deputes_actifs_mandats_actifs_organes/AMO10_deputes_actifs_mandats_actifs_organes_XIV.json.zip"
const LIST_LOCATION = "/home/flo/.assnat/deputy_list.zip"
const LOCATION = "/home/flo/.assnat/"
const PARSE_SCRIPT = "/home/flo/.assnat/scripts/parse.py"

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

func (deputyController DeputyController) UpdateList() {
    fmt.Println("Updating deputy list")

    fetchList()
    unPackageData()
    cleanUpDownloadArtifact()
    parseData(deputyController.session)
}

func cleanUpDownloadArtifact() {
    err := os.Remove(LIST_LOCATION)
    if err != nil {
        fmt.Println("can't remove donwload artifact: ", err)
    }
}

func fetchList() {
    fmt.Println("Fetch deputy list")
    out, fileCreateError := os.Create(LIST_LOCATION)
    defer out.Close()

    if fileCreateError != nil {
        fmt.Println("can't create destination file: ", fileCreateError)
    }

    response, fetchError := http.Get(DEPUTY_LIST_URL)
    defer response.Body.Close()

    if fetchError != nil {
        fmt.Println("can't fetch deputy list: ", fetchError)
    }

    _, writeError := io.Copy(out, response.Body)
    if writeError != nil {
        fmt.Println("can't write list: ", writeError)
    }
}

func unPackageData() error {
    fmt.Println("Unpackaging data")

    reader, err := zip.OpenReader(LIST_LOCATION)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(LOCATION, 0755); err != nil {
        return err
    }

    for _, file := range reader.File {
        path := filepath.Join(LOCATION, "list.json") //TODO: change file name?
        if file.FileInfo().IsDir() {
            os.MkdirAll(path, file.Mode())
            continue
        }

        fileReader, err := file.Open()
        if err != nil {
            return err
        }
        defer fileReader.Close()

        targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
        if err != nil {
            return err
        }
        defer targetFile.Close()

        if _, err := io.Copy(targetFile, fileReader); err != nil {
            return err
        }
    }

    return nil
}

func parseData(session *mgo.Session) {
    fmt.Println("parsing data")
    cmd := exec.Command("python", PARSE_SCRIPT, "members", "/home/flo/.assnat/list.json")
    out, err := cmd.Output()

    if err != nil {
        println("error while parsing:", err.Error())
        return
    }

    data := string(out)
    result := strings.Split(data, "\n")

    for i := range result {
        if i > 0 && i < len(result) {
            fmt.Println(result[i])
            writeEntryToDatabase(session, result[i])
        }
    }
}
func writeEntryToDatabase(session *mgo.Session, entry string) {
    result := strings.Split(entry, " ")

    deputy := model.Deputy{
        Id: bson.NewObjectId(),
        Title: result[1],
        FirstName: result[2],
        LastName: result[3],
        Province: result[4],
    }

    fmt.Println(deputy)
    session.DB("ass_nat").C("deputy").Insert(deputy)
}
