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
	"github.com/spf13/viper"
)

var Version = ""

// Functions
func readSerial(p serial.Port) string {
	buf := make([]byte, 1024)
	r, _ := p.Read(buf)
	
	return strings.TrimSpace(string(buf[:r]))
}

func writeSerial(p serial.Port, s string) (err error){
	_, err = p.Write([]byte(s + "\n"))
	if err != nil {
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
		_, _ = sendCommand(p, fmt.Sprintf("file.write(\"%s\")", scanner.Text()))
		fmt.Sprintf("file.write(\"%s\")", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	_, _ = sendCommand(p, "file.close()")

	return err
}

func flushSerial(p serial.Port) {
	res := "flush"
	// read until we get no response
	for res != "" {
		res = readSerial(p)
	}	
}

func openSerial() *serial.Port {
	//c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	c := &serial.Config{Name: viper.GetString("serial"), Baud: 115200, ReadTimeout: time.Second * 1}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	flushSerial(*s)
	
	return s	
}


// Commands
var cmdWmcVersion = &cobra.Command{
	Use: "ver",
	Short: "Get the current version",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, "print(mcu.ver())")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Result: %s", res)
	},
}

var cmdList = &cobra.Command{
	Use: "ls",
	Short: "List all files on device",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, "file.slist()")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("File list...\n")
		fmt.Printf("%s\n", res)
	},
}

var cmdPut = &cobra.Command{
	Use: "put [filename]",
	Short: "Send a file to the device",
	Run: func(cmd *cobra.Command, args []string) {
		p := openSerial()
		err := uploadFile(*p, args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdRm = &cobra.Command{
	Use: "rm",
	Short: "Remove a file",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("removing %s\n", args[0])
		s := openSerial()
		cmdString := fmt.Sprintf("file.remove(\"%s\")", args[0])
		_, err := sendCommand(*s, cmdString)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdConfig = &cobra.Command{
	Use: "config",
	Short: "Display current config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WMC_SERIAL: %s\n", viper.Get("serial"))
	},
}

var cmdVersion = &cobra.Command{
	Use: "version",
	Short: "Display wmc version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %#v\n", Version)
	},
}


func main() {
	fmt.Printf("Version: %s\n", Version)
	
	// setup our settings
	viper.SetEnvPrefix("wmc")
	viper.BindEnv("serial")
	viper.SetDefault("serial", "/dev/ttyUSB0")
	
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(cmdWmcVersion, cmdList, cmdPut, cmdRm, cmdConfig, cmdVersion)
	rootCmd.Execute()

}