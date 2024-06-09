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
handler.AccountHandler
### Accounts Subdomain
```
name: accounts
endpoint: /v1/accounts
method: GET
parameters: None
handler: HandleGetAccounts
description: returns array of accounts using daos model query.

name: register
endpoint: /v1/accounts/register
method: POST
parameters: None
handler: HandleAddAccount
description: adds account and returns json with account data (removes encrypted password)

name: verify
endpoint: /v1/accounts/verify
method: PUT
parameters: None
handler: HandleVerifyAccount
description: verifies if user email exists in the database and provides simple 200 valid response if true.

name: id
endpoint: /v1/accounts/:id
method: GET
rqeuired parameters: id
handler: HandleGetAccountsById
description: returns account information if it exists.

name: attach
endpoint: /v1/accounts/attach
method:  POST
parameters:
handler: HandleAttachAccount
description: attaches account to parent id? creates account if email not found?

name: parent id
endpoint: /v1/accounts/attached/:parent_id
method: GET
required parameters: parent_id
handler: HandleGetAttachedAccounts
description: returns record of parent id, if it exists.

name: update
endpoint: /v1/accounts/update
method: PUT
required parameters: infoType (only accepts 'personal' or 'authentication') 
handler: HandleUpdateAccount
description: updates account information

name: login
endpoint: /v1/accounts/login
method: POST
parameters: None
handler: HandleLogIn
description: validates the account. 
```
### Events Subdomain
```
name: events
endpoint: /v1/events
method: GET
required parameters: healthSpecialistId or healthSpecialistId (can only chooose one, else 400 error)
optional parameters: after 
handler: HandleGetEvents
description: returns records of parient or specialist events. 'after' optional parameter provides a cutoff for event's date

name: schedule
endpoint: /v1/events/schedule
method: POST
parameters: None
handler: HandleScheduleEvent
description: schedules an event, returns scheduled event

name: reschedule
endpoint: v1/events/reschedule
parameters: 
handler: HandleRescheduleEvent
description: reschedules event to next date.
```
### Planner Subdomain
```
name: add meal
endpoint: /v1/planner/addMeal
method: POST
parameters:
handler: HandleAddMeal
description: add meal to meals table, if it exists.

name: add meal plan
endpoint: /v1/planner/addMealPlan
method: POST
parameters:
handler: HandleAddMealPlan
description: wtf, ask angelo

name: get meal
endpoint: /v1/planner/getMeal
method: GET
required parameters: healthSpecialistId or mealId (can only chooose one, else 400 error)
handler: HandleGetMeal
description:  returns meal or list of meals depending on parameters

name: get meal plan
endpoint: /v1/planner/getMealPlan
method: GET
required parameters: healthSpecialistId or patientId (can only chooose one, else 400 error)
handler: HandleGetMealPlan
description: returns meal plan or list of meals depending on parameters.
```
