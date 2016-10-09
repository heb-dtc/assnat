package controller

import (
	"gopkg.in/mgo.v2"
	"os"
	"fmt"
	"os/user"
)

const APPLICATION_DIRECTORY = ".assnat"
const APPLICATION_SCRIPT_DIRECTORY = "script"
const DATABASE_LOCATION =  "mongodb://localhost"

type (
	SetupController struct {
		ApplicationHomeDirectory string
	}
)

func NewSetupController() *SetupController {
	return &SetupController{
		ApplicationHomeDirectory: getApplicationHomeDirectory(),
	}
}

func getApplicationHomeDirectory() string {
	currentUser, userError := user.Current()
	if userError != nil {
		fmt.Println("Can't retrieve current user: ", userError)
		panic(userError)
	}
	return fmt.Sprintf("%s/%s", currentUser.HomeDir, APPLICATION_DIRECTORY)
}

func (setupController SetupController) Init() {
	fmt.Println("Checking application home directory: ", setupController.ApplicationHomeDirectory)
	if _,err := os.Stat(setupController.ApplicationHomeDirectory); os.IsNotExist(err) {
		settingsDirError := os.Mkdir(setupController.ApplicationHomeDirectory, 0777)
		if settingsDirError != nil {
			fmt.Println("Can't create settings directory: ", settingsDirError)
			panic(settingsDirError)
		}
	}

	applicationScriptDirectory := fmt.Sprintf("%s/%s", setupController.ApplicationHomeDirectory, APPLICATION_SCRIPT_DIRECTORY)
	fmt.Println("Checking application script directory: ", applicationScriptDirectory)
	if _,err := os.Stat(applicationScriptDirectory); os.IsNotExist(err) {
		scriptDirError := os.Mkdir(applicationScriptDirectory, 0777)
		if scriptDirError != nil {
			fmt.Println("Can't create script directory: ", scriptDirError)
			panic(scriptDirError)
		}
	}
}

func (setupController SetupController) GetDatabaseConnection() *mgo.Session {
	return getSession()
}

func getSession() *mgo.Session {
	fmt.Println("Retrieving database session")
	session, err := mgo.Dial(DATABASE_LOCATION)

	if err != nil {
		panic(err)
	}
	return session
}

func (setupController SetupController) UpdateLocalData(dataController *DataController) {

	deputyFileLocation := fmt.Sprintf("%s/%s", setupController.ApplicationHomeDirectory, dataController.GetDeputyFileName())
	fmt.Println("Checking that local data is fresh: ", deputyFileLocation)
	if _,err := os.Stat(deputyFileLocation); os.IsNotExist(err) {
		fmt.Println("No data file found, update required")
		dataController.UpdateList(setupController)
	} else {
		fmt.Println("Data file found, no update required")
	}

}