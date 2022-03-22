# AuthSvc

An User **Auth**entication **S**er**v**i**c**e that will allow users to Login via Email/Phone and Password combination and receive a short lived Access Token that will allow them to access some authenticated routes in other services.

## Outline

- [AuthSvc](#authsvc)
  - [Outline](#outline)
  - [Postman collection](#postman-collection)
  - [Endpoints](#endpoints)
  - [Project run instructions](#project-run-instructions)
  - [Project Structure](#project-structure)
    - [Configuration](#configuration)
      - [Example](#example)
    - [Structure](#structure)
      - [Files/Folders Map](#filesfolders-map)

## Postman collection

Import the url (<https://www.getpostman.com/collections/f4dc6a39771cb8945120>) into Postman to get the API collection and test the endpoints. Please follow the [link][1] to know more.

## Endpoints

> Note: Use the [Postman collection](#postman-collection) to test the endpoints.

|Endpoint|Description|Method|Authorization|Request body Example|Response body Example|
|--------|-----------|------|-------------|---------------|----------------|
| /  | Home page containing server configurations | **GET** | N/A |  | ```<html>...</html>```
| _/auth/register_ | To register a user. If user is successfully registered, an email verification link  will be sent to the registered email | **POST** | N/A | <code>{"firstname": "Test",<br>"lastname": "User",<br>"email": "test.user1@testmail.com",<br>"password": "giv_Me_1_Pine@pple"}</code> | Please check your email to verify. <br> **Note**: Check the [SMTP Mock server](http://localhost:8025) to get email verification link |
| */auth/email_verification?email=$email&verfication_code=$verificationCode* | To verify the email | **GET** | N/A | | _user is successfully verified!!_ |
| _/auth/login_ | To login a user. After a successful login, user will get an access token and a refresh token | **POST** | N/A | <code>{"email": "admin.user@testmail.com", "password": "_LaRa08CRoft"}</code> | <code>{"access_token": "eyJhbGc...", "refresh_token": "eyJhI....", "token_type":"bearer",<br>"expires": 300}</code> |
| _/auth/token/verify_ | To verify an Access Token. Verified Access token will return the User's profile, role, permission etc. | **POST** | N/A | <code>{"access_token": "eyJhbGciO..."}</code> | <code>{"firstname": "Admin",<br>"lastname": "User",<br>"email": "admin.user@testmail.com",<br>"roles": ["Admin"],<br>"permissions": ["GetPost", "AddPost", "UpdatePost", "DeletePost"]}</code> |
| _/auth/token/refresh_ | To acquire a new Access Token using the Refresh Token generated upon Login | **POST** | N/A | <code>{"refresh_token": "eyJhbGciO..."}</code> | <code>{"access_token": "eyJhbGciO...",<br>"refresh_token": "eyJhbG...",<br>"token_type": "bearer",<br>"expires": 300}</code> |

## Project run instructions
<!-- + change Server -> Bind of **app.json**
+ change Db -> Password of **app.json** -->
- run **`go get` or `go mod download`** command
- Optional run `go mod vendor` to make copies of all packages needed to support builds and tests of packages in the main module
- run **`go build`**
- run **`./authsvc ${configFilePath}`** if no configFilePath is provided default [config](./authsvc.log) will be used instead.

## Project Structure

### Configuration

Service depends on configuration so that service can be run with different environment configurations. This configuration is passed in  **CLI arguments** when the service is going to run (*see the [run instructions](#project-run-instructions)*). e.g. configuration - [authsvc.json][2]

#### Example

```json
{
  "Name": "AuthSvc", // service name
  "Description": "Configuration for the authentication service.", // service description
  "AllowCORS": true, // enable or disable CORS
  "Server": { // server configuration
    "Bind": "", // binding address
    "Port": 8080, 
    "SSLCertificate": { // SSL certificates definition to allow HTTPS requests
      "ServerKey": "/home/shuvojit-kaz/Desktop/Learning/auth-system/certificates/server.key",
      "ServerCrt": "/home/shuvojit-kaz/Desktop/Learning/auth-system/certificates/server.crt"
    }
  },
  "DB": { // Database configurations
    "User": "user", // database user
    "Password": "password", // database user's password
    "Host": "localhost", // database host address
    "Port": 3306, // database port
    "Database": "AuthDB" // database name
  },
  "TokenDB": { // Redis cache token database configuration
    "Host": "localhost", // database host address
    "Port": 6379, // database port
    "Password": "password", // database password
    "Database": 1 // redis database
  },
  "JWTDef": { // JWT token definition
    "AccessToken": { // Access token
      "Secret": "#LaRa_cR0ft$", // Secret
      "Exp": 5 // Expire time in Minutes
    },
    "RefreshToken": { // Refresh token
      "Secret": "scr1bus1nt3rp@r3s",  // Secret
      "Exp": 10 // Expire time in Minutes
    }
  },
  "SmtpServer": { // SMTP server definition
    "Host": "localhost", // Server address
    "Port": 1025, // Port
    "from": "authsvc@testmail.com" // client email address
  },
  "Logging": { // logging definition
    "Filename": "./authsvc.log", // log file path
    "Level": "DEBUG" // log level
  },
  "Indent": true // Enable/disable HTTP response body indentation
}
```

### Structure

Using [Clean Architecture][1] to structure the go project's files and folders

![Clean Architecture][4]

#### Files/Folders Map

```
├── cache                <- cache database repository module (redis)
│   ├── auth.go          <- refresh token store
│   └── tokendb.go       <- connection setup and managing connection instance
├── cfg                  <- project configuration module related on authsvc.json
│   ├── config.go        
├── db                   <- database repository module (MySQL)
│   ├── authdb.go        <- authdb connection setup and managing connection instance
│   └── permission.go    <- Permission store
│   └── role.go          <- Role store
│   └── permission.go    <- User store
├── email                <- SMTP email client module
│   ├── emailclient.go   <- Use for sending new mail
├── log4u                <- logging module; much like log4j has
│   ├── log4u.go
├── render               <- HTTP response renderer module
│   └── jsonrenderer.go  <- HTTP JSON response definition
│   └── renderer.go      <- Renderer interface
├── resource             <- REST API endpoints's (resource) request handler module
│   └── auth.go          <- Request handlers for auth resource e.g. /auth
│   └── common.go        <- resource utility
│   └── errors.go        <- HTTP request ERROR responses
│   └── home.go          <- / endpoint request handler
│   └── protect.go       <- Route protector
│   └── token.go         <- Request handlers for token resource e.g. /auth/token
└── route                <- Route builder module
│   └── routebuilder.go
├── table                <- Database entity/tables
│   └── permission       <- Permission table module consists of its definition and related DB operations
│       └── table.go     
│   └── role             <- Role table module consists of its definition and related DB operations
|       └── table.go
│   └── user             <- User table module consists of its definition and related DB operations
|       └── table.go
└── token                <- token service module
│   └── service.go
│   └── token.go
└── uc                   <- Use cases
│   └── adm              <- Admin related use cases
│       └── handler.go     
│   └── permission       <- Permission related use cases
|       └── handler.go
│   └── role             <- Role related use cases
|       └── handler.go
│   └── token            <- Token related use cases
|       └── handler.go
│   └── user             <- User related use cases
|       └── handler.go
│   └── common.go        <- Use case utilities
└── validator            <- validator module with custom validators's tag e.g. `validPwd`
│   └── validator.go
└── authsvc.go           <- entry point of the service
└── authsvc.json         <- service config
└── go.mod               <- list dependent packages
└── go.sum               <- list checksum of downloaded go modules and their dependencies
└── version.go           <- project versioning
```

[4]: https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg
[3]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[2]: ./authsvc.json
[1]: https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#:~:text=to%20import%20your%20api%20specifications%20into%20postman%3A
