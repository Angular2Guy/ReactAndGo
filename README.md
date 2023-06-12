# ReactAndGo
This is a project to compare a single page application frontend based on React and a Rest backend based on Go, to Angular frontends and Spring Boot/Java backends. It enables the car drivers to find gas prices below an user defined target price. It sends notifications if the target price is reached and shows the prices matches in a table. The gas prices of the location of the user are shown in a table and on a map with with pins and hovers. The frontend uses React/Typescript  with Recoil for the state, Mui for the components and Openlayers for the map. The backend uses Go with Gin for the controllers with Jwt Token for security and Gorm for database access. Postgresql is used as database and for the MQTT messaging is the Paho library used. For development Apache Artemis is used as messaging system. The gas stations are imported every night with help of the go-cron library. 

Technologies: Go/Golang, Gin, Gorm, Paho, React, Recoil, Typescript, Recoil, Mui, Openlayers, Structurizr

## Articles
* [Cron Jobs and MQTT Messaging in Go](https://angular2guy.wordpress.com/2023/03/27/cron-jobs-and-mqtt-messaging-in-go/)
* [The ReactAndGo Architecture and Gorm DB access](https://angular2guy.wordpress.com/2023/02/26/the-reactandgo-architecture-and-gorm-db-access/)
* [Notifications from React frontend to Go/Gin/Gorm backend](https://angular2guy.wordpress.com/2023/03/09/notifications-from-react-frontend-to-go-gin-gorm-backend/)

## Features
1. Automatic database init on startup.
2. Nightly cronjob with http data import.
3. Datamanagement with Gorm for querys, mutations and automatic mapping.
4. Messaging is done with MQTT with the Paho library and Apache Artemis for development. 
5. Serve the frontend with a controller for the frontend and routes to the rest endpoints. 
6. The security is done with Jwt Tokens that can be revoked.
7. The frontend shows the prices matches and the local prices in a Mui Table with React/Typescript.
8. Price matches are shown as notifications. 
9. The local prices are shown on a map with openlayers at the locations of the gas stations.

## Mission Statement 
The ReactAndGo project serves as example for the integration of React, Go, Gin, Gorm and Postgresql in a structured architecture. The build is integrated in one Makefile and the application can be build in a Docker image with the Dockerfile. As documentation are the structurizr diagrams as images and sources available.

## Postgresql setup
The database can be run as Docker image with the commands in the 'docker-postgres.sh' script. 

## Apache MQ Artemis
The Messaging server can be run as Docker image with the commands in the 'docker-artemis.sh' script. 

## C4 Architecture Diagrams
The project has a [System Context Diagram](structurizr/diagrams/structurizr-1-SystemContext.svg), a [Container Diagram](structurizr/diagrams/structurizr-1-Containers.svg) and a [Component Diagram](structurizr/diagrams/structurizr-1-Components.svg). The Diagrams have been created with Structurizr. The file runStructurizr.sh contains the commands to use Structurizr and the directory structurizr contains the dsl file.

## Development Environment
Visual Studio Code with the Go Extension works well as IDE for the frontend and the backend.