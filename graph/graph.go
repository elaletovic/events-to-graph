package graph

import (
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	dbName               = "example_db"
	graphName            = "example_graph"
	usersCollection      = "users"
	itemsCollection      = "items"
	viewedEdgeCollection = "viewed"
)

// Store --
type Store struct {
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
func Init(client driver.Client) *Store {
	db, graph := create(client)

	return initStore(db, graph)
}

func create(client driver.Client) (driver.Database, driver.Graph) {
	exists, err := client.DatabaseExists(nil, dbName)
	if err != nil {
		log.Fatalf("failed to check database %s: %v", dbName, err)
	}
	if exists {
		r, err := client.Database(nil, dbName)
		if err != nil {
			log.Fatalf("failed to get database %s: %v", dbName, err)
		}
		r.Remove(nil)
	}

	db := createDB(client)
	graph := createGraph(db)

	return db, graph
}

func getDB(client driver.Client) driver.Database {
	db, err := client.Database(nil, dbName)
	if err != nil {
		log.Fatalf("failed to get a database %s: %v", dbName, err)
	}
	return db
}

func createDB(client driver.Client) driver.Database {
	db, err := client.CreateDatabase(nil, dbName, nil)
	if err != nil {
		log.Fatalf("failed to create a database %s: %v", dbName, err)
	}
	return db
}

func removeGraph(db *driver.Database) {
	graph, err := (*db).Graph(nil, graphName)
	if err != nil {
		log.Fatalf("failed to get a graph %s: %v", graphName, err)
	}
	err = graph.Remove(nil)
	if err != nil {
		log.Fatalf("failed to remove graph %s: %v", graphName, err)
	}
	log.Println("graph removed")
}

func createGraph(db driver.Database) driver.Graph {
	var definition driver.EdgeDefinition
	definition.Collection = viewedEdgeCollection
	definition.From = []string{usersCollection}
	definition.To = []string{itemsCollection}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{definition}
	graph, err := db.CreateGraph(nil, graphName, &options)
	if err != nil {
		log.Fatalf("failed to create a graph %s: %v", graphName, err)
	}
	return graph
}

func initStore(db driver.Database, graph driver.Graph) *Store {
	users := initCollection(graph, usersCollection)
	items := initCollection(graph, itemsCollection)
	usersToItems := initEdgeCollection(graph, viewedEdgeCollection)

	return &Store{
		DB:           db,
		Graph:        graph,
		UsersToItems: usersToItems,
		Users:        users,
		Items:        items,
	}
}

func initCollection(graph driver.Graph, name string) driver.Collection {
	col, err := graph.VertexCollection(nil, name)
	if err != nil {
		log.Fatalf("failed to get vertex %s: %v", name, err)
	}
	return col
}

func initEdgeCollection(graph driver.Graph, name string) edgeCollection {
	col, constraints, err := graph.EdgeCollection(nil, name)
	if err != nil {
		log.Fatalf("failed to get edge %s: %v", name, err)
	}

	return edgeCollection{
		Collection:  col,
		Constraints: constraints,
	}
}
