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

// Map for Trip Records
type Trip struct {
	TripID      int    `json:"TripID"`
	PickupPC    string `json:"PickupPostalCode"`
	DropoffPC   string `json:"DropoffPostalCode"`
	DriverID    int    `json:"DriverID"`
	PassengerID int    `json:"PassengerID"`
	TripStatus  string `json:"TripStatus"`
}

// Map for Passenger Records
type Passenger struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	MobileNo    int    `json:"MobileNo"`
	Email       string `json:"EmailAddress"`
}

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

var trips map[string]tripInfo

type tripInfo struct {
	Title string `json:"Trips"`
}

func main() {
	router := mux.NewRouter()
	trips = make(map[string]tripInfo)
	router.HandleFunc("/api/v1/trips/{tripid}", trip).Methods("GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/api/v1/trips", alltrips)

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
}

func alltrips(w http.ResponseWriter, r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, v)
	}
	json.NewEncoder(w).Encode(trips)
}

// Function to allow GET, PUT, POST, and DELETE methods for Trips
func trip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Connect to MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/RideSharingDB")

	//Error handling:
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database successfully opened :)")
	}

	// Retrieve trip data
	if r.Method == "GET" {
		fmt.Println(params)
		Tripid, err := strconv.Atoi(params["tripid"])
		if err != nil {
			fmt.Println(err)
		}
		tripInfo := GetAllTripsRecord(db, Tripid)
		if err != nil {
			fmt.Printf("Error in JSON encoding. Error is %s", err)
		} else if tripInfo.TripID == 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Invalid trip!"))
			return
		} else {
			json.NewEncoder(w).Encode(GetAllTripsRecord(db, tripInfo.TripID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
		// Create a new trip
		if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			// Defer the close until the main function is done executing
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &newTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("%s - Error in JSON encoding", err)
				} else {
					if newTrip.PassengerID == 0 || newTrip.DriverID == 0 {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("422 - Please supply trip " + "information " + "in JSON format"))
						return
					}

					if validateTrips(db, newTrip.PassengerID, newTrip.DriverID) {
						InsertTripRecord(db, newTrip)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Trip added"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("409 - Ongoing trip"))
						return
					}
				}
			}
			// Update existing trips
		} else if r.Method == "PUT" {
			var updateTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			// Defer the close until the main function is done executing
			defer r.Body.Close()
			if err == nil {
				err := json.Unmarshal(reqBody, &updateTrip)
				if err != nil {
					println(string(reqBody))
					fmt.Printf("%s - Error in JSON encoding", err)
				} else {
					if validateTrips(db, updateTrip.DriverID, updateTrip.PassengerID) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("404 - No trip found"))
					} else {
						if EditTripRecord(db, updateTrip) {
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("201 - Trip is updated"))
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("401 - Invalid information"))
						}
					}
				}
			}
		}
		// Unable to delete trips as per requirements.
		w.Write([]byte("404 - You cannot delete your trip due to auditing purposes"))
	}
}

// Retrieve all trips requested
func GetAllTripsRecord(db *sql.DB, TID int) Trip {
	results, err := db.Query("SELECT * FROM Trip WHERE TripID=?", TID)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip
}

// Retrieve trips requested by a Passenger (based on their PID)
func GetAllTrips(db *sql.DB, PID int) []Trip {
	query := fmt.Sprintf("SELECT * FROM Trip WHERE PassengerID ='%d'", PID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() {
		var trip Trip
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
		trips = append(trips, trip)
	}
	return trips
}

func InsertTripRecord(db *sql.DB, trip Trip) bool {
	query := fmt.Sprintf("INSERT INTO Trip (TripID, PickupPC, DropoffPC, DriverID, PassengerID, TripStatus) VALUES ('%d','%s','%s','%d','%d','%s')",
		trip.TripID, trip.PickupPC, trip.DropoffPC, trip.DriverID, trip.PassengerID, trip.TripStatus)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditTripRecord(db *sql.DB, trip Trip) bool {
	query := fmt.Sprintf("UPDATE Trip SET PickupPC='%s', DropoffPC='%s', DriverID='%d', PassengerID='%d', TripStatus='%s' WHERE TripID='%d'",
		trip.PickupPC, trip.DropoffPC, trip.DriverID, trip.PassengerID, trip.TripStatus, trip.TripID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func DeleteTrip(db *sql.DB, TID int) {
	fmt.Println("You cannot delete your trip record due to auditing purposes.")
}

// Validation
func validateTrips(db *sql.DB, PID int, DID int) bool {
	query := fmt.Sprintf("SELECT * FROM Trip WHERE PassengerID= '%d' OR DriverID='%d'", PID, DID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip Trip
	for results.Next() {
		err = results.Scan(&trip.TripID, &trip.PickupPC, &trip.DropoffPC, &trip.DriverID, &trip.PassengerID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		} else if trip.TripStatus != "Finished" {
			return false
		}
	}
	return true
}

func GetAllDriverRecords() Driver {
	response, err := http.Get("http://localhost:5001/api/v1/GetAllDriverRecords/")
	if err != nil {
		fmt.Print(err.Error())
	}
	var driverTrip Driver
	if response.StatusCode == http.StatusAccepted {
		responseData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(responseData))
		fmt.Println(response.StatusCode)
		response.Body.Close()
		json.Unmarshal(responseData, &driverTrip)
	} else {
		fmt.Printf("404 - There are no available drivers at the moment")
		return driverTrip
	}
	return driverTrip
}

func GetDriver(DriverID string) string {
	url := "http://localhost:5001/api/v1/GetDriver/" + DriverID

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
}

// Retrieve Drivers who are available
func GetAvailDriver(db *sql.DB, DID int) int {
	results, err := db.Query("SELECT Driver.DriverID FROM Driver INNER JOIN Trip ON Driver.DriverID = Trip.DriverID WHERE Trip.TripStatus='Finished'")
	if err != nil {
		panic(err.Error())
	}
	var DriverID int
	for results.Next() {
		err = results.Scan(&DriverID)
		if err != nil {
			panic(err.Error())
		}
	}
	return DriverID
}
