package main

import (
	"fmt"
	"log"
	"time"

	"github.com/argandas/serial"
)

var (
	atgps *serial.SerialPort
)

func main() {

	atgps = serial.New()
	err := atgps.Open("/dev/ttyACM0", 9600, 5*time.Second)
	if err != nil {
		log.Println("PORT BUSY")
	} else {
		log.Println("SUCCESS OPEN PORT")
	}

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			var resume string

			fmt.Println("============================================================================")
			for i := 0; i < 13; i++ {
				rawline, _ := atgps.ReadLine()
				resume += rawline + "\n"
			}
			fmt.Println("%s", resume)
			fmt.Println("============================================================================")
			fmt.Println()
		}
	}
}
