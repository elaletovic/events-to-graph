package main

import (
	"context"

	"github.com/elaletovic/events-to-graph/graph"
	"github.com/elaletovic/events-to-graph/processors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/elaletovic/events-to-graph/generators"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func main() {
	conn := graph.Connect("http://localhost:8529")

	client := graph.GetClient(conn)

	graph.Init(client)

	//configure router
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Recoverer,
	)

	// init publishers and subscribers, using only one for everything
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	//init handler struct
	eventGeneratorHandler := generators.GeneratorHandler{}

	//configure handlers
	router.AddHandler(
		"initial_events_handler",
		generators.InitialEventsTopic,
		pubSub,
		generators.CheckoutTopic,
		pubSub,
		eventGeneratorHandler.InitialEventsHandler,
	)

	router.AddHandler(
		"purchased_events_handler",
		generators.CheckoutTopic,
		pubSub,
		generators.DeliveryTopic,
		pubSub,
		eventGeneratorHandler.PurchasedEventsHandler,
	)

	//handlers for processors (save to SQL and graph DBs)
	processor := processors.EventProcessor{}
	//sqlPubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	//handle initial events
	router.AddNoPublisherHandler(
		"save_to_graph_initial_events_handler",
		generators.InitialEventsTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	//handle purchase events
	router.AddNoPublisherHandler(
		"save_to_graph_purchase_events_handler",
		generators.CheckoutTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	//handle delivery events
	router.AddNoPublisherHandler(
		"save_to_graph_delivery_events_handler",
		generators.DeliveryTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	go generators.GenerateEvents(pubSub)

	//run the router
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
