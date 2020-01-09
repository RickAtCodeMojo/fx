package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//LadderPair is a value provided for a trade volume
type LadderPair struct {
	Volume float64
	Rate   float64
}

//Quote is the banklink value that we send over the wire
type Quote struct {
	Base    string
	Counter string
	Ladder  []LadderPair
}

//ParseFIX parses a fix message
func ParseFIX(fix string, sep string) {
	fmt.Println(fix)
	messages := strings.Split(fix, sep)
	for _, m := range messages {
		parts := strings.Split(m, "=")
		fmt.Println(parts)
		if len(parts) == 2 {
			switch parts[0] {
			case "9":
				{
					// fmt.Println("Value:", parts[1])
				}
			}
		}
	}

}
func main() {

	file, err := os.Open("../fix/fix.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	oneAndDone := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ParseFIX(scanner.Text(), "\001")
		if oneAndDone {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
