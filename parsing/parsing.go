package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	nmea "github.com/adrianmo/go-nmea"
	"github.com/argandas/serial"
)

var (
	atgps            *serial.SerialPort
	rawline          string
	gnmrc            string
	gngga            string
	strNumSatellites string
)

type paket struct {
	Timestamputc string
	Heading      string
	Latitude     string
	Longitude    string
	Speed        string
}

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
			for i := 0; i < 13; i++ {
				rawline, _ = atgps.ReadLine()
				if strings.Contains(rawline, "$GNRMC") {
					gnmrc = rawline
				} else if strings.Contains(rawline, "$GNGGA") {
					gngga = rawline
				}
			}
			//log.Printf("gnmrc:%s ", gnmrc)
			//log.Printf("gngga:%s ", gngga)

			if len(gngga) > 10 {
				m, err := nmea.Parse(gngga)
				if err != nil {
					log.Println("err GNGGA")
				} else {
					s := m.(nmea.GNGGA)
					strNumSatellites = strconv.FormatInt(s.NumSatellites, 10)
					//log.Println("NumSatellites: ", strNumSatellites)
				}
			}

			if len(gnmrc) > 10 {
				m, err := nmea.Parse(gnmrc)
				if err != nil {
					log.Println("err GNRMC")
				} else {
					s := m.(nmea.GNRMC)

					date := s.Date
					time := s.Time
					speed := s.Speed
					lat := s.Latitude
					long := s.Longitude
					Course := s.Course

					if s.Validity == "A" {
						speed *= 1.28
						//log.Printf("date:%v", date)
						//log.Printf("time:%v", time)
						//log.Printf("Course:%v", Course)
						//log.Printf("speed:%v", speed)
						//log.Printf("lat:%v", lat)
						//log.Printf("long:%v", long)

						strdate := date.String()
						strtime := time.String()
						strspeed := fmt.Sprintf("%.2f", speed)
						strlat := fmt.Sprintf("%.6f", lat)
						strlong := fmt.Sprintf("%.6f", long)
						strcorse := fmt.Sprintf("%.2f", Course)
						strspeed += " kmh"
						strcorse += " degre"
						strtime = string(strtime[0:strings.Index(strtime, ".")])
						timestamp := strdate + " " + strtime

						//log.Printf("timestamp:%s", timestamp)
						//log.Printf("strspeed:%s", strspeed)
						//log.Printf("strlat:%s", strlat)
						//log.Printf("strlong:%s", strlong)
						//log.Printf("strcorse:%s", strcorse)

						bufCD := paket{timestamp, strcorse, strlat, strlong, strspeed}
						WriteCD, err := json.MarshalIndent(bufCD, " ", " ")
						if err == nil {
							strpaket := string(WriteCD)
							strpaket += "\n"
							fmt.Printf("strpaket:%s", strpaket)
						} else {
							log.Println("err json")
						}
						fmt.Println()
					}
				}
			}
		}
	}
}
