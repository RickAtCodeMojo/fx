//
//  Pubsub envelope publisher
//

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	json "encoding/json"

	zmq "github.com/pebbe/zmq4"
)

//PairType is for enumerating types of PaIRS
type PairType int
type id int64

const (
	//USD is any pair with USD as the base
	USD PairType = iota
	//Euro is any pair with Euro as the base
	Euro
	//Cross is non Euro, now USD
	Cross
)

//Currency is a description of a currency that can be included in a pair
type Currency struct {
	Key    int    `json:"Key"`
	Name   string `json:"Name"`
	ISO    string `json:"ISO"`
	Symbol string `json:"Symbol"`
}

//Rate is the bid and ask for a pair at a given time
type Rate struct {
	Bid     float64
	BidPips uint64
	Ask     float64
	AskPips uint64
	At      time.Time
}

//Pair is a combination of currencies that are traded at a rate
type Pair struct {
	Base       Currency
	Counter    Currency
	BankRate   Rate
	OurRate    Rate
	ClientRate Rate
}

//Trade is a transaction where a currency trades hands for a price
type Trade struct {
	Trader  id
	Base    string
	Counter string
	Bid     float64
	Ask     float64
	At      time.Time
}

//Trader is an entity who trades on our system
type Trader struct {
	ID            string `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	IP            string `json:"ip_address"`
	Cell          string `json:"cell"`
	BusinessPhone string `json:"business"`
	BankAccount   string `json:"bank"`
	CreditCard    string `json:"credit_card"`
}

//Set creates fills in a new pair with the base and counter currency
func (p *Pair) Set(BaseCurrency Currency, CounterCurrency Currency) {
	p.Base = BaseCurrency
	p.Counter = CounterCurrency

}

func fractionalPart(rate float64) uint64 {

	var dollars int64
	var cents uint64
	s := fmt.Sprintf("%f", rate)
	if _, err := fmt.Sscanf(s, "%d.%d", &dollars, &cents); err != nil {
		panic(err)
	}
	return cents
}

//GetOurRate refreshes the rate for the pair that we get
func (p *Pair) GetOurRate() {
	ours := Rate{Bid: 0.89567, BidPips: 567, Ask: 0.89564, AskPips: 564, At: time.Now()}
	p.OurRate = ours
	bidRate := ours.Bid * 1.05
	askRate := ours.Ask * 1.05
	p.ClientRate = Rate{Bid: bidRate, BidPips: fractionalPart(bidRate), Ask: askRate, AskPips: fractionalPart(askRate), At: ours.At}
}

func main() {
	traders, err := getTraders("../customer/customer.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(traders[0])
	pairs, err := makePairs()
	if err != nil {
		fmt.Println(err)
		return
	}
	context, _ := zmq.NewContext()
	defer context.Term()

	publisher, _ := context.NewSocket(zmq.PUB)
	defer publisher.Close()
	publisher.Bind("tcp://*:5563")

	for {
		for _, pr := range pairs {
			publisher.Send(pr.Base.ISO, zmq.SNDMORE)
			rate := fmt.Sprintf("%s/%s Bid:%f, Ask:%f", pr.Base.ISO, pr.Counter.ISO, pr.ClientRate.Bid, pr.ClientRate.Ask)
			publisher.Send(rate, 0)
		}
	}
}
func getTraders(filePath string) ([]Trader, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var traders []Trader
	json.Unmarshal(raw, &traders)
	return traders, nil
}

func getCurrencies(filePath string) ([]Currency, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var currencies []Currency
	json.Unmarshal(raw, &currencies)
	return currencies, nil
}

func makePairs() ([]Pair, error) {
	var currencies, err = getCurrencies("../currency/currency.json")
	if err != nil {
		return nil, err
	}
	if len(currencies) == 0 {
		return nil, errors.New("no currencies could be found")
	}
	usd := currencies[0]
	euro := currencies[1]
	pairs := make([]Pair, 180)
	for i, c := range currencies {
		if i > 0 {
			var p Pair
			p.Set(usd, c)
			p.GetOurRate()
			pairs = append(pairs, p)
		}
	}
	for i, c := range currencies {
		if i != 1 {
			var p Pair
			p.Set(euro, c)
			p.GetOurRate()
			pairs = append(pairs, p)

		}
	}
	// for _, p := range pairs {
	// 	fmt.Println(p)
	// }
	return pairs, nil
}

/*
	[]Currency{
		{Rank: 1, Name: "United States dollar", ISO: "USD", Symbol: "US$", Frequency: 87.6},
		{Rank: 2, Name: "Euro", ISO: "EUR", Symbol: "€", Frequency: 31.4},
		{Rank: 3, Name: "Japanese yen", ISO: "JPY", Symbol: "¥", Frequency: 21.6},
		{Rank: 4, Name: "Pound sterling", ISO: "GBP", Symbol: "£", Frequency: 2.8},
		{Rank: 5, Name: "Australian dollar", ISO: "AUD", Symbol: "A$", Frequency: 6.9},
		{Rank: 6, Name: "Canadian dollar", ISO: "CAD", Symbol: "C$", Frequency: 5.1},
		{Rank: 7, Name: "Swiss franc", ISO: "CHF", Symbol: "Fr", Frequency: 4.8},
		{Rank: 8, Name: "Renminbi", ISO: "CNY", Symbol: "元", Frequency: 4.0},
		{Rank: 9, Name: "Swedish krona", ISO: "SEK", Symbol: "kr", Frequency: 2.2},
		{Rank: 10, Name: "New Zealand dollar", ISO: "NZD", Symbol: "NZ$", Frequency: 2.1},
		{Rank: 11, Name: "Mexican peso", ISO: "MXN", Symbol: "$", Frequency: 1.9},
		{Rank: 12, Name: "Singapore dollar", ISO: "SGD", Symbol: "S$", Frequency: 1.8},
		{Rank: 13, Name: "Hong Kong dollar", ISO: "HKD", Symbol: "HK$", Frequency: 1.7},
		{Rank: 14, Name: "Norwegian krone", ISO: "NOK", Symbol: "kr", Frequency: 1.7},
		{Rank: 15, Name: "South Korean won", ISO: "KRW", Symbol: "₩", Frequency: 1.7},
		{Rank: 16, Name: "Turkish lira", ISO: "TRY", Symbol: "₺", Frequency: 1.4},
		{Rank: 17, Name: "Russian ruble", ISO: "RUB", Symbol: "₽", Frequency: 1.1},
		{Rank: 18, Name: "Indian rupee", ISO: "INR", Symbol: "₹", Frequency: 1.1},
		{Rank: 19, Name: "Brazilian real", ISO: "BRL", Symbol: "R$", Frequency: 1.0},
		{Rank: 20, Name: "South African rand", ISO: "ZAR", Symbol: "R", Frequency: 1.0},
		{Rank: 21, Name: "Danish krone", ISO: "DKK", Symbol: "kr", Frequency: 0.8},
		{Rank: 22, Name: "Polish złoty", ISO: "PLN ", Symbol: "zł", Frequency: 0.7},
		{Rank: 23, Name: "New Taiwan dollar", ISO: "TWD", Symbol: "NT$", Frequency: 0.6},
		{Rank: 24, Name: "Thai baht", ISO: "THB", Symbol: "฿", Frequency: 0.4},
		{Rank: 25, Name: "Malaysian ringgit", ISO: "MYR", Symbol: "RM", Frequency: 0.4},
		{Rank: 26, Name: "Hungarian forint", ISO: "HUF", Symbol: "Ft", Frequency: 0.3},
		{Rank: 27, Name: "Saudi riyal", ISO: "AR", Symbol: "﷼", Frequency: 0.3},
		{Rank: 28, Name: "Czech koruna", ISO: "CZK", Symbol: "Kč", Frequency: 0.3},
		{Rank: 29, Name: "Israeli shekel", ISO: "ILS", Symbol: "₪", Frequency: 0.3},
		{Rank: 30, Name: "Chilean peso", ISO: "CLP", Symbol: "CLP$", Frequency: 0.2},
		{Rank: 31, Name: "Indonesian rupiah", ISO: "IDR", Symbol: "Rp", Frequency: 0.2},
		{Rank: 32, Name: "Colombian peso", ISO: "COP", Symbol: "COL$", Frequency: 0.2},
		{Rank: 33, Name: "Philippine peso", ISO: "PHP", Symbol: "₱", Frequency: 0.1},
		{Rank: 34, Name: "Romanian leu", ISO: "RON ", Symbol: "L", Frequency: 0.1},
		{Rank: 35, Name: "Peruvian sol", ISO: "PEN", Symbol: "S/", Frequency: 0.1},
	}

*/
