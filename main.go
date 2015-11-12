package main

import serial "github.com/tarm/serial"
import "log"
import "time"

func readSerial(p serial.Port) string {
	buf := make([]byte, 128)
	r, err := p.Read(buf)
	if err != nil {
		 log.Printf("Can't read from serial port\n")
		 log.Fatal(err)
	 }
	
	return string(buf[:r])
}

func writeSerial(p serial.Port, s string) (err error){
	_, err = p.Write([]byte(s + "\n"))
	if err != nil {
		log.Printf("Can't send command\n")
		log.Fatal(err)
	}
	return err
}

func sendCommand(p serial.Port, cmd string) (out string, err error) {
	err = writeSerial(p, cmd)
	_ = readSerial(p)
	out = readSerial(p)
	return out, err
}


func main() {
	//c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("Can't open serial port\n")
		log.Fatal(err)
	}


	res, err := sendCommand(*s, "print(mcu.info())")
	if err != nil {
		log.Printf("comms error")
		log.Fatal(err)
	}
	log.Printf("Cmd: print(mcu:info())\n")
	log.Printf("Result: %#v\n", res)
}
