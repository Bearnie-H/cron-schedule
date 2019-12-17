package main

import (
	"flag"
	"log"

	cronschedule "github.com/Bearnie-H/cron-schedule"
)

// Input flags.
var (
	InputTimestampFlag = flag.String("v", "* * * * *", "This flag represents the Cron Timestamp to check for validity, and to parse to an array of integer arrays representing the times to action the Cron task at.")
)

func main() {
	flag.Parse()

	if Schedule, err := cronschedule.ParseSchedule(*InputTimestampFlag); err != nil {
		log.Printf("Cron Timestamp [ %s ] is not valid - %s", *InputTimestampFlag, err)
	} else {
		log.Printf("Cron Timestamp of [ %s ] corresponds to the following times:\nMinutes: %v\nHours: %v\nDays of the Month: %v\nMonths: %v\nDays of the Week: %v\n", *InputTimestampFlag, Schedule[0], Schedule[1], Schedule[2], Schedule[3], Schedule[4])
	}
}
