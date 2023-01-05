curl -X GET http://localhost:3000/gasstation/33d6adbc-e6f5-4c4e-f7a5-a132626d3736 -H 'Content-Type: application/json'

curl -X POST http://localhost:3000/gasstation/search/place/ -H 'Content-Type: application/json' -d '{"StationName": "Berlin"}'

curl -X POST http://localhost:3000/gasstation/search/place/ -H 'Content-Type: application/json' -d '{"StationName": "Berlin", "PostCode": "13088"}'