package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	truc, err := os.ReadDir(".\\chapitres")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range truc {
		if f.IsDir() {
			var pa string = ".\\chapitres\\" + f.Name()
			fi, _ := os.ReadDir(pa)
			for _, g := range fi {
				fmt.Println(g.Name())
			}
		}
		fmt.Println(f.Name())
	}
}
