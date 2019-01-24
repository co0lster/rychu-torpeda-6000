package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type NewsAggPage struct {
	IsWinning string
	Advantage string
	Votes     int64
}

type Data struct {
	Data     []Datum `json:"data"`
	Results  int64   `json:"results"`
	NextPage bool    `json:"nextPage"`
	PrevPage bool    `json:"prevPage"`
	//	Refs     Refs          `json:"refs"`
	Debug []interface{} `json:"debug"`
}

type Datum struct {
	Name       string `json:"name"`
	Position   int64  `json:"position"`
	VotesCount int64  `json:"votesCount"`
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://plebiscite.ppapi.pl/plebiscite/v1/plebiscites/1002867/accepted-top-voted?groupId=25097&page=1&limit=500", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://dzienniklodzki.pl/p/kandydaci/czlowiek-roku-2018%2C1002867/?groupId=25097&fbclid=IwAR3W9HzhRyw8HoS99RClU0d0J8HFNchkUUKWVSne2obQcYzeKgu3spO4LPo")
	req.Header.Set("Origin", "https://dzienniklodzki.pl")
	req.Header.Set("JWT-Access-Token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6InZlcjEifQ.eyJ1c2VybmFtZSI6Im5vYm9keSIsImtpZCI6InZlcjEiLCJpYXQiOjE1NDgyNjc1NTksIm5iZiI6MTU0ODI2NzU1OSwiZXhwIjoxNTQ4MzAzNTU5LCJqdGkiOjEsInNjb3BlIjpbInBsZWJpc2NpdGU6YW5vbmltIl0sInBlcnNvbiI6eyJlbWFpbCI6bnVsbCwiaWQiOm51bGwsImZhY2Vib29rSWQiOjAsInBob25lIjpudWxsLCJrbm93blBob25lIjpmYWxzZSwibmFtZSI6IiAiLCJzZXNzaW9uSWQiOiJzc28tMTNiMzM1MDkzMjdhZmQ2NzAxNzZmM2JkYzNlNTY5NGMuMjZlMWEwYjgifSwibG9naW5VcmwiOiJodHRwczpcL1wvZHppZW5uaWtsb2R6a2kucGxcL2xvZ293YW5pZSIsImZhY2Vib29rTG9naW5VcmwiOiJodHRwczpcL1wvZHppZW5uaWtsb2R6a2kucGxcL2xvZ293YW5pZVwvZmFjZWJvb2tcLyIsIm5vdEZvdW5kVXJsIjoiaHR0cHM6XC9cL2R6aWVubmlrbG9kemtpLnBsXC9ibGVkbnlhZHJlc1wvIiwic2l0ZVVybCI6Imh0dHBzOlwvXC9kemllbm5pa2xvZHpraS5wbCIsImNvaW5QdXJjaGFzZVVybCI6Imh0dHBzOlwvXC9wbHVzLmR6aWVubmlrbG9kemtpLnBsXC93eWt1cC1kb3N0ZXBcL3BsYXRub3NjIiwidXNlcldhbGxldFVybCI6Imh0dHBzOlwvXC91c2x1Z2kuZ3JhdGthLnBsXC9yZXN0XC91enl0a293bmlrXC9cL3BvcnRmZWwiLCJob21lVXJsIjoiaHR0cHM6XC9cL2R6aWVubmlrbG9kemtpLnBsIiwic3ViIjpudWxsfQ.XPVOwWbAlTuWFdrV8CTSM8aukdMaeieq32Uw_nQwldB3b4no23SpyILgNf4bO6Crk5lzc8y-GqjMevMPFceQYenLaKpZC41E4HhJ4ZD_h-LDgvkiU5R58Q_N0Kn6aGAfMfZ-56-AD9biUK2LpAPsU9xlUUI4Cj87HTYO_pi-5Ug")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	req.Header.Set("X-Request-Id", "619527C717C2691A")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	// bodyText, err := ioutil.ReadFile("./rysz.json")
	print(resp)
	if err != nil {
		log.Fatal(err)
	}

	var adventage, isWinning string
	var votes int64
	var data Data
	var Ryszard, Opponent Datum
	json.Unmarshal(bodyText, &data)
	if data.Data[0].Name == "Ryszard Gawroński" {
		Ryszard = data.Data[0]
		Opponent = data.Data[1]
		isWinning = "TAK"
	} else {
		isWinning = "NIE"
		Opponent = data.Data[0]
		for _, candidate := range data.Data {

			if candidate.Name == "Ryszard Gawroński" {
				Ryszard = candidate
			}
		}

	}

	diff := int(Ryszard.VotesCount - Opponent.VotesCount)
	if diff > 0 {
		adventage = "i ma " + strconv.Itoa(diff) + " głosów przewagi"
	} else {
		adventage = "Brakuje mu " + strconv.Itoa(diff*-1) + " głosów do wygranej"
	}

	votes = Ryszard.VotesCount

	p := NewsAggPage{isWinning, adventage, votes}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", newsAggHandler)
	http.ListenAndServe(GetPort(), nil)
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
