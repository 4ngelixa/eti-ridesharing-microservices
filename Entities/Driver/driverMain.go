package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Map for Driver Records
type Driver struct {
	DriverID  string `json:"DriverID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MobileNo  int    `json:"MobileNo"`
	Email     string `json:"Email"`
	IcNo      string `json:"IcNo"`
	LicenseNo int    `json:"LicenseNo"`
}

var drivers map[string]driverInfo

type driverInfo struct {
	Title string `json:"Driver"`
}

var (
	db  *sql.DB
	err error
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/driverMain", allDrivers)
	router.HandleFunc("/api/v1/validateDriverRecord/{icno}", validateDriver)
	router.HandleFunc("/api/v1/GetDriver/{driverid}", GetDriverID)
	router.HandleFunc("/api/v1/drivers/{driverid}/{icno}", driverMain).Methods("GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

// Function to allow Drivers to call GET, PUT, POST, and DELETE methods
func driverMain(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Connect to MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")

	//Error handling:
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database successfully opened :)")
	}

	// Retrive driver data
	if r.Method == "GET" {
		var getDrivers Driver
		reqBody, err := ioutil.ReadAll(r.Body)
		// Defer the close until the main function is done executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &getDrivers)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("%s - Error in JSON encoding", err)
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information"))
				return
			}
		}
		json.NewEncoder(w).Encode(GetDriverRecords(db, params["driverid"], params["icno"]))
		w.WriteHeader(http.StatusAccepted)
		return
	}
	// Create a new driver record
	if r.Method == "POST" {
		var newDriver Driver
		reqBody, err := ioutil.ReadAll(r.Body)
		// Defer the close until the main function is done executing
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &newDriver)
			if err != nil {
				println(string(reqBody))
				fmt.Printf("%s - Error in JSON encoding", err)
			} else {
				if newDriver.IcNo == "" { // IC No. validation - IC No. cannot be empty
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver information " + "in JSON format"))
					return
				} else { // Run validation funcs
					if !validateDriverRecord(db, newDriver.IcNo) {
						InsertDriverRecord(db, newDriver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - New Driver added"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("409 - Duplicate Driver ID"))
						return
					}
				}
			}
		}
		// Update existing driver records
	} else if r.Method == "PUT" {
		var updateDriver Driver
		reqbody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqbody, &updateDriver)

			// Name validation - First name cannot be empty.
			if updateDriver.FirstName == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " + "information " + "in JSON format"))
				return
			} else { // Run validation funcs
				if !validateDriverRecord(db, updateDriver.IcNo) {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("No driver found with: " + updateDriver.IcNo))
				} else {
					EditDriverRecord(db, updateDriver)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver updated :)"))
				}
			}
		}
	} // Driver users are unable to delete their account as per requirements.
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("404 - You cannot delete your account due to auditing purposes"))
	}
}

// Retrieve all driver records
func allDrivers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")
	//Error handling:
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database successfully opened :)")
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(GetAvailDriver(db)))
	//json.NewEncoder(w).Encode(GetDriverRecords(db, params["driverid"], params["email"]))
}

func GetDriverID(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")
	//Error handling:
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["driverid"] == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Please enter your driver ID"))
		return
	} else {
		println(params["driverid"])
		query := fmt.Sprintf("SELECT LicenseNo FROM Driver WHERE DriverID='%s'", params["driverid"])
		results, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var LicenseNo string
		for results.Next() {
			err = results.Scan(&LicenseNo)
			if err != nil {
				panic(err.Error())
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(LicenseNo))
		return
	}
}

func validateDriverIC(db *sql.DB, IC string) string {
	query := fmt.Sprintf("SELECT * FROM Driver WHERE IcNo= '%s'", IC)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.Email, &driver.IcNo, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver.DriverID
}

func validateDriver(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["email"] == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply driver IC No. " + "information " + "in JSON format"))
		return
	} else if validateDriverRecord(db, params["icno"]) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(validateDriverIC(db, params["icno"])))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

// Check if the driver has already created an account using their IC number
func validateDriverRecord(db *sql.DB, IC string) bool {
	query := fmt.Sprintf("SELECT * FROM Driver WHERE IcNo= '%s'", IC)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.Email, &driver.IcNo, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		} else if driver.IcNo == IC {
			return true
		}
	}
	return false
}

func GetDriverRecords(db *sql.DB, DID string, IC string) Driver {
	results, err := db.Query("SELECT * FROM Driver WHERE DriverID=? AND IcNo=?", DID, IC)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.Email, &driver.IcNo, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver
}

func InsertDriverRecord(db *sql.DB, driver Driver) bool {
	query := fmt.Sprintf("INSERT INTO Driver (DriverID, FirstName, LastName, MobileNo, Email, IcNo, LicenseNo) VALUES ('%s','%s','%s','%d','%s', '%s', '%d');",
		driver.DriverID, driver.FirstName, driver.LastName, driver.MobileNo, driver.Email, driver.IcNo, driver.LicenseNo)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditDriverRecord(db *sql.DB, driver Driver) bool {
	query := fmt.Sprintf("UPDATE Driver SET FirstName='%s', LastName='%s', MobileNo=%d, Email='%s', IcNo='%s', LicenseNo='%d' WHERE DriverID='%s'",
		driver.FirstName, driver.LastName, driver.MobileNo, driver.Email, driver.IcNo, driver.LicenseNo, driver.DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

// Retrieve records of Drivers who are not driving
func GetAvailDriver(db *sql.DB) string {
	query := "SELECT DriverID FROM Driver"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DriverID string
	for results.Next() {
		var ID string
		err = results.Scan(&ID)
		if err != nil {
			panic(err.Error())
		}
		DriverID += ID + ","
	}
	return DriverID
}

func DeleteDriver(db *sql.DB, DID int) {
	fmt.Println("You cannot delete your account due to auditing purposes.")
}
