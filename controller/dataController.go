package controller

import (
	"gopkg.in/mgo.v2"
	"os"
	"fmt"
	"net/http"
	"io"
	"archive/zip"
	"path/filepath"
	"os/exec"
	"strings"
	"assnat/model"
	"gopkg.in/mgo.v2/bson"
)

const DEPUTY_LIST_URL = "http://data.assemblee-nationale.fr/static/openData/repository/AMO/deputes_actifs_mandats_actifs_organes/AMO10_deputes_actifs_mandats_actifs_organes_XIV.json.zip"
const ZIP_LIST_NAME = "deputy_list.zip"
const PARSE_SCRIPT_NAME = "/script/parse.py"
const DEPUTY_LIST_NAME = "list.json"

type (
	DataController struct {
		session *mgo.Session
	}
)

func NewDataController(session *mgo.Session) *DataController {
	return &DataController{session}
}

func (dataController DataController) GetDeputyFileName() string {
	return DEPUTY_LIST_NAME
}

func (dataController DataController) UpdateList(setupController SetupController) {
	fmt.Println("Updating deputy list")

	listLocation := fmt.Sprintf("%s/%s", setupController.ApplicationHomeDirectory, ZIP_LIST_NAME)
	fetchList(listLocation)
	unPackageData(listLocation, setupController.ApplicationHomeDirectory)
	cleanUpDownloadArtifact(listLocation)

	pathToScript := fmt.Sprintf("%s/%s", setupController.ApplicationHomeDirectory, PARSE_SCRIPT_NAME)
	pathToList := fmt.Sprintf("%s/%s", setupController.ApplicationHomeDirectory, DEPUTY_LIST_NAME)
	parseData(dataController.session, pathToScript, pathToList)
}

func cleanUpDownloadArtifact(listLocation string) {
	err := os.Remove(listLocation)
	if err != nil {
		fmt.Println("can't remove donwload artifact: ", err)
	}
}

func fetchList(listLocation string) {
	fmt.Println("Fetch deputy list")
	out, fileCreateError := os.Create(listLocation)
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

func unPackageData(listLocation string, applicationHomeDir string) error {
	fmt.Println("Unpackaging data")

	reader, err := zip.OpenReader(listLocation)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(applicationHomeDir, DEPUTY_LIST_NAME)
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

func parseData(session *mgo.Session, pathToScript string, pathToList string) {
	fmt.Println("parsing data")
	cmd := exec.Command("python3", pathToScript, "members", pathToList)
	out, err := cmd.Output()

	if err != nil {
		println("error while parsing:", err.Error())
		return
	}

	data := string(out)
	result := strings.Split(data, "\n")
	session.DB("ass_nat").C("deputy").DropCollection()

	for i := range result {
		if i > 0 && i < len(result) {
			fmt.Println(result[i])
			writeEntryToDatabase(session, result[i])
		}
	}
}
func writeEntryToDatabase(session *mgo.Session, entry string) {
	result := strings.Split(entry, ";")

	if len(result) == 5 {
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
}
