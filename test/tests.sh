#!/bin/sh
### basic api test. Other tools should be used for a more accurate testing

apk update && apk add curl jq

echo "---Getting category list---"
curl -s -X GET "http:/localhost:8000/categories" |jq

echo "---Adding scooter a category in the list---"
curl -s -X POST "http:/localhost:8000/categories" \
   -H 'Content-Type: application/json' \
   -d '{"title":"scooter"}'
sleep 1
curl -s -X GET "http:/localhost:8000/categories" |jq

echo "---Deleting scooter category from the list---"
curl -s -X DELETE "http:/localhost:8000/categories/3" 
sleep 1
curl -s -X GET "http:/localhost:8000/categories" |jq
