package graph

import (
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	dbName                     = "example_db"
	graphName                  = "example_graph"
	usersCollection            = "users"
	itemsCollection            = "items"
	viewedEdgeCollection       = "viewed"
	purchasedEdgeCollection    = "purchased"
	droppedEdgeCollection      = "dropped"
	deliveredEdgeCollection    = "delivered"
	notDeliveredEdgeCollection = "not_delivered"
)

// Store --
type Store struct {
	DB               driver.Database
	Graph            driver.Graph
	Users            driver.Collection
	Items            driver.Collection
	ViewedEdge       edgeCollection
	PurchasedEdge    edgeCollection
	DroppedEdge      edgeCollection
	DeliveredEdge    edgeCollection
	NotDeliveredEdge edgeCollection
}

type edgeCollection struct {
	Collection  driver.Collection
	Constraints driver.VertexConstraints
	Label       string
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
	edgeDefinitions := []driver.EdgeDefinition{}
	edgeDefinitions = append(edgeDefinitions,
		createEdgeDefinitions([]string{viewedEdgeCollection, purchasedEdgeCollection, droppedEdgeCollection}, []string{usersCollection}, []string{itemsCollection})...)
	edgeDefinitions = append(edgeDefinitions,
		createEdgeDefinitions([]string{deliveredEdgeCollection, notDeliveredEdgeCollection}, []string{itemsCollection}, []string{usersCollection})...)

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = edgeDefinitions
	graph, err := db.CreateGraph(nil, graphName, &options)
	if err != nil {
		log.Fatalf("failed to create a graph %s: %v", graphName, err)
	}
	return graph
}

func initStore(db driver.Database, graph driver.Graph) *Store {
	return &Store{
		DB:               db,
		Graph:            graph,
		ViewedEdge:       initEdgeCollection(graph, viewedEdgeCollection),
		PurchasedEdge:    initEdgeCollection(graph, purchasedEdgeCollection),
		DroppedEdge:      initEdgeCollection(graph, droppedEdgeCollection),
		DeliveredEdge:    initEdgeCollection(graph, deliveredEdgeCollection),
		NotDeliveredEdge: initEdgeCollection(graph, notDeliveredEdgeCollection),
		Users:            initCollection(graph, usersCollection),
		Items:            initCollection(graph, itemsCollection),
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
		Label:       name,
	}
}

func createEdgeDefinitions(edgeCollections []string, fromVertices []string, toVertices []string) []driver.EdgeDefinition {
	edgeDefinitions := []driver.EdgeDefinition{}
	for _, edge := range edgeCollections {
		var definition driver.EdgeDefinition
		definition.Collection = edge
		definition.From = fromVertices
		definition.To = toVertices
		edgeDefinitions = append(edgeDefinitions, definition)
	}

	return edgeDefinitions
}
