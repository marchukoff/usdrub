package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const url = "https://www.cbr-xml-daily.ru/daily_json.js"

var prompt = "Enter a value in ₽ or in $, e.g., 125 or 90$ or $90, " +
	"or %s to quit."
var transferRate float64 = 0.1

func init() {
	if runtime.GOOS == "windows" {
		prompt = fmt.Sprintf(prompt, "Ctrl+Z, Enter")
	} else { // Unix-подобная система
		prompt = fmt.Sprintf(prompt, "Ctrl+D")
	}
	cbr := getCbr()
	transferRate = getCbrRate("USD", cbr)
}

func main() {
	money := make(chan Money)
	defer close(money)
	currency := currencyConverter(money)
	convert(money, currency)
}

func currencyConverter(money chan Money) chan Money {
	currency := make(chan Money)
	go func() {
		for {
			value := <-money
			switch value.CharCode {
			case "RUB":
				currency <- Money{"USD", value.Value / transferRate}
			case "USD":
				currency <- Money{"RUB", value.Value * transferRate}
			}
		}
	}()
	return currency
}

const result = "%.2f %s = %.2f %s\n"

func convert(money chan Money, currency chan Money) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		fmt.Printf("Money to convert: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var value float64
		code := "RUB"
		if strings.HasPrefix(line, "$") || strings.HasSuffix(line, "$\n") {
			code = "USD"
			line = strings.Replace(line, "$", "", 1)
		}
		if _, err := fmt.Sscanf(line, "%f", &value); err != nil {
			fmt.Fprintln(os.Stderr, "invalid input")
			continue
		}
		money <- Money{code, value}
		foreign := <-currency
		fmt.Printf(result, value, code, foreign.Value, foreign.CharCode)
	}
	fmt.Println()
}

func getCbr() *CBR {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	cbr := new(CBR)
	err = json.Unmarshal(body, cbr)
	if err != nil {
		panic(err)
	}
	return cbr
}

func getCbrRate(code string, base *CBR) float64 {
	var res float64 = 0.0
	currency, ok := base.Valute[code]
	if ok == true {
		res = currency.Value
	}
	return res
}
