package main

import (
	"cloudMonitor/package/cloudFrame"
	"cloudMonitor/package/comDevice"
	"cloudMonitor/package/config"
	"cloudMonitor/package/constant"
	httpClient "cloudMonitor/package/http"
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	// input arguments parse
	cfgFile := flag.String("config", "config.yaml", "configuration file in Yaml format")
	refFile := flag.String("reference", "reference.txt", "cloud reference data file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Syntax: %s [option]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "option:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if 0 == flag.NFlag() {
		flag.Usage()
		return
	}

	fmt.Print("Loading config...")
	config := new(config.ConfigData)
	err := config.LoadFromYaml(*cfgFile)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("done.")

	fmt.Print("Loading cloud reference...")
	cloudFrame.LoadConfig(config)
	err = cloudFrame.LoadReference(*refFile)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("done.")

	fmt.Print("Establishing connection to the infrared camera...")
	infraredCam := new(comDevice.SerialPortData)
	err = infraredCam.OpenPort(config.ComPort)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("connected to", config.ComPort)

	defer infraredCam.Close()

	var waitGroup sync.WaitGroup

	fmt.Print("Starting the receiver routine...")
	waitGroup.Add(1)
	go comDevice.Receiver(infraredCam, &waitGroup)
	fmt.Println("done.")

	fmt.Print("Initializing the infrared camera...")
	infraredCam.SendCommand(constant.COMCMD_STOP_OUTPUT)
	infraredCam.WaitForString(constant.COMCMD_OK, true)
	infraredCam.SendCommand(constant.COMCMD_TURN_OFF_LED)
	infraredCam.WaitForString(constant.COMCMD_OK, true)
	fmt.Println("done.")

	dataChannel := make(chan []float32, 2)

	fmt.Print("Starting the frame fetcher routine...")
	waitGroup.Add(1)
	go cloudFrame.FrameFetcher(infraredCam, config.SampleSecond, &waitGroup, dataChannel)
	fmt.Println("done.")

	scoreChannel := make(chan float32, 2)

	fmt.Print("Starting the cloud evaluator routine...")
	waitGroup.Add(1)
	go cloudFrame.CloudEvaluator(&waitGroup, dataChannel, scoreChannel)
	fmt.Println("done.")

	fmt.Print("Starting the score processing routine...")
	waitGroup.Add(1)
	go httpClient.ScoreUploader(config.DataServer, &waitGroup, scoreChannel)
	fmt.Println("done.")

	waitGroup.Wait()

	fmt.Println("Done")
}
