package main

import "math/rand"

var choices = []string{"sten", "sax", "påse"}



// Kommer vi få samma random varje gång vi startar servern på nytt?

func getComputerChoice() string {
	i := rand.Intn(len(choices))
	return choices[i]
}

func getResults(playerChoice string, computerChoice string) string {

	if playerChoice == computerChoice {
		return "oavgjort"
	}

	if playerChoice == "sten" && computerChoice == "sax" {
		return "vinst"
	}
	
	if playerChoice == "sax" && computerChoice == "påse" {
		return "vinst"
	}

	if playerChoice == "påse" && computerChoice == "sten"{
		return "vinst"
	}


	return "förlust"

}
