package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	// Create kqueue instance
	kq, err := unix.Kqueue()
	if err != nil {
		fmt.Println("Error creating kqueue:", err)
		os.Exit(1)
	}
	defer unix.Close(kq)

	// Start a TCP listener
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Get file descriptor of the listener
	listenerFD, err := listener.(*net.TCPListener).File()
	if err != nil {
		fmt.Println("Error getting file descriptor:", err)
		os.Exit(1)
	}

	// Configure kqueue event
	event := unix.Kevent_t{
		Ident:  uint64(listenerFD.Fd()),
		Filter: unix.EVFILT_READ, // Monitor for read events
		Flags:  unix.EV_ADD | unix.EV_ENABLE,
	}

	// Register event with kqueue
	events := []unix.Kevent_t{event}
	_, err = unix.Kevent(kq, events, nil, nil)
	if err != nil {
		fmt.Println("Error registering event:", err)
		os.Exit(1)
	}

	fmt.Println("Server listening on port 8080 (using kqueue)")

	// Event loop
	for {
		// Wait for events
		kevents := make([]unix.Kevent_t, 10)
		n, err := unix.Kevent(kq, nil, kevents, nil)
		if err != nil {
			fmt.Println("Kevent wait error:", err)
			continue
		}

		// Handle events
		for i := 0; i < n; i++ {
			if kevents[i].Ident == uint64(listenerFD.Fd()) {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Accept error:", err)
					continue
				}
				buf := make([]byte, 1024)
				conn.Read(buf) // Read a byte to prevent blocking
				fmt.Println(string(buf))
				fmt.Println("New connection:", conn.RemoteAddr())
				conn.Close() // Close immediately for demo
			}
		}
	}
}
