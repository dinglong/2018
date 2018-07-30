package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"
)

type Event struct {
	XMLName xml.Name `xml:"Event"`
	System  struct {
		EventID  string `xml:"EventID"`
		Level    string `xml:"Level"`
		Computer string `xml:"Computer"`
	}
	EventDatas []struct {
		Name string `xml:"Name,attr"`
		Data string `xml:",chardata"`
	} `xml:"EventData>Data"`
}

type Day struct {
	Start *time.Time
	Stop  *time.Time
}

const timeFormat = "2006/01/02 15:04"

func main() {
	f, _ := os.Open("out.log")
	defer f.Close()

	days := make(map[string]*Day)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		process(line, days)
	}

	o, err := os.Create("result.csv")
	if err != nil {
		log.Fatalf("open result.csv error, %v\n", err)
	}
	defer o.Close()

	for _, d := range days {
		if d.Start != nil && d.Stop != nil {
			log.Printf("%s - %s - %s\n", d.Start.Format(time.RFC3339), d.Stop.Format(time.RFC3339), d.Start.Weekday().String())
			fmt.Fprintf(o, "%s,%s,%s\n", d.Start.Format(timeFormat), d.Stop.Format(timeFormat), d.Start.Weekday().String())
		}
	}
}

func process(data string, days map[string]*Day) {
	var e Event
	if err := xml.Unmarshal([]byte(data), &e); err != nil {
		log.Fatal(err)
	}

	for _, d := range e.EventDatas {

		if d.Name == "StartTime" || d.Name == "StopTime" {
			t, err := time.Parse(time.RFC3339Nano, d.Data)
			if err != nil {
				continue
			}

			t = t.Local()

			key := t.Format("2006-01-02")
			day, exist := days[key]
			if !exist {
				day = &Day{}
				days[key] = day
			}

			if d.Name == "StartTime" {
				day.Start = &t
			} else {
				day.Stop = &t
			}
		}
	}
}
