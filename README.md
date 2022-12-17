# Generation
Reading the book "Cloud Native Go(edition O'REILLY)" by Matthew A.Titmus(Thank you), following along I applied Matthew's teaching and ended up with a key value store of type "Least Recently Used"". This kind of Key-Value storage tracks how recently each of its keys have been used and delete the (oldest) one standing outside the predefined maximum records. It has the advantage to make sure the predefined capacity of the Key-Value-Store is respected, keeping only the most recent records. This project is mean to be a Cloud Native service so it needed discovery and data-replications between its various instances, luckily for me I was reading a second book "Distributed Services with Go(edition O'REILLY)" by Travis Jeffrey(Thank you) which taught me how to implement it.


## Generation features
- Storage: The records are stored via a Doubly-Linked-List data structure, which is itself wrapped into a ShardedMap pattern, preventing potential bottleneck resulting from lock contention. Note that the number of shards and number of records per shard are configurable.

- Transport: HTTP and gRPC transport layer protocol are supported. But right now using HTTP does not allow the replication of the datas when using few replicas(see the next point).

- Replication: This key value store is made to run as a cloud native service so replication has been implemented. Gossip protocol(using hashicorp/serf library) allows service discovery and orchestration. The replication of the datas(between nodes) uses the Raft protocol(using hashicorp/raft library) to provide consensus. 

- Security: some TLS certificates protect the datas and secure the connections between the various end-point.

- Observability: For observability concerns tracing(Jaeger) and metrics(Prometheus) are already instrumented into the code. Still they remain optional so if they(or only one of them) are needed, add the corresponding flag.


## Configuration flags
```
  -s, --shards          number of shards (default 10)
  -i, --itemPerShard    number of shards (default 100)
  -d, --dbLogger        enable the database logging (default disabled)
  -m, --isMetrics       enable Prometheus metrics (default disabled)
  -t, --isTracing       enable Jaeger tracing (default disabled)
  -l, --loggerMode      logger mode can be prod, development, debug (default "prod")
  -j, --jaeger          the Jaeger endpoint to connect (default "http://jaeger:14268/api/traces")
```


## Start the Key-Value-Store
- `docker-compose up --build` build(add the configuration flags needed in the Dockerfile) and run the Key Value Store container, a Postgres  container, a Jaeger container, a Prometheus container. If the tracing has been enabled see the traces at `http://localhost:16686/search`, if the metrics has been enabled see them at `http://localhost:9090`  
- As Transaction-logs, Tracing and Metrics are by default disabled, the Key Value Store will work fine with
```
cd cmd
go run .
```


## Available requests
- The actual available requests
```
// add a value
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
client put key value

// get the value
curl -X GET http://localhost:8080/v1/key-a
client get key

// delete a value
curl -X DELETE http://localhost:8080/v1/key-a
client delete key

// get all keys in storage
curl -X GET http://localhost:8080/v1/util/keys
client getkeys "" ""

// get all keys values in storage(is a stream), no implementation for HTTP
client getkeysvalues "" ""
```
- A few files to stress test the service are available in `testScript`



