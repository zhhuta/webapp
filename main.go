package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"

	"golang.org/x/net/context"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	PubsubClient  *pubsub.Client //  global
)

// topic to publich
const PubsubTopicID = "one" // read topic

type templateParams struct {
	Date    string
	Time    string
	Notice  string
	Warning string
	Name    string
}
type eventMessage struct {
	name    string
	message string
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	params := templateParams{}
	currentDate := time.Now().Local()
	params.Date = currentDate.Format("2006-02-01")
	params.Time = currentDate.Format("3:04 PM")

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
	}

	if r.Method == "POST" {

		name := r.FormValue("name")
		params.Name = name
		if name == "" {
			name = "Anonymous"
		}

		message := r.FormValue("message")
		if r.FormValue("message") == "" {
			w.WriteHeader(http.StatusBadRequest)

			params.Warning = "No message"
			indexTemplate.Execute(w, params)
			return
		}
		go publishUpdate(eventMessage{name: name, message: message})
		params.Notice = fmt.Sprintf("Message from %s: %s", name, message)
		indexTemplate.Execute(w, params)

	}
}

func publishUpdate(event eventMessage) {
	ctx := context.Background()

	PubsubClient, err := configurePubsub("riverlife-197216")
	if err != nil {
		log.Fatalf("Issue during configuring pubsub")
	}

	topic := PubsubClient.Topic(PubsubTopicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatal("Error cheking for topic %v", err)
	}
	if !exists {
		if _, err := PubsubClient.CreateTopic(ctx, PubsubTopicID); err != nil {
			log.Fatal("Failed to create Topic: %v", err)
		}
	}

	b, err := json.Marshal(event)
	if err != nil {
		return
	}

	_, err = topic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	log.Printf("Published update to Pub/Sub for Event ID %d: %v", event.name, err)
}
func configurePubsub(projectID string) (*pubsub.Client, error) {
	//For beginign we have to configure PubSub.Clinet base on our PROJECT_ID
	ctx := context.Background()
	//Creating a new clinent
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	//Create topic if it's not exit.
	if exists, err := client.Topic(PubsubTopicID).Exists(ctx); err != nil {
		return nil, err
	} else if !exists {
		if _, err := client.CreateTopic(ctx, PubsubTopicID); err != nil {
			return nil, err
		}
	}
	return client, err
}
