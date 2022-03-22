# AuthSvc

An User **Auth**entication **S**er**v**i**c**e that will allow users to Login via Email/Phone and Password combination and receive a short lived Access Token that will allow them to access some authenticated routes in other services.

## Outline

- [AuthSvc](#authsvc)
  - [Outline](#outline)
  - [Postman collection](#postman-collection)
  - [Endpoints](#endpoints)
  - [Project run instructions](#project-run-instructions)

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

[1]: https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#:~:text=to%20import%20your%20api%20specifications%20into%20postman%3A
