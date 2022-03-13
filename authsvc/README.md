# AuthSvc

An User **Auth**entication **S**er**v**i**c**e that will allow users to Login via Email/Phone and Password combination and receive a short lived Access Token that will allow them to access some authenticated routes in other services.

## Postman collection

Import the url <https://www.getpostman.com/collections/f4dc6a39771cb8945120> into Postman to get the API collection. Please follow the [link][1] to know more.

## Project run instructions
<!-- + change Server -> Bind of **app.json**
+ change Db -> Password of **app.json** -->
+ run **`go get` or `go mod download`** command
+ Optional run `go mod vendor` to make copies of all packages needed to support builds and tests of packages in the main module
+ run **`go build`**
+ run **`./authsvc ${configFilePath}`** if no configFilePath is provided default [config](./authsvc.log) will be used instead.


[1]: https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#:~:text=to%20import%20your%20api%20specifications%20into%20postman%3A
