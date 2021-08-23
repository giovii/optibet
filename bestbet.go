package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Rng struct {
	Min float64
	Max float64
}
type Bets struct {
	rng Rng
	Inc float64
}
type Percentage struct {
	Min int
	Max int
	Lin int
}

type Wallet struct {
	start   float64
	current float64
}
type Player struct {
	wallet Wallet
	Stake  float64
}

type Result struct {
	Avgquote float64 `json:"avgquote",omitempty`
	Avgbid   float64 `json:"stake",omitempty`
	//	Martingale bool      `json:"martingale",omitempty`
	//Failed bool `json:"failed",omitempty`
	Times int `json:"times",omitempty`
}

type SyncWriter struct {
	m      sync.Mutex
	Writer io.Writer
}

func (w *SyncWriter) Write(b []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.Writer.Write(b)
}
func d(n float64) (a float64) {
	a, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", n), 64)

	return
}
func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func tos(x float64) (y string) {
	y = strconv.FormatFloat(x, 'f', 2, 64)
	return
}

func main() {
	bets := Bets{rng: Rng{Min: 1.1, Max: 1.5}, Inc: 0.05}
	quotes := make([]float64, 0)
	for i := bets.rng.Min; i <= bets.rng.Max; i += bets.Inc {
		quotes = append(quotes, d(i))
	}
	rand.Seed(time.Now().UnixNano())
	os.Create("nomartinngale.json")
	os.Create("nomartinngale.csv")

	for t := 5; t <= 90; t += 5 {
		for x := bets.rng.Min; x <= bets.rng.Max; x += d(bets.Inc) {

			for l := 70; l <= 99; l++ {

				player := Player{wallet: Wallet{current: 100, start: 100}, Stake: t}
				p := Percentage{Min: 0, Max: 100, Lin: l}

				bids := make([]float64, 0)
				times := 0
				for player.wallet.current < 100000 && player.wallet.current > 20 {

					esito := rand.Intn(p.Max-p.Min) + p.Min
					if esito > p.Lin {
						player.wallet.current = d(player.wallet.current - ((player.wallet.current * player.Stake) / 100))

					} else {

						player.wallet.current = d(player.wallet.current * x)
					}
					if player.Stake == 30 {
						fmt.Println((player.wallet.current))
					}
					bids = append(bids, player.wallet.current)
					times++
				}
				f, _ := os.OpenFile("nomartinngale.json", os.O_RDWR|os.O_APPEND, 0660)
				c, _ := os.OpenFile("nomartinngale.csv", os.O_RDWR|os.O_APPEND, 0660)

				if player.wallet.current > 20 {
					result, _ := json.Marshal(Result{Avgquote: d(x), Avgbid: player.Stake, Times: times})
					f.WriteString(string(result) + "\n")
					csvWriter := csv.NewWriter(c)
					csvWriter.Write([]string{tos(d(x)), tos(d(player.Stake)), strconv.Itoa(times)})
					csvWriter.Flush()

				}

				f.Close()
			}
		}
	}

	//fmt.Printf("%+v\n", Results)
	rad := false
	oldstake := 0.0

	os.Create("martinngale.json")
	os.Create("martinngale.csv")

	for t := 5.0; t <= 40; t += 5 {
		for x := bets.rng.Min; x <= bets.rng.Max; x += d(bets.Inc) {

			for l := 70; l <= 99; l++ {

				player := Player{wallet: Wallet{current: 100, start: 100}, Stake: t}
				p := Percentage{Min: 0, Max: 100, Lin: l}

				bids := make([]float64, 0)
				times := 0
				for player.wallet.current < 100000 && player.wallet.current > 20 {

					esito := rand.Intn(p.Max-p.Min) + p.Min
					if esito > p.Lin {
						player.wallet.current = d(player.wallet.current - ((player.wallet.current * player.Stake) / 100))
						rad = true
						oldstake = player.Stake
						player.Stake = player.Stake * 2
					} else {

						player.wallet.current = d(player.wallet.current * x)
						if rad == true {
							player.Stake = oldstake
						}
					}
					bids = append(bids, player.wallet.current)
					times++
				}
				f, _ := os.OpenFile("martinngale.json", os.O_RDWR|os.O_APPEND, 0660)
				c, _ := os.OpenFile("martinngale.csv", os.O_RDWR|os.O_APPEND, 0660)

				if player.wallet.current > 20 {
					result, _ := json.Marshal(Result{Avgquote: d(x), Avgbid: player.Stake, Times: times})
					f.WriteString(string(result) + "\n")
					csvWriter := csv.NewWriter(c)
					csvWriter.Write([]string{tos(d(x)), tos(d(player.Stake)), strconv.Itoa(times)})
					csvWriter.Flush()

				}
				//fmt.Printf("%+v\n", Results)
			}
		}
	}
}
