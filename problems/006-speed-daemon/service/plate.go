package service

import (
	"fmt"
	"log"
	"math"
	"net"
)

// Function to handle plate message
func (ts *TicketSystem) HandlePlate() {

	for sp := range ts.spch {
		log.Println(sp)
		day := sp.timestamp / DAY_DIV

		plate_day := fmt.Sprintf("%v_%v", sp.plate, day)

		_, ok := ts.hasTicketOnDay[plate_day]
		if ok {
			continue
		}

		road_plate := fmt.Sprintf("%v_%v", sp.road, sp.plate)
		sps, ok := ts.roadPlateSpeedProps[road_plate]
		if ok {
			speedLimit := sp.speedLimit

			for _, sp2 := range sps {
				speed := getSpeed(sp, sp2)
				log.Println(speed, speedLimit)

				if speed >= float64(speedLimit)+0.5 {
					// Issue a ticket
					ticket := &Ticket{
						plate: sp.plate,
						road:  sp.road,
						speed: Limit(math.Round(speed * 100)),
					}
					if sp.timestamp < sp2.timestamp {
						ticket.mile1 = sp.mile
						ticket.timestamp1 = sp.timestamp
						ticket.mile2 = sp2.mile
						ticket.timestamp2 = sp2.timestamp
					} else {
						ticket.mile1 = sp2.mile
						ticket.timestamp1 = sp2.timestamp
						ticket.mile2 = sp.mile
						ticket.timestamp2 = sp.timestamp
					}
					log.Println("Sending Ticket")
					sendTicket(ticket, ts.dispatcher[sp.road])
					ts.hasTicketOnDay[plate_day] = 1
					break
				}
			}
			sps = append(sps, sp)

		} else {
			sps := make([]*SpeedProps, 0)
			sps = append(sps, sp)
			ts.roadPlateSpeedProps[road_plate] = sps
		}
	}

}

func sendTicket(ticket *Ticket, dispatchers chan net.Conn) {
	conn := <-dispatchers
	log.Println(conn.RemoteAddr())
	_, err := conn.Write(ticket.Encode())
	log.Println(ticket.Encode())
	if err != nil {
		conn.Close()
		return
	}
	dispatchers <- conn
}

func getSpeed(sp1, sp2 *SpeedProps) float64 {
	d := math.Abs(float64(sp1.mile) - float64(sp2.mile))
	t := math.Abs(float64(sp1.timestamp) - float64(sp2.timestamp))

	if t == 0 {
		return 0
	}
	return d / (t / 3600)
}
