package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func main() {
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
	eventGeneratorHandler := generatorHandler{}

	//configure handlers
	router.AddHandler(
		"initial_events_handler",
		initialEventsTopic,
		pubSub,
		checkoutTopic,
		pubSub,
		eventGeneratorHandler.InitialEventsHandler,
	)

	router.AddHandler(
		"purchased_events_handler",
		checkoutTopic,
		pubSub,
		deliveryTopic,
		pubSub,
		eventGeneratorHandler.PurchasedEventsHandler,
	)

	//handlers for processors (save to SQL and graph DBs)
	processor := eventProcessor{}
	//sqlPubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	sqlDBPublisher := createSQLPublisher()

	//handle initial events
	router.AddHandler(
		"save_to_sql_initial_events_handler",
		initialEventsTopic,
		pubSub,
		sqlEventsTopic,
		sqlDBPublisher,
		processor.SaveToSQLDB,
	)
	router.AddNoPublisherHandler(
		"save_to_graph_initial_events_handler",
		initialEventsTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	//handle purchase events
	router.AddHandler(
		"save_to_sql_purchase_events_handler",
		checkoutTopic,
		pubSub,
		sqlEventsTopic,
		sqlDBPublisher,
		processor.SaveToSQLDB,
	)
	router.AddNoPublisherHandler(
		"save_to_graph_purchase_events_handler",
		checkoutTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	//handle delivery events
	router.AddHandler(
		"save_to_sql_delivery_events_handler",
		deliveryTopic,
		pubSub,
		sqlEventsTopic,
		sqlDBPublisher,
		processor.SaveToSQLDB,
	)
	router.AddNoPublisherHandler(
		"save_to_graph_delivery_events_handler",
		deliveryTopic,
		pubSub,
		processor.SaveToGraphDB,
	)

	go generateEvents(pubSub)

	//run the router
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
