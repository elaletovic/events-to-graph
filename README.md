# event-to-graph

An example service used for converting events to useful graph data.

## Description
The service uses [watermill][1] and [arangoDB][2] for showcasing how to convert events, coming from an event stream, to useful graph data, that can be searched and queried easily. It performs two main tasks:
* tries to mimick an event stream by generating series of events 
* converts events to graph data; to vertices and edges

### Generating series of events

This part of the service tries to mimick an event stream by generating events. The domain context of this service is a simple marketplace where users can view and purchase items. Two types of events are generated. Creational events like `user_registered` and `item_created`, and relational events like `item_viewed`, `item_dropped`, `item_purchased`, `item_delivered` and `item_not_delivered`.

### Event to graph data conversion

This part of the service receives events and converts the event payload to graph data. Creational events are converted to matching vertices, while relations events are converted to matching edges between vertices.

## Running the service

The easiest way to run the service is to run `docker-compose up`. After ArangoDB and the service are set up, the service will start generating events and saving them to the graph database. The service will log some info for each step of the process, finishing with `done generating events!`.

If you run the service and the graph database separately, bear in mind that service deletes the database and its data every time the service is started.

## Viewing the graph data

The graph data can be checked out with ArangoDB's web interface on `http://localhost:8529`. Access it with credentials `root:rootpassword` and choose the `example_db` database.

ArangoDB offers the data in three formats: `key-value`, `json` and `graph`. To view the entire graph that the service generated, go to the Graphs section and choose `example_graph`. Part of the graph should be loaded. 

Additionally, you can load the entire graph by clicking on the graph icon in the toolbar. Or you can style your vertices and nodes by adding colors, names, attributes, etc.

### Queries
For more details on how to use AQL, the ArangoDB query language, you can visit this [link][3]. Until then, here are some queries you can try out in the `Queries` tab.

**Get all available fields for available items**
```
FOR item in items
RETURN item
```

**Get information on all users**
```
FOR user IN users
RETURN {
    name: user.name,
    age: user.age
}
```

**Get information on failed deliveries**
```
FOR user IN users
    FOR item, delivery IN 
    INBOUND user not_delivered
    RETURN {
        user_name: user.name,
        item_title: item.title,
        user_address: delivery.address,
        reason: delivery.reason
    }
```
If you want to see that displayed as a graph, try out this one below (to display results as a graph, returning data has to have `_from` and `_to` fields)

```
FOR v in not_delivered
RETURN v 
```

**Get stats on all items**

```
FOR item IN items
    LET p = (FOR e IN 1..1 INBOUND item purchased RETURN 1)
    LET v = (FOR e IN 1..1 INBOUND item viewed RETURN 1)
    LET d = (FOR e IN 1..1 INBOUND item dropped RETURN 1)
    LET del = (FOR e IN 1..1 OUTBOUND item delivered RETURN 1)
    LET n = (FOR e IN 1..1 OUTBOUND item not_delivered RETURN 1)
    RETURN {
        item_id: item.id,
        item_name: item.title,
        viewed: COUNT(v),
        dropped: COUNT(d),
        purchased: COUNT(p),
        delivered: COUNT(del),
        not_delivered: COUNT(n)
        }
```

[1]: https://watermill.io/
[2]: https://www.arangodb.com/
[3]: https://www.arangodb.com/docs/stable/aql/