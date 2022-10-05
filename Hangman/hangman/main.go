package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Tested struct {
	letters [30]string
}
type User struct {
	Wordtested string
	Health     int
	Fullword   string
	Hiddenword string
}

func randomword(s []byte) string {
	count := 0
	for _, v := range s {
		if v == '\n' {
			count++
		}
	}
	rand.Seed(time.Now().UnixNano())
	index := 0
	tabs := make([]string, count)
	for _, v := range s {
		if v == '\n' {
			index++
		}
		if v != '\n' && v != 13 && index != len(tabs) {
			tabs[index] += string(v)
		}
	}
	return tabs[rand.Intn(len(tabs))]
}
func changestring(s string, news string, nb int) string {
	nnews := ""
	for i := 0; i <= len(s)-1; i++ {
		if i != nb {
			nnews += string(news[i])
		}
		if i == nb {
			nnews += string(s[i])
		}
	}
	return nnews
}

func newstring(s string, nb int) string {
	var newmsg string
	for range s {
		newmsg += "_"
	}
	if nb != 0 {
		for i := 0; i < nb; i++ {
			indexts := rand.Intn(len(s) - 1)
			newmsg = changestring(s, newmsg, indexts)
		}
	}

	return newmsg
}
func hangman(hp int) {
	f, _ := os.Open("hangman.txt")
	scanner := bufio.NewScanner(f)
	var line int
	a := 71 - hp*8
	b := 78 - hp*8
	for scanner.Scan() {
		if line >= a && line <= b {
			fmt.Println(scanner.Text())
		}
		line++
	}
}
func verif(ucl string, word string) bool {
	for _, v := range word {
		if ucl == string(v) {
			return true
		}
	}
	return false
}
func goodletter(ucl string, word string, hword string) string {
	newhword := ""
	for i, v := range word {
		if ucl == string(v) {
			newhword += string(v)
		} else {
			newhword += string(hword[i])
		}
	}
	return newhword
}
func istested(ucl string, s [30]string) bool {
	for _, v := range s {
		if v == ucl {
			return true
		}
	}
	return false
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error syntax : (go run main.go <file.txt>)")
		os.Exit(0)
	}
	var a Tested
	index := 0
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: File not found")
		os.Exit(0)
	}
	b, _ := ioutil.ReadAll(file)
	word := randomword(b)
	nletter := len(word)/2 - 1
	hword := newstring(word, nletter)
	hp := 10
	tested := ""
	fmt.Println("Bienvenue sur le jeu hangman.")
	fmt.Println("Vous avez 10 vies, bonne chance =D")
	fmt.Println("Veuillez appuyer sur -<entrer>- ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("veuillez rentrer une valeur")
		fmt.Println("lettres testées: ", a.letters)
		ucl := scanner.Text()
		fmt.Println(hword)
		if ucl == "stop" {
			os.Create("save.txt")
			u, _ := json.Marshal(User{Wordtested: tested, Health: hp, Fullword: word, Hiddenword: hword})
			os.WriteFile("save.json", u, 5)
			fmt.Println("game stopped")
			os.Exit(0)
		}
		tested += ucl
		tested += ","
		if ucl == "" {
			fmt.Println("veuillez rentrer une valeur svp")
		}
		if ucl == word {
			fmt.Println("vous avez gagné!")
			os.Exit(0)
		}
		if len(ucl) > 1 && ucl != word {
			fmt.Print("Mauvais mot, -2 vies!")
			hp -= 2
			hangman(hp)
			if hp <= 0 {
				fmt.Println("Vous avez perdu")
				fmt.Println("Le mot était: ", word)
				os.Exit(0)
			} else {
				fmt.Println("Nombre de vies restantes: ", hp)
			}
		}
		if len(ucl) == 1 {
			if istested(ucl, a.letters) {
				fmt.Println("lettre déjà testée")
			} else if verif(ucl, word) {
				a.letters[index] = ucl
				index++
				fmt.Println("bien joué")
				hword = goodletter(ucl, word, hword)
				fmt.Println(hword)
				if hword == word || hword == word+" " {
					fmt.Println("vous avez gagné")
					os.Exit(0)
				}
			} else {
				a.letters[index] = ucl
				index++
				if hp > 1 {
					fmt.Println("Mauvaise lettre, -1 vie")
					hp--
					hangman(hp)
					fmt.Println("nombre de vies restantes", hp)
					fmt.Println("lettres testées: ", a.letters)
					fmt.Println(hword)
				} else {
					hp--
					hangman(hp)
					fmt.Println("perdu")
					fmt.Println("le mot était: ", word)
					os.Exit(0)
				}
			}
		}
	}
}
