package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MyStruct struct {
	name string
}

type Buffer struct {
	structsByWeek map[int][]*MyStruct
	structs       map[*MyStruct]bool
}

func (b *Buffer) Add(element *MyStruct) {
	b.structsByWeek[1] = append(b.structsByWeek[1], element)
	b.structs[element] = true
}

func main() {

	s := []int{1, 2, 3, 4, 5}
	s = getElementsFromEnd(s, 2, 2)
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}

}
func getElementsFromEnd(slice []int, n int, offset int) []int {
	// Проверяем, что offset не выходит за границы слайса

	startEl := len(slice) - (n * offset)
	if startEl < 0 {
		return nil
	}
	endEl := startEl - n
	if endEl < 0 {
		endEl = 0
	}

	return slice[endEl:startEl]
}

func sleepyGopher(id int, c chan int) { // Объявляет канал как аргумент
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Second)
	fmt.Println("... ", id, " snore ...")
	c <- id // Отправляет значение обратно к main
}
