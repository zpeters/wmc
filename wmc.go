package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	serial "github.com/tarm/serial"
)

var Version = ""

// Functions
func cleanString(dirty string) string {
	clean := strings.Replace(dirty, "\"", "\\\"", -1)
	return clean
}
func readSerial(p serial.Port) string {
	buf := make([]byte, 1024)
	r, _ := p.Read(buf)

	return strings.TrimSpace(string(buf[:r]))
}

func writeSerial(p serial.Port, s string) (err error) {
	_, err = p.Write([]byte(s + "\n"))
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func sendCommand(p serial.Port, cmd string) (out string, err error) {
	var res string

	err = writeSerial(p, cmd)

	// drop the first line, it's always an
	// echo of the command being sent
	_ = readSerial(p)

	// read until we get a '>'
	for !strings.Contains(res, ">") {
		res = readSerial(p)
		if !strings.Contains(res, ">") {
			out = out + res + "\n"
		}
	}
	return out, err
}

func uploadFile(p serial.Port, filename string) (err error) {
	fmt.Printf("Uploading file: %s\n", filename)
	_, destFileName := filepath.Split(filename)
	fmt.Printf("Destination name: %s\n", destFileName)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Print("Writing file\n")
	_, _ = sendCommand(p, fmt.Sprintf("file.open('%s', 'w+')", destFileName))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := fmt.Sprintf("file.writeline(\"%s\")", cleanString(scanner.Text()))
		//fmt.Printf(cmd)
		_, _ = sendCommand(p, cmd)
		fmt.Printf(".")
	}
	fmt.Printf("\n")
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
	Use:   "ver",
	Short: "Get the current version",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, "print(mcu.ver())")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WifiMCU Firmware: %s", res)
	},
}

var cmdList = &cobra.Command{
	Use:   "ls",
	Short: "List all files on device",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, "for k,_ in pairs(file.list()) do print(k) end")
		//res, err := sendCommand(*s, "file.slist()")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("File list...\n")
		fmt.Printf("%s\n", res)
	},
}

var cmdPut = &cobra.Command{
	Use:   "put [filename]",
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
	Use:   "rm",
	Short: "Remove a file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("removing %s\n", args[0])
		s := openSerial()
		cmdString := fmt.Sprintf("file.remove(\"%s\")", args[0])
		_, err := sendCommand(*s, cmdString)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Display current config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WMC_SERIAL: %s\n", viper.Get("serial"))
	},
}

var cmdRead = &cobra.Command{
	Use:   "read",
	Short: "read a file",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		_, _ = sendCommand(*s, "file.open(\""+args[0]+"\")")
		_, _ = sendCommand(*s, "data=file.read()")
		_, _ = sendCommand(*s, "file.close()")
		res, _ := sendCommand(*s, "print(data)")

		fmt.Printf("%s\n", res)
	},
}

var cmdCommand = &cobra.Command{
	Use:   "cmd",
	Short: "Run an arbitrary command, needs to be in double-quotes",
	Run: func(cmd *cobra.Command, args []string) {
		s := openSerial()
		res, err := sendCommand(*s, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", res)
	},
}

func main() {
	// setup our settings
	viper.SetEnvPrefix("wmc")
	viper.BindEnv("serial")
	viper.SetDefault("serial", "/dev/ttyUSB0")

	var rootCmd = &cobra.Command{
		Use:   "wmc",
		Short: "wmc is a cross-platform command line utility for the WifiMCU Platform",
		Long: `wmc: ` + Version + `
wmc is a cross-platform command line utility for the WifiMCU Platform
GPL Licensed source code, downloads and issue tracker at
https://github.com/zpeters/wmc`,
	}
	rootCmd.AddCommand(cmdWmcVersion, cmdList, cmdPut, cmdRm, cmdConfig, cmdCommand, cmdRead)
	rootCmd.Execute()

}
