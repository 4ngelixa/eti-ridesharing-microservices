package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Map for Passenger Records
type Passenger struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	MobileNo    int    `json:"MobileNo"`
	Email       string `json:"EmailAddress"`
}

var passengers map[string]passengerInfo

type passengerInfo struct {
	Title string `json:"Passenger"`
}

var (
	db  *sql.DB
	err error
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passengerMain", allPassengers)
	router.HandleFunc("/api/v1/validatePassengerRecord/{id}", validatePassenger)
	router.HandleFunc("/api/v1/passenger/{passengerid}/{email}", passengerMain).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Function to allow Passengers to call GET, PUT, POST, and DELETE methods
func passengerMain(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Connect to MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")

	//Error handling:
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database successfully opened :)")
	}

	// Retrive passsenger data
	if r.Method == "GET" {
		var getAllPassengers Passenger
		reqBody, err := ioutil.ReadAll(r.Body)

		// Defer the close until the main function is done executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getAllPassengers)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("%s - Error in JSON encoding", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetPassengerRecords(db, params["passengerid"], params["email"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	// Creaste a new passenger record
	if r.Method == "POST" {
		var newPassenger Passenger
		reqBody, err := ioutil.ReadAll(r.Body)

		// Defer the close until the main function is done executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &newPassenger)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("%s - Error in JSON encoding", err)
			} else { // Email validation - Email address cannot be empty
				if newPassenger.Email == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				} else { // Run validation funcs
					if !validatePassengerRecord(db, newPassenger.Email) {
						InsertPassengerRecord(db, newPassenger)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - New Passenger added"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("409 - Duplicate Passenger ID"))
						return
					}
				}
			}
		}

		// Update existing passenger records
	} else if r.Method == "PUT" {
		var updatePassenger Passenger
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &updatePassenger)

			// Name validation - First name cannot be empty.
			if updatePassenger.FirstName == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Invalid First Name"))
				return
			} else { // Run validation funcs
				if !validatePassengerRecord(db, updatePassenger.Email) {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("No passenger found with the email: " + updatePassenger.Email))
				} else {
					EditPassengerRecord(db, updatePassenger.PassengerID, updatePassenger.FirstName, updatePassenger.LastName, updatePassenger.MobileNo, updatePassenger.Email)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger details have been updated"))
					return
				}
			}
		}
	} // Passenger users are unable to delete their account as per requirements.
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("404 - You cannot delete your account due to auditing purposes"))
	}
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	json.NewEncoder(w).Encode(passengers)
}

// Check for duplicate Passenger email addresses in the database
func validatePassengerRecord(db *sql.DB, EML string) bool {
	query := fmt.Sprintf("SELECT * FROM Passenger WHERE Email= '%s'", EML)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger Passenger
	for results.Next() {
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNo, &passenger.Email)
		if err != nil {
			panic(err.Error())
		} else if passenger.Email == EML {
			return true
		}
	}
	return false
}

// Check if Passenger exists (using PassengerID)
func validatePassengerID(db *sql.DB, PID string) int {
	query := fmt.Sprintf("SELECT * FROM Passenger WHERE PassengerID=%s", PID)
	var passenger Passenger
	row := db.QueryRow(query)
	if err := row.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNo, &passenger.Email); err != nil {
		panic(err.Error())
	} else {
		return passenger.PassengerID
	}
}

func validatePassenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")
	//Error handling:
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if _, err := strconv.Atoi(params["id"]); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply passenger information " + "information " + "in JSON format"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(validatePassengerID(db, params["id"]))))
	}
}

func GetPassengerRecords(db *sql.DB, PID string, EML string) Passenger {
	results, err := db.Query("SELECT * FROM Passenger WHERE PassengerID=? AND Email=?", PID, EML)
	if err != nil {
		panic(err.Error())
	}
	var passenger Passenger
	for results.Next() {
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNo, &passenger.Email)
		if err != nil {
			panic(err.Error())
		}
	}
	return passenger
}

func InsertPassengerRecord(db *sql.DB, passenger Passenger) bool {
	query := fmt.Sprintf("INSERT INTO Passenger (PassengerID, FirstName, LastName, MobileNo, Email) VALUES ('%d','%s','%s','%d','%s');",
		passenger.PassengerID, passenger.FirstName, passenger.LastName, passenger.MobileNo, passenger.Email)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditPassengerRecord(db *sql.DB, PID int, FN string, LN string, MN int, EML string) bool {
	query := fmt.Sprintf("UPDATE Passenger SET FirstName='%s', LastName='%s', MobileNo=%d, Email='%s' WHERE PassengerID=%d", FN, LN, MN, EML, PID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func DeletePassenger(db *sql.DB, PID int) {
	fmt.Println("You cannot delete your account due to auditing purposes.")
}
