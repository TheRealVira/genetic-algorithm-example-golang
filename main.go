package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

const MutationRate = 5
const PopulationCount = 100
const PredatorCount = 20

type Child struct  {
	text string
	dif int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter goal-text: ")
	goalText, _ := reader.ReadString('\n')

	// Little hack since certain systems use "\n" and others "\r\n"
	goalText = strings.TrimSuffix(strings.TrimSuffix(goalText, "\n"), "\r")

	fmt.Print("\"")
	fmt.Print(goalText)
	fmt.Println("\"")

	myGeneration := calcNewGeneration(PopulationCount, goalText)

	for i := 0; i < len(myGeneration) ; i++ {
		fmt.Println(myGeneration[i])
	}

	generation:=1

	// While the set of children does not contain the wanted string:
	for !sliceContains(goalText, myGeneration){
		// Sort by fitness
		sort.Slice(myGeneration, func(i, j int) bool {
			return myGeneration[i].dif < myGeneration[j].dif
		})

		// Monitor fittest entry
		fmt.Println("Top Entry:\t", myGeneration[0])

		var predators []Child
		// Store top entries inside a slice
		for i := 0; i < PredatorCount; i++ {
			predators = append(predators, myGeneration[i])
		}

		// Breed all parents to gather new set of children
		newChildren := breedAllParents(predators, goalText)

		// Mutate new set of children
		mutateChildren(newChildren, goalText)

		// Update the current generation by the sum of new children and predators
		myGeneration = append(append(newChildren, calcNewGeneration(PopulationCount, goalText)...), predators...)

		generation++
	}

	fmt.Println(goalText)

	fmt.Println("Finished! Generation:\t",generation)
}

func calcNewGeneration(count int, goal string) []Child {
	toRet := make([]Child, count)

	for count--; count >= 0; count-- {
		temp:=calcRandomString(len(goal))
		toRet[count] = Child{temp, calcScore(temp, goal)}
	}

	return toRet
}

func calcRandomString(length int)string{
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(byte(randInt(32, 127)))
	}
	return string(bytes)
}

func calcScore(input, goal string) int{
	dif:=0
	for i := 0; i < len(input); i++ {
		if string(input[i])!=string(goal[i]) {
			dif++
		}
	}

	return dif
}

func sliceContains(entry string, children []Child) bool{
	for i := 0; i < len(children); i++ {
		if entry == children[i].text {
			return true
		}
	}
	return false
}

func breedAllParents(parents []Child, goalText string) []Child {
	var toRet []Child

	for i := 0; i < len(parents); i++ {
		for j := 0; j < len(parents); j++ {
			if parents[i]==parents[j] {
				continue
			}

			toRet = append(toRet, breedTwoParents(parents[i], parents[j], goalText)...)
		}
	}

	return toRet
}

func breedTwoParents(p1, p2 Child, goalText string) []Child {
	toRet:=[]Child{
		mixTwoStrings(p1.text, p2.text, goalText),
		mixTwoStrings(p2.text, p1.text, goalText),
	}

	return toRet
}

func mixTwoStrings(s1, s2, goalText string) Child {
	toRet:= Child{}

	for i := 0; i < len(s1); i++ {
		if rand.Intn(2)==0 {
			toRet.text+=string(s1[i])
		} else
		{
			toRet.text+=string(s2[i])
		}
	}

	toRet.dif = calcScore(toRet.text, goalText)

	return toRet
}

func mutateChildren(children []Child, goalText string){
	for i := 0; i < len(children); i++ {
		if rand.Intn(MutationRate)==1 {
			newText:=mutateString(children[i].text)
			children[i]= Child{
				text:newText,
				dif:calcScore(newText, goalText),
			}
		}
	}
}

func mutateString(input string) string{
	randPosition:=rand.Intn(len(input))

	b := make([]rune, len(input))
	for i := range b {
		if i==randPosition {
			b[i] = rune(byte(randInt(32, 127)))
		} else{
			b[i]= rune(input[i])
		}
	}
	return string(b)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}