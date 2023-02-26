# AngularAndGo

This is a learning project to test how an Angular frontend and a Golang backend can be integrated. 

## Project Status

Experimental project to gather experience.
The backend provides a signin and login for the users. It uses currently a Postgresql DB for storage and the file docker-postgresql.sh contains the Docker commands to download and start/stop it. The Messaging is provided by an Apache ArtemisMQ server. The docker-artemis.sh file contains the Docker commands to download start/stop it. The application uses the Gin framework to provide the rest endpoints for the frontend and uses the Golang-Jwt library for the tokens. The Messages are processed with the Paho-MQTT library. 
The frontend provides a login/signin dialog that uses React and MUI components.
A combined build is still open.

## Articles
[The ReactAndGo Architecture and Gorm DB access](https://angular2guy.wordpress.com/2023/02/26/the-reactandgo-architecture-and-gorm-db-access/)