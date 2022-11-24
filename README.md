# GoLRU (or Generation0 ?)
Reading the book "Cloud Native Go(edition O'REILLY)" by Matthew A.Titmus(thank you very much), following along I applied Matthew's teaching and ended up with a key value store of type "Least Recently Used"". This kind of Key-Value storage tracks how recently each of its keys have been used and delete the (oldest) one standing outside the predefined boundary memory space. It has the advantage to make sure the predefined capacity of the Key-Value-Store is respected, keeping only the most recent records. 


## GoLRU features
- An optional Transaction logs(into db or file-system or both) maintains a history of mutating changes executed by the data store, this feature allows, in case of service crashes for eg, the service to replay the transactions to reconstruct its functional state.
- The data store is wrapped into a ShardedMap pattern, preventing potential bottleneck resulting from lock contention. Note that the number of shards and number of records per shard are configurable.
- HTTP and gRPC transport layer protocol are supported.
- For observability concerns tracing(Jaeger) and metrics(Prometheus) are already instrumented into the code. Still they remain optional so if they(or only one of them) are needed, add the corresponding flag.


## Configuration flags
```
  -s, --shards          number of shards (default 10)
  -i, --itemPerShard    number of shards (default 100)
  -d, --dbLogger        enable the database logging (default disabled)
  -f, --fileLogger      enable the file logging (default disabled)
  -m, --isMetrics       enable Prometheus metrics (default disabled)
  -t, --isTracing       enable Jaeger tracing (default disabled)
  -l, --loggerMode      logger mode can be prod, development, debug (default "prod")
  -j, --jaeger          the Jaeger endpoint to connect (default "http://jaeger:14268/api/traces")
  -e, --encryptK        an encoding key to encrypt data to file logs (default "xxxxxxxxxxx")
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

// get the value
curl -X GET http://localhost:8080/v1/key-a

// get all keys in storage
curl -X GET http://localhost:8080/v1/util/keys

// delete a value
curl -X DELETE http://localhost:8080/v1/key-a
```

- A few files to stress test the service are available in `testScript`


## Conclusion
This is not really a project(at this point) but more a fun try to apply what I have learned, however I personally like the result.






