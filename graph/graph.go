package graph

import (
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	dbName                   = "example_db"
	graphName                = "example_graph"
	usersCollection          = "users"
	itemsCollection          = "items"
	usersItemsEdgeCollection = "users_items"
)

type store struct {
	DB           driver.Database
	Graph        driver.Graph
	Users        driver.Collection
	Items        driver.Collection
	UsersToItems edgeCollection
}

type edgeCollection struct {
	Collection  driver.Collection
	Constraints driver.VertexConstraints
}

var s = store{}

// Connect connects to the graph server
func Connect(address string) driver.Connection {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{address},
	})
	if err != nil {
		log.Fatalf("failed to create HTTP connection: %v", err)
	}
	return conn
}

// GetClient creates a graph client
func GetClient(conn driver.Connection) driver.Client {
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "rootpassword"),
	})
	if err != nil {
		log.Fatalf("failed to create graph client: %v", err)
	}

	return c
}

// Init inits the graph database and graph
func Init(client driver.Client) {
	exists, err := client.DatabaseExists(nil, dbName)
	if err != nil {
		log.Fatalf("failed to check for a database %s: %v", dbName, err)
	}
	if exists {
		get(client)
	} else {
		create(client)
	}
	initCollections()
}

func get(client driver.Client) {
	getDB(client)
	graphExists, err := s.DB.GraphExists(nil, graphName)
	if err != nil {
		log.Fatalf("failed to check if graph %s exists: %v", graphName, err)
	}
	if graphExists {
		getGraph()
	} else {
		createGraph()
	}
}

func create(client driver.Client) {
	createDB(client)
	createGraph()
}

func getDB(client driver.Client) {
	db, err := client.Database(nil, dbName)
	if err != nil {
		log.Fatalf("failed to get a database %s: %v", dbName, err)
	}
	s.DB = db
}

func createDB(client driver.Client) {
	db, err := client.CreateDatabase(nil, dbName, nil)
	if err != nil {
		log.Fatalf("failed to create a database %s: %v", dbName, err)
	}
	s.DB = db
}

func getGraph() {
	graph, err := s.DB.Graph(nil, graphName)
	if err != nil {
		log.Fatalf("failed to get a graph %s: %v", graphName, err)
	}
	s.Graph = graph
}

func createGraph() {
	var definition driver.EdgeDefinition
	definition.Collection = usersItemsEdgeCollection
	definition.From = []string{usersCollection}
	definition.To = []string{itemsCollection}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{definition}
	graph, err := s.DB.CreateGraph(nil, graphName, &options)
	if err != nil {
		log.Fatalf("failed to create a graph %s: %v", graphName, err)
	}
	s.Graph = graph

	log.Println(s.Graph.Name())
}

func initCollections() {
	s.Users = initCollection(usersCollection)
	s.Items = initCollection(itemsCollection)
	s.UsersToItems = initEdgeCollection(usersItemsEdgeCollection)
}

func initCollection(name string) driver.Collection {
	col, err := s.Graph.VertexCollection(nil, name)
	if err != nil {
		log.Fatalf("failed to get vertex %s: %v", name, err)
	}
	return col
}

func initEdgeCollection(name string) edgeCollection {
	col, constraints, err := s.Graph.EdgeCollection(nil, name)
	if err != nil {
		log.Fatalf("failed to get edge %s: %v", name, err)
	}

	return edgeCollection{
		Collection:  col,
		Constraints: constraints,
	}
}
