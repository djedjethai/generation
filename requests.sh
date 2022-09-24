#!/bin/bash

curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e
curl -X PUT -d 'Hello, key-value store!' -v http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeA2!' -v http://localhost:8080/v1/key-a
curl -X PUT -d 'Hello, key-value storeB!' -v http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-a
curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeC!' -v http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
# curl -X DELETE http://localhost:8080/v1/key-b
curl -X PUT -d 'Hello, key-value storeD!' -v http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-c
curl -X PUT -d 'Hello, key-value storeE!' -v http://localhost:8080/v1/key-e
curl -X GET http://localhost:8080/v1/key-a
curl -X GET http://localhost:8080/v1/key-b
curl -X GET http://localhost:8080/v1/key-c
curl -X GET http://localhost:8080/v1/key-d
curl -X GET http://localhost:8080/v1/key-e
curl -X DELETE http://localhost:8080/v1/key-c
curl -X DELETE http://localhost:8080/v1/key-e


echo "finex, should show C and E"





