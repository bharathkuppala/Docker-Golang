package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
)

var (
	logger        *log.Logger
	port          *string
	mgoDatabase   *mgo.Database
	mgoSession    *mgo.Session
	username      = ""
	password      = ""
	database      = "dockerTest"
	collection    = "visitCount"
	connectionURL = "mongodb://localhost:27017"
)

// Visit ...
type Visit struct {
	VisitCount int
	isVisited  bool
}

func main() {
	var documentResult map[string]interface{}
	logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	// port = flag.String("port", "", "port to listen on")

	// flag.Parse()
	// if *port == "" {
	// 	log.Fatal("-port is required")
	// }

	// info := &mgo.DialInfo{
	// 	Addrs:    []string{connectionURL},
	// 	Username: username,
	// 	Password: password,
	// 	Database: database,
	// }

	// mgoSession, err := mgo.DialWithInfo(info)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	mgoSession, err := mgo.Dial("mongo:27017")
	if err != nil {
		logger.Fatalln(err.Error())
	}

	defer mgoSession.Close()

	mgoSession.SetMode(mgo.Monotonic, true)
	mgoSession.SetMode(mgo.Strong, true)
	mgoSession.SetSafe(&mgo.Safe{
		WMode: "majority",
	})

	mgoDatabase = mgoSession.DB(database)
	router := mux.NewRouter()
	var isVisited = false
	var vCount = 0
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Started go server using docker container."))
		mgoDatabaseToCollection := mgoDatabase.C(collection)
		var visit Visit
		if isVisited == false {
			visit.VisitCount = 0
			mgoDatabaseToCollection.Insert(visit)
			isVisited = true

			err = mgoDatabaseToCollection.Find(bson.M{
				"visitcount": 0,
			}).Limit(1).One(&documentResult)
		}

		vCount++

		err = mgoDatabaseToCollection.UpdateId(documentResult["_id"], bson.M{
			"$set": bson.M{
				"visitcount": vCount,
			},
		})

		if err != nil {
			log.Println(err.Error())
		}

		var visits []Visit
		iter := mgoDatabaseToCollection.Find(nil).Limit(50).Iter()
		err = iter.All(&visits)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		b, err := json.Marshal(visits)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		w.Write(b)
		w.Write([]byte("Your Site Visit Count is:" + strconv.FormatInt(int64(visits[0].VisitCount), 10)))
	})

	server := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	// flag.VisitAll(func(flag *flag.Flag) {
	// 	log.Println(flag.Name, "->", flag.Value)
	// })

	if err := server.ListenAndServe(); err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Go-server Started on port 8000 .....")
}
