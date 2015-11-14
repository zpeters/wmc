package main

import (
	"strings"
	"log"
	"fmt"
	"time"
	"os"
	"bufio"
	"path/filepath"
)

import (
	serial "github.com/tarm/serial"
	"github.com/spf13/cobra"
)


// Functions
func readSerial(p serial.Port) string {
	buf := make([]byte, 1024)
	r, err := p.Read(buf)
	if err != nil {
		 log.Printf("Can't read from serial port\n")
		 log.Fatal(err)
	 }
	
	return strings.TrimSpace(string(buf[:r]))
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

func uploadFile(p serial.Port, filename string) (err error) {
	log.Printf("Uploading file: %s\n", filename)
	_, destFileName := filepath.Split(filename)
	log.Printf("Destination name: %s\n", destFileName)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	log.Printf("Writing file...\n")
	_, _ = sendCommand(p, fmt.Sprintf("file.open('%s', 'w+')", destFileName))
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Print(scanner.Text())
		//_, _ = sendCommand(p, fmt.Sprintf("file.write(\"%s\")", scanner.Text()))
		fmt.Sprintf("file.write(\"%s\")", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	_, _ = sendCommand(p, "file.close()")

	return err
}

func openSerial() *serial.Port {
	//c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("Can't open serial port\n")
		log.Fatal(err)
	}
	return s
	
}

// Commands
var cmdVersion = &cobra.Command{
	Use: "version",
	Short: "Get the current version",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, "print(mcu.ver())")
		if err != nil {
			log.Printf("Comm error")
			log.Fatal(err)
		}
		log.Printf("Result: %s", res)
	},
}


func main() {
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(cmdVersion)
	rootCmd.Execute()

}
