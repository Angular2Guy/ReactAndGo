workspace "ReactAndGo" "This is project to explore the integration of a React Frontend and a Go/Gin/Gorm Backend" {

    model {
        user = person "User"
        gasPriceFinderSystem = softwareSystem "Gas Price Finder System" "System to import and process gasprices and alert users of cheap prices" {
            gasPriceFinder = container "Gas Price Finder" "App to integrate the React SPA and the Go/Gin/Gorm Backend" {
                reactFrontend = component "React Frontend" "The SPA to show the current gas prices, latest price updates and send alerts for cheap prices." tag "Browser"
                backendCron = component "Scheduler" "Start the scheduled jobs." tag "Scheduler"
                backendMessageConsumer = component "Message Consumer" "Consume the MQTT gas price messages." tag "Consumer"
                backendGasStationClient = component "Gas station client" "Rest Request for the current gas stations."
                backendGasStationController = component "Gas station controller" "Provides the interfaces for the gas station related requests."
                backendJwtTokenService = component "Jwt Token service" "Provides the Jwt Token handling for security."
                backendAppUserController = component "App user controller" "Provides the interfaces for the app user related requests."
                backendNotificationController = component "Notification controller" "Provides the interfaces for the notification related requests."
                backendGasStationService = component "Gas station service" "Implements the gas station related logic."
                backendNotificationService = component "Notification service" "Implements the notification related logic."
                backendAppUserService = component "App user service" "Implements the app user related logic (login/signin/logout)."
            }
            database = container "Postgresql Db" "Postgresql Db stores all the data of the system." tag "Database"
        }
        gasStationProviderSystem = softwareSystem "Gas Station Provider System" "System that provides a rest interface to the currently supported gas stations."
        gasPriceProviderSystem = softwareSystem "Gas Price Provider System" "System that sends MQTT messages with the gas price updates."

        # relationships system context
        user -> gasPriceFinderSystem "Uses"
        gasPriceFinderSystem -> gasStationProviderSystem
        gasPriceFinderSystem -> gasPriceProviderSystem
        
        # relationships containers
        user -> gasPriceFinder
        gasPriceFinder -> database
        gasPriceFinder -> gasStationProviderSystem "requests the gas stations once a day to update db."
        gasPriceFinder -> gasPriceProviderSystem "receives/processes the gas price update messages."

        # relationships components
        reactFrontend -> backendGasStationController "rest requests"
        reactFrontend -> backendAppUserController "rest requests"
        reactFrontend -> backendNotificationController "rest requests"
        backendGasStationController -> backendJwtTokenService
        backendAppUserController -> backendJwtTokenService
        backendNotificationController -> backendJwtTokenService
        backendGasStationController -> backendGasStationService
        backendAppUserController -> backendAppUserService
        backendNotificationController -> backendNotificationService
        backendCron -> backendGasStationClient "trigger gas stations update."
        backendGasStationClient -> backendGasStationService
        backendMessageConsumer -> backendGasStationService "gas price updates"
        backendGasStationService -> backendNotificationService "create notifications from gas price updates."
    }

    views {
        systemContext gasPriceFinderSystem "ContextDiagram" {
            include *
            autoLayout
        }

        container gasPriceFinderSystem "Containers" {
        	include *
            autoLayout
        }

        component gasPriceFinder "Components" {
        	include *
            autoLayout
        }   

        styles {
        	element "Person" {            
            	shape Person
        	}        	
        	element "Database" {
                shape Cylinder                
            }
            element "Browser" {
                shape WebBrowser
            }
            element "Scheduler" {
            	shape Circle
            }
            element "Consumer" {
            	shape Pipe
            }
    	}
    }

}