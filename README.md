# WellnessWaveApi

### Folder Structure
```plaintext
/WellnessWaveApi
│
├── cmd
│   ├── service1  # Entry point for the first service
│   │   └── main.go
│   │
│   └── service2  
│
├── internal
│   ├── service1  # Internal package for the first service
│   │   ├── app
│   │   │   └── handler
│   │   │       └── ...  Contains HTTP handlers that process requests and return responses. Handlers often marshal and unmarshal data, validate requests, and call upon services to perform business logic.
│   │   ├── model
│   │   │   └── ...  Defines the data structures used by the application. These models represent the schema of your data and are used across your application.
│   │   ├── service
│   │   │   └── ...  Contains the business logic of your application. Services use models and repositories to perform operations and implement the core functionalities of your application.
│   │   └── repository
│   │       └── ...  The data access layer responsible for interacting with the database. Repositories abstract the data source and provide a clean API for fetching and storing data.
│   │
│   └── service2  # Internal package for the second service
│
├── pkg  # Shared packages used across services
│   └── ...
│
├── config  # Configuration files and management
│   └── ...
│
├── scripts  # Scripts for tasks like DB migrations
│   └── ...
│
├── migrations  # Database migration files
│   └── ...
│
└── tests  # Test files
    ├── service1  # Tests specific to service1
    │   └── ...
    └── service2  # Tests specific to service2
        └── ...
```
