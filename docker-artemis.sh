docker run --name local-artemis-reactandgo -p 1883:1883 -p 8161:8161 -e ARTEMIS_USER=artemis1 -e ARTEMIS_PASSWORD=artemis1 angular2guy/artemis-ubuntu

docker start local-artemis-reactandgo

docker stop local-artemis-reactandgo