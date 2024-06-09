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
│   │
|   └── main.go #Entry point for the application
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

### Run Locally

Navigate to the ```cmd``` directory and use the ```go run main.go serve``` command
The command will spin up two services:
- the api at localhost:8090/api
- the admin dashboard at localhost:8090/_/

In order to access the admin dashboard you will need to register yourself.

## API

### Accounts Subdomain
```
name: accounts
endpoint: /v1/accounts
method: GET
parameters: 

name: register
endpoint: /v1/accounts/register
method: POST
parameters: 

name: verify
endpoint: /v1/accounts/verify
method: PUT
parameters: 

name: id
endpoint: /v1/accounts/:id
method: GET
parameters:

name: attach
endpoint: /v1/accounts/attach
method:  POST
parameters:

name: parent id
endpoint: /v1/accounts/attached/:parent_id
method: GET
parameters:

name: update
endpoint: /v1/accounts/update
method: PUT
parameters:

name: login
endpoint: /v1/accounts/login
method: POST
parameters:
```
### Events Subdomain
```
name: events
endpoint: /v1/events
method: GET
parameters:

name: schedule
endpoint: /v1/events/schedule
method: POST
parameters: 
```
### Planner Subdomain
```
name: add meal
endpoint: /v1/planner/addMeal
method: POST
parameters:

name: add meal plan
endpoint: /v1/planner/addMealPlan
method: POST
parameters: 

name: get meal
endpoint: /v1/planner/getMeal
method: GET
parameters: 

name: get meal plan
endpoint: /v1/planner/getMealPlan
method: GET
parameters: 
```
