package main

import "github.com/tarm/goserial"
import "log"

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("Can't open serial port\n")
		log.Fatal(err)
	}

	_, err = s.Write([]byte("print(mcu.info())\n"))
	if err != nil {
		log.Printf("Can't send command\n")
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	r, err := s.Read(buf)
	if err != nil {
		log.Printf("Can't read from serial port\r")
		log.Fatal(err)
	}
	
	log.Printf("Reading from Serial\n")
	log.Printf("%#v", string(buf[:r]))

	buf = make([]byte, 128)
	r, err = s.Read(buf)
	if err != nil {
		log.Printf("Can't read from serial port\r")
		log.Fatal(err)
	}
	
	log.Printf("Reading from Serial\n")
	log.Printf("%#v", string(buf[:r]))
}
