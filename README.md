# ReactAndGo

The ReactAndGo project is used to compare a single page application frontend based on React and a Rest backend based on Go, to Angular frontends and Spring Boot/Java backends. 

## Project Status

The ReactAndGo project is used to compare a single page application frontend based on React and a Rest backend based on Go, to Angular frontends and Spring Boot/Java backends. 

The goal of the project is to send out notifications to car drivers if the gas price falls below their target price. The gas prices are imported from a provider via MQTT messaging and stored in the database. 

For development 2 test messages are provided that are send to an Apache Artemis server to be processed in the project. The Apache Artemis server can be run as Docker image and the commands to download and run the image can be found in the 'docker-artemis.sh' file. As database Postgresql is used and it can be run as Docker image too. The commands can be found in the 'docker-postgres.sh' file.

The frontend provides a login/signin dialog that uses React and MUI components.

The combined build is implemented with a MakeFile

## Articles
[The ReactAndGo Architecture and Gorm DB access](https://angular2guy.wordpress.com/2023/02/26/the-reactandgo-architecture-and-gorm-db-access/)

[Notifications from React frontend to Go/Gin/Gorm backend](https://angular2guy.wordpress.com/2023/03/09/notifications-from-react-frontend-to-go-gin-gorm-backend/)

## C4 Architecture Diagrams
The project has a [System Context Diagram](structurizr/diagrams/structurizr-1-SystemContext.svg), a [Container Diagram](structurizr/diagrams/structurizr-1-Containers.svg) and a [Component Diagram](structurizr/diagrams/structurizr-1-Components.svg). The Diagrams have been created with Structurizr. The file runStructurizr.sh contains the commands to use Structurizr and the directory structurizr contains the dsl file.