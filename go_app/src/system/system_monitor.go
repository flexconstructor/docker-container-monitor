package system

import (
	"log"
	"net"
	"os"
	"time"
)

// Establishes connection to UNIX socket for write media-server container system
// info.
type SystemMonitor struct {
	socket_file   string             // UNIX socket file path.
	close_channel chan bool          // Channel for close signal.
	info_factory  *SystemInfoFactory // System info factory.
}

// Returns new instance of system monitor.
func NewSystemMonitor(socket_file_name string) *SystemMonitor {
	return &SystemMonitor{
		socket_file:   socket_file_name,
		close_channel: make(chan bool),
		info_factory:  NewSystemInfoFactory("nginx"),
	}
}

// Run system monitor.
func (m *SystemMonitor) Run() {
	defer os.Exit(1)
	connection, err := net.Listen("unix", m.socket_file)
	if err != nil {
		log.Printf("listen error %v", err)
		return
	}
	log.Printf("connection %v", connection.Addr().String())
	defer connection.Close()
	go m.listenConnection(connection)

	if err != nil {
		log.Printf("Dial error %v", err)
		return
	}
	for {
		select {
		case <-m.close_channel:
			log.Println("close listener")
			return

		}
	}
}

// Listens established connection and starts timer for send information for
// client.
func (m *SystemMonitor) listenConnection(listener net.Listener) {

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("can not acept connection %v", err)
			m.Stop()
			return
		}
		log.Printf("ConnectionAcepted %v", connection.LocalAddr().String())
		for {
			select {
			case <-m.close_channel:
				log.Println("close connection")
				connection.Close()
				return
			case <-time.After(2 * time.Second):
				_, err := connection.Write(m.info_factory.GetSystemInfo())
				if err != nil {
					log.Printf("Can not write message to socket: %v", err)
					m.Stop()
					return
				}
			}
		}
	}
}

// Writes command to close channel for close connection and listener.
func (m *SystemMonitor) Stop() {
	log.Println("Stop Monitor")
	m.close_channel <- true
}
