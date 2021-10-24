package main

import (
	"os"
	"sync"
	"testing"
)

var connMap = &sync.Map{} // Mock in setup()

func TestWhoAmI(t *testing.T) {

}

func TestWhoElse(t *testing.T) {

}

func TestSendCommand(t *testing.T) {

}

func TestParseSendCommand(t *testing.T) {

}

func setup() {
	// Mock client connections and store in connMap perhaps using net.Pip

	// This would require implementing net.Conn interface. In interest of time, I outlined the methods I would test once I have the mocks.
}

func shutdown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
