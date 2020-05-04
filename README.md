# Ekas Portal RESTful Application

The api provides the following features right out of the box 

* RESTful endpoints in the widely accepted format
* Standard CRUD operations of a database table
* JWT-based authentication
* Application configuration via environment variable and configuration file
* Structured logging with contextual information
* Panic handling and proper error response generation
* Automatic DB transaction handling
* Data validation
* Full test coverage
 
The api uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing framework: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing/v2)
* Database: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [logrus](https://github.com/sirupsen/logrus)
* Configuration: [viper](https://github.com/spf13/viper)
* Dependency management: [dep](https://github.com/golang/dep)
* Testing: [testify](https://github.com/stretchr/testify)



Now you can build and run the application by running the following command under the
`$GOPATH/ekas-portal-api` directory:

```shell
go run server.go
```

or simply the following if you have the `make` tool:

```shell
make
```

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `GET /ping`: a ping service mainly provided for health check purpose
* `POST /v1/auth`: authenticate a user

For example, if you access the URL `http://localhost:8080/ping` in a browser, you should see the browser
displays something like `OK v0.1#bc41dce`.

If you have `cURL` or some API client tools (e.g. Postman), you may try the following more complex scenarios:

```shell
# authenticate the user via: POST /v1/auth
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:8080/v1/auth
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the artist resources, such as: GET /v1/artists
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/artists
# should return a list of artist records in the JSON format
```

## Next Steps

In this section, we will describe the steps you may take to make use of this starter kit in a real project.
You may jump to the [Project Structure](#project-structure) section if you mainly want to learn about 
the project structure and the recommended practices.



### Implementing CRUD of Another Table
 
To implement the CRUD APIs of another database table (assuming it is named as `album`), 
you will need to develop the following files which are similar to the `artist.go` file in each folder:

* `models`: contains the data structure representing a row in the new table.
* `services`: contains the business logic that implements the CRUD operations.
* `daos`: contains the DAO (Data Access Object) layer that interacts with the database table.
* `apis`: contains the API layer that wires up the HTTP routes with the corresponding service APIs.

Then, wire them up by modifying the `serveResources()` function in the `server.go` file.

### Implementing a non-CRUD API

* If the API uses a request/response structure that is different from a database model,
  define the request/response model(s) in the `models` package.
* In the `services` package create a service type that should contain the main service logic for the API.
  If the service logic is very complex or there are multiple related APIs, you may create
  a package under `services` to host them.
* If the API needs to interact with the database or other persistent storage, create
  a DAO type in the `daos` package. Otherwise, the DAO type can be skipped.
* In the `apis` package, define the HTTP route and the corresponding API handler.
* Finally, modify the `serveResources()` function in the `server.go` file to wire up the new API.

## Project Structure

This starter kit divides the whole project into four main packages:

* `models`: contains the data structures used for communication between different layers.
* `services`: contains the main business logic of the application.
* `daos`: contains the DAO (Data Access Object) layer that interacts with persistent storage.
* `apis`: contains the API layer that wires up the HTTP routes with the corresponding service APIs.

[Dependency inversion principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
is followed to make these packages independent of each other and thus easier to test and maintain.

The rest of the packages in the kit are used globally:
 
* `app`: contains routing middlewares and application-level configurations
* `errors`: contains error representation and handling
* `util`: contains utility code

The main entry of the application is in the `server.go` file. It does the following work:

* load external configuration
* establish database connection
* instantiate components and inject dependencies
* start the HTTP server

Genretae Cert (.pem)
    openssl req -newkey rsa:2048 \
    -new -nodes -x509 \
    -days 3650 \
    -out cert.pem \
    -keyout key.pem \
    -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost"