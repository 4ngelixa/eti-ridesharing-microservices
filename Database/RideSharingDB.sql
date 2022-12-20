-- Setting up MySQL - Create an account named 'user' with the password 'password'
/* CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost' */

-- Creating database and using it:
CREATE DATABASE RideSharingDB;
USE RideSharingDB;

-- To prevent conflicts:
DROP TABLE IF EXISTS Driver;
DROP TABLE IF EXISTS Passenger;
DROP TABLE IF EXISTS Trip;

-- Create Tables for Driver, Passenger, and Trip
CREATE TABLE Driver
(
	DriverID VARCHAR(5) NOT NULL PRIMARY KEY,
	FirstName VARCHAR(50),
	LastName VARCHAR(50),
	MobileNo VARCHAR(8),
	Email VARCHAR(50),
	IcNo VARCHAR(9) NOT NULL,
	LicenseNo VARCHAR(6) NOT NULL
);

CREATE TABLE Passenger
(
	PassengerID VARCHAR(5) NOT NULL PRIMARY KEY, 
	FirstName VARCHAR(50),
	LastName VARCHAR(50),
	MobileNo VARCHAR(8),
	Email VARCHAR(50)
);

CREATE TABLE Trip (
	TripID VARCHAR(5) NOT NULL PRIMARY KEY, 
    Pickup VARCHAR(10) NOT NULL, 
    Dropoff VARCHAR(10) NOT NULL, 
    DriverID VARCHAR(5), 
    PassengerID VARCHAR(5), 
    TripStatus VARCHAR(15) NOT NULL,
    FOREIGN KEY (DriverID) REFERENCES Driver(DriverID),
    FOREIGN KEY (PassengerID) REFERENCES Passenger(PassengerID)
);

-- Inserting values into Driver, Passenger, and Trip tables
INSERT INTO Driver
VALUES
    (1, 'Angelica', 'Sim', 81289617, 'sima@gmail.com', 'S1345672F', '112234'),
    (2, 'Quincy', 'Lee', 97753908, 'quincee@gmail.com', 'S1234532E', '443221'),
    (3, 'Hazel', 'Tay', 94521775, 'htay@gmail.com', 'S1457591A', '567765'),
    (4, 'Jeff', 'Oon', 88110345, 'oonjeff@gmail.com', 'S1087613D', '980912');
/* SELECT * FROM Driver; */

INSERT INTO Passenger
VALUES
    (1, 'Joey', 'Low', 89086258, 'jolo@gmail.com'),
    (2, 'Abel', 'Khoo', 85333987, 'akhoo@gmail.com'),
    (3, 'Jax', 'Tan', 98551335, 'tjax@hotmail.com'),
    (4, 'Yin', 'Lee', 98367710, 'leeyin@outlook.com');
/* SELECT * FROM Passenger; */

INSERT INTO Trip
VALUES
    (1, 123456, 218031, 1, 4, 'Finished'),
    (2, 348745, 432136, 2, 3, 'Processing'),
    (3, 335597, 123456, 3, 2, 'In Progress'),
    (4, 200213, 623310, 4, 1, 'Processing');
/* SELECT * FROM Trip; */