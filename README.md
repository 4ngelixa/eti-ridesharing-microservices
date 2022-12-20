## ETI Assignment 1 - 30%
By: Angelica Sim (S10205409D)

## Brief Description
An application to simulate a ride-sharing platform using:
1. Golang
2. MySQL Database
3. HTML, JavaScript, CSS (Front-end)

Assignment Requirements:
- At least 2 microservices
- Persistent storage of information using a Database
- Demonstrate the ability to develop REST APIs
- Make conscientious consideration in designing microservices

## Design Considerations
The following microservices have been implemented:
1. Driver
2. Passenger
3. Trip

Each microservice is able to perform CRUD (Create, Read, Update, Delete) operations.

Benefits of using Microservice architecture:
- Easier to manage due to smaller code base
- Each service scales independently when needed
- Changes can be applied to each service independently
- Each service runs separately and communicates with each other using lightweight mechanisms

The Driver and Passenger microservices can call POST, GET, and PUT HTTP methods.
The Trip microservice can call POST, GET, and PUT HTTP methods, as well as the Passenger and Driver microservices when required.

## Architecture Diagram

## Setup Instructions
1. Install [GO](https://go.dev/dl/) and [MySQL Community Edition](https://dev.mysql.com/downloads/installer/).
2. Launch MySQL Workbench and create a new connection. 
3. Run the following command:
```
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost'
 
```
This will create an account named user with the password 'password'.

4. Run `RideSharingDB.sql` database.

5. Clone the repository. Install [GitHub desktop](https://desktop.github.com/) and/or follow the steps [here](https://docs.github.com/en/desktop/contributing-and-collaborating-using-github-desktop/adding-and-cloning-repositories/cloning-and-forking-repositories-from-github-desktop)

## Utilising the code
1. Install the relevant packages needed to run the code, with the exception of the "strconv" package for drivermain.go:
```
"database/sql"
"encoding/json"
"fmt"
"io/ioutil"
"log"
"net/http"
"strconv"

_ "github.com/go-sql-driver/mysql"
"github.com/gorilla/handlers"
"github.com/gorilla/mux"
 
```
2. Run the following codes either on an IDE or on Command Prompt. Type the following commands:
```
cd ETI_Assg1\Entities\Driver
go run driverMain.go

```
cd ETI_Assg1\Entities\Passenger
go run passengerMain.go

```
cd ETI_Assg1\Entities\Trip
go run tripMain.go

```
3. Once the microservices are running, go to HTML folder > `Main.html` to test the ride-sharing platform!