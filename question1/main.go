package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Train struct {
	TrainName      string           `json:"trainName"`
	TrainNumber    string           `json:"trainNumber"`
	DepartureTime  Time             `json:"departureTime"`
	SeatsAvailable SeatAvailability `json:"seatsAvailable"`
	Price          TrainPrice       `json:"price"`
	DelayedBy      int              `json:"delayedBy"`
}

type Time struct {
	Hours int `json:"Hours"`

	Minutes int `json:"Minutes"`
	Seconds int `json:"Seconds"`
}

type SeatAvailability struct {
	Sleeper int `json:"sleeper"`
	AC      int `json:"AC"`
}

type TrainPrice struct {
	Sleeper float64 `json:"sleeper"`
	AC      float64 `json:"AC"`
}

func main() {
	router := gin.Default()

	router.GET("/trains", func(c *gin.Context) {
		trainsList, err := GETTRAINS()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, trainsList)
	})

	router.Run(":8000")
}

func GETTRAINS() ([]Train, error) {
	apiURL := "http://104.211.219.98/train/trains"
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODczMjgwODksImNvbXBhbnlOYW1lIjoiVmFtc2kgQ2VudHJhbCIsImNsaWVudElEIjoiMjVmZTAwMjMtY2U4Ny00MmFjLTk2MjMtOTU3Y2U3OTdkYjRiIiwib3duZXJOYW1lIjoiIiwib3duZXJFbWFpbCI6IiIsInJvbGxObyI6IjIwNDgxYTEyODYifQ.7gBXV_RBmnnBggR6LD_UhqGsrgNW29n-yQAyGAelgcc" // Replace with your actual bearer token

	cli := &http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, error := cli.Do(request)
	if error != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trains []Train
	json.NewDecoder(resp.Body).Decode(&trains)

	now := time.Now()
	endTime := now.Add(12 * time.Hour)
	FINAL := make([]Train, 0)

	for _, trainName := range trains {
		departureTime := time.Date(now.Year(), now.Month(), now.Day(), trainName.DepartureTime.Hours, trainName.DepartureTime.Minutes, trainName.DepartureTime.Seconds, 0, now.Location())
		if departureTime.After(now.Add(30*time.Minute)) && departureTime.Before(endTime) {
			FINAL = append(FINAL, trainName)
		}
	}

	return FINAL, nil
}
