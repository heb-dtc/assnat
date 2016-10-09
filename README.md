### AssNat

1. setup environement
> install golang sdk
  create a directory to host the `go workspace` (i.e. ~/go_workspace)  
  define `GOPATH̀` environnement variable pointing to previously defined workpace  
  create a `src/` directory inside the worksapce  
  clone the project there  
  
2. install the application
> install all the dependencies and the application  
  `$ go get gopkg.in/mgo.v2`  
  `$ go get github.com/julienschmidt/httprouter`  
  `$ go get gopkg.in/mgo.v2/bson`  
  `$ go install`   
