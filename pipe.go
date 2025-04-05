package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	bufferSize    = 10
	purgeInterval = 5 * time.Second
)

type RingIntBuffer struct {
	array []int
	pos   int
	size  int
	m     sync.Mutex
}

func NewRingIntBuffer(size int) *RingIntBuffer {
	return &RingIntBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

func (r *RingIntBuffer) Push(el int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}

}

func (r *RingIntBuffer) Get() []int {
	if len(r.array) < 1 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.array[:r.pos+1]
	r.pos = -1
	return output
}

func read(inChannel chan<- int, done chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var data string
	for scanner.Scan() {
		data = scanner.Text()
		if strings.EqualFold(data, "exit") {
			fmt.Println("Выход из программы")
			close(done)
			return
		}
		i, err := strconv.Atoi(data)
		if err != nil {
			fmt.Println("Вводите только целые числа")
			continue
		}

		inChannel <- i
	}
}

func NegativeFilter(inChannel <-chan int, outChannel chan<- int, done chan bool) {
	for {
		select {
		case data := <-inChannel:
			if data > 0 {
				outChannel <- data
			}
		case <-done:
			return
		}
	}
}

func nonThreeDividedFilter(inChannel <-chan int, outChannel chan<- int, done chan bool) {
	for {
		select {
		case data := <-inChannel:
			if data%3 == 0 {
				outChannel <- data
			}
		case <-done:
			return
		}
	}
}

func bufferFunc(inChannel <-chan int, outChannel chan<- int, done chan bool, size int, interval time.Duration) {
	buffer := NewRingIntBuffer(size)
	for {
		select {
		case data := <-inChannel:
			buffer.Push(data)
		case <-time.After(interval):
			bufferData := buffer.Get()
			for _, data := range bufferData {
				outChannel <- data
			}
		case <-done:
			return
		}
	}
}

func main() {
	fmt.Println("Введите целые числа, после каждого нажимайте Enter, для выхода из программы введите exit")
	input := make(chan int)
	done := make(chan bool)
	go read(input, done)

	negativeChannel := make(chan int)
	go NegativeFilter(input, negativeChannel, done)

	nonThreeDividedChannel := make(chan int)
	go nonThreeDividedFilter(negativeChannel, nonThreeDividedChannel, done)

	bufferedChannel := make(chan int)
	go bufferFunc(nonThreeDividedChannel, bufferedChannel, done, bufferSize, purgeInterval)

	for {
		select {
		case data := <-bufferedChannel:
			fmt.Println("Обработанные данные:", data)
		case <-done:
			return
		}
	}

}
