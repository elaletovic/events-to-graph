package main

import (
	"context"
	"fmt"

	"github.com/elaletovic/events-to-graph/config"
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
	//load config
	cfg := config.Load()

	conn := graph.Connect(cfg.DBAddress)

	client := graph.GetClient(conn)

	store := graph.Init(client)

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
	//handlers for processors (save to SQL and graph DBs)
	processor := processors.NewEventProcessor(store)

	//save events to graph
	for _, topic := range generators.Topics {
		router.AddNoPublisherHandler(
			fmt.Sprintf("save_to_graph_%s_handler", topic),
			topic,
			pubSub,
			processor.SaveToGraph,
		)
	}

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

	go generators.GenerateEvents(pubSub)

	//run the router
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
