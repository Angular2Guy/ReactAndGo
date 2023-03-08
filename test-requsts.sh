curl -X GET http://localhost:3000/gasstation/33d6adbc-e6f5-4c4e-f7a5-a132626d3736 -H 'Content-Type: application/json'  -H 'Authorization: Bearer injectJwtTokenString'

curl -X GET http://localhost:3000/config/updatepc?filename=plz-5stellig.geojson.gz -H 'Content-Type: application/json'  -H 'Authorization: Bearer injectJwtTokenString'

curl -X POST http://localhost:3000/gasstation/search/place -H 'Content-Type: application/json' -H 'Authorization: Bearer injectJwtTokenString' -d '{"StationName": "Berlin"}'

curl -X POST http://localhost:3000/gasstation/search/place -H 'Content-Type: application/json' -H 'Authorization: Bearer injectJwtTokenString' -d '{"StationName": "Berlin", "PostCode": "13088"}'

curl -X POST http://localhost:3000/gasstation/search/location -H 'Content-Type: application/json' -H 'Authorization: Bearer injectJwtTokenString' -d '{"Longitude": 13.519067, "Latitude": 52.522686, "Radius": 25.0}'

curl -X POST http://localhost:3000/appuser/signin -H 'Content-Type: application/json' -d '{"Username": "Max123","Password": "Password123","Latitude": 54.824158,"Longitude": 8.346131}'

curl -X POST http://localhost:3000/appuser/login -H 'Content-Type: application/json' -d '{"Username": "Max123","Password": "Password123"}'

curl -X POST http://localhost:3000/appuser/locationradius -H 'Content-Type: application/json' -d '{"Username": "Max123", "Latitude": 12.12, "Longitude": 21.21, "SearchRadius": 10.0}'

curl -X POST http://localhost:3000/appuser/targetprices -H 'Content-Type: application/json' -d '{"Username": "Max123", "TargetDiesel": "1.750", "TargetE10": "1.760", "TargetE5": "1.770"}'