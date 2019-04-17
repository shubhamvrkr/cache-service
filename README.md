# Cache Service

Cache service is an simple golang project to store data items (employees) into the cache and also in persistent storage. Mongo DB is used as persistent storage and [freecache](https://github.com/coocood/freecache) library is used as in memory cache library
________________________________________________________
## Installation

- Git clone the project using
```
git clone https://github.com/shubhamvrkr/cache-service.git
```
- Run
```` cd cache-service ````
and then run
````  go install ```` and then run the executable created in bin folder
- You can also run by running the docker-compose.yml file using the command
```` docker-compose -f docker-compose.yml up -d ````
- You can also deploy the deployment files using kubectl commands present in kube-deployments folder
- Go to [http://localhost:8080/swaggerui/]() for API details
_______________________________________________________

## Project Structure

### Architecture

![Architecture](docs/arch.png)

The cache service contains multiple packages. Each package contains test case script written that can be executed by running
```` go test ```` command. Below is the description of each package

- #### Config
   - This package loads the configuration from yaml file for database, cache memory to be used, messaging queue, server. It also checks if environment variable are set, if so environment variables overrides config.yml properties
- #### Database
    - This package handles all database operation such as save, get by primary key, get by custom filter criteria.
- #### Cache
    - This package handles two types of cache operations such as set and get, it stores the bytes of employee along with the employee id as a key in memory.
- #### Messaging Queue
    - This package handles sending messages to messaging queue and receiving the message from messaging queue. RabbitMQ is used as a messaging queue to demonstrate reloading of cache data on reload event trigger through API.
- #### Handler
  - This package implements API functionalities such as POST employee, GET employee by ID, Get employees by gender and trigger reload data event to messaging queue. This package depends on database package, cache package and messaging package
  - POST employee API: stores the employee data in the DB store and also updadtes the cache
  - GET employee by ID: Hits the cache for employee details, if miss loads from database and update cache, if cache hit directly returns employee.
  - Get employee by Gender: Gets list of employee IDs matching the query criteria, loads employee details for each from cache, if not present in cache loads from database and update caches, finally return details of all employees
  - Trigger event: submits a reload event in the messaging queue which is listened by an listener on which data in the cache is being loaded from database
 ________________________________________________________
## How to try

- After setting up, open the swagger UI
- There are four API's provided, One for posting employee details, one for getting employee details by ID, one for getting employees based on gender and last one is to trigger a reload event on which cache data is reloaded from database.
- Employee ID is assumed to be int32, Gender may be single string (i.e M/F)
- Getting all the employees based on gender API demonstrate pagination. You can set limit for the no. of employee details to be fetched, API return next URL string along with limited no of employee. For stable pagination lastid of the employee fetched is used as a query params. Initlally it can be set to -1 assuming this wont be a valid employee ID. For next page result next url returned by previous API call can be used.
 ________________________________________________________
