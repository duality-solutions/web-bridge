package main

// one goroutine is the main
// goroutine that comes by default
import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type message struct {
	counter int
}

type messages struct {
	mapMessages map[string]chan *message
	mapMess     map[string]*message
	*sync.RWMutex
}

func (m *messages) updateMessageMap(id string, message *message) {
	m.Lock()
	defer m.Unlock()
	m.mapMessages[id] <- message
}

func (m *messages) updateMessMap(id string, message *message) {
	m.Lock()
	defer m.Unlock()
	m.mapMess[id] = message
}

func newMessages() *messages {
	m := messages{}
	m.mapMessages = make(map[string]chan *message)
	m.mapMess = make(map[string]*message)
	m.RWMutex = new(sync.RWMutex)
	return &m
}

var wgIns sync.WaitGroup

func main() {
	wgIns.Add(6)
	var msgs = newMessages()
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	msg := message{
		counter: r1.Intn(100),
	}
	for i := 0; i < 3; i++ {
		// goroutines are made
		go func() {
			for j := 0; j < 3; j++ {
				// shared variable execution
				msgs.mapMessages["1"] <- &msg
				msgs.mapMess["1"] = &msg //no race
			}
		}()
		wgIns.Done()
	}
	for i := 0; i < 3; i++ {
		// goroutines are made
		go func() {
			for j := 0; j < 3; j++ {
				// shared variable execution
				msgs.mapMess["1"] = &msg //race
			}
		}()
		wgIns.Done()
	}
	fmt.Println("The number of goroutines before wait =", runtime.NumGoroutine())
	wgIns.Wait()
	fmt.Println("The number of goroutines after wait = ", runtime.NumGoroutine())
}
