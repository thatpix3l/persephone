package entrypoint

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/thatpix3l/persephone/pkg/command"
	"tinygo.org/x/bluetooth"
)

var (
	isScanning   bool
	adapter      = bluetooth.DefaultAdapter
	adapterScans = make(map[string]*bluetooth.ScanResult)
	goproDevice  *bluetooth.Device

	reader = *bufio.NewReader(os.Stdin)

	commandRequestHandle *bluetooth.DeviceCharacteristic
	commandResponse      command.Response
)

const (
	HANDLE_REQUEST_COMMAND   = "b5f90072-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_RESPONSE_COMMAND  = "b5f90073-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_NOTIFY_COMMAND    = "00002902-0000-1000-8000-00805f9b34fb"
	HANDLE_REQUEST_SETTING   = "b5f90074-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_RESPONSE_SETTINGS = "b5f90075-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_NOTIFY_SETTING    = HANDLE_NOTIFY_COMMAND
	HANDLE_REQUEST_QUERY     = "b5f90076-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_RESPONSE_QUERY    = "b5f90077-aa8d-11e3-9046-0002a5d5c51b"
	HANDLE_NOTIFY_QUERY      = HANDLE_NOTIFY_COMMAND
)

func readLine() (string, error) {
	text, err := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	return text, err
}

// Start scanning for and store references to bluetooth devices with "GoPro" substring
func startScan() {
	go func() {
		log.Println("Starting bluetooth adapter...")
		isScanning = true
		adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) {
			if strings.Contains(sr.LocalName(), "GoPro") {
				adapterScans[sr.LocalName()] = &sr
			}

		})
	}()
}

// Stop scanning for bluetooth devices
func stopScan() {
	log.Println("Stopping bluetooth adapter...")
	isScanning = false
	adapter.StopScan()
}

// Connect to a particular bluetooth device with its broadcasted/friendly name
func connectDevice() (*bluetooth.Device, error) {

	input := "GoPro 0357"

	scanResult := adapterScans[input]
	if scanResult == nil {
		return nil, fmt.Errorf("no device called \"%s\" is available", input)

	}

	// Connect to device with address from scan result
	addr := scanResult.Address
	device, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, fmt.Errorf("error connecting \"%s\": %s", input, err)

	}

	return device, nil

}

func getCharacteristic(device *bluetooth.Device, uuid string) (*bluetooth.DeviceCharacteristic, error) {

	// Store slice of services from device
	services, err := device.DiscoverServices(nil)
	if err != nil {
		log.Printf("Error discovering services from device: %s", err)
		return nil, err
	}

	var devChar *bluetooth.DeviceCharacteristic

	// For each service from device...
	for serviceIdx := 0; serviceIdx < len(services); serviceIdx++ {

		// Current service
		service := &services[serviceIdx]

		// Store slice of characteristics from service
		devChars, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			log.Printf("Error dicovering characteristics from device: %s", err)
			continue

		}

		// For each device characteristic from current service...
		for devCharIdx := 0; devCharIdx < len(devChars); devCharIdx++ {

			// If the characteristic's UUID matches what was given, store reference to it
			if devChars[devCharIdx].UUID().String() == uuid {
				devChar = &devChars[devCharIdx]
				break
			}

		}

	}

	if devChar == nil {
		return nil, fmt.Errorf("given UUID \"%s\" does not exist on this device", uuid)

	}

	return devChar, nil

}

func Start() {

	// Enable bluetooth adapter
	if err := adapter.Enable(); err != nil {
		log.Println(err)
		return

	}

	// Start adapter in background
	startScan()

	// Log the list of available GoPro devices every 3 seconds
	go func() {
		for {
			if isScanning {
				goproNames := ""
				for name := range adapterScans {

					if goproNames == "" {
						goproNames = name

					} else {
						goproNames += ", " + name

					}

				}
				log.Println("Available GoPro devices:", goproNames)
			}

			time.Sleep(time.Second * 3)

		}

	}()

	// Blocking loop for user input and actions
	for {

		// User input
		var input, err = readLine()
		if err != nil {
			continue
		}

		// Switch cases on possible user actions
		switch input {

		case "start": // Start scanning for devices if not already
			if !isScanning {
				go startScan()
			}

		case "stop": // Stop scanning for devices
			stopScan()

		case "connect": // Connect to device
			{

				if device, err := connectDevice(); err != nil {
					log.Println(err)
					break

				} else {
					goproDevice = device

				}

				commandResponseHandle, err := getCharacteristic(goproDevice, HANDLE_RESPONSE_COMMAND)
				if err != nil {
					continue
				}

				commandResponseHandle.EnableNotifications(func(buf []byte) {
					if err := commandResponse.Unmarshal(buf); err != nil {
						log.Println(err)
					}
					log.Println(commandResponse)
				})

				respChars := []*bluetooth.DeviceCharacteristic{}

				for _, handle := range []string{HANDLE_RESPONSE_QUERY, HANDLE_RESPONSE_SETTINGS} {

					devChar, err := getCharacteristic(goproDevice, handle)
					if err != nil {
						continue
					}

					respChars = append(respChars, devChar)

				}

				if len(respChars) == 0 {
					break
				}

				for i := 0; i < len(respChars); i++ {

					char := respChars[i]

					log.Printf("Enabling notifications on query response handle \"%s\"", char.UUID().String())
					char.EnableNotifications(func(buf []byte) {
						log.Printf("Characteristic \"%s\" has data %v", char.UUID().String(), buf)

					})

				}

				commandRequestHandle, _ = getCharacteristic(goproDevice, HANDLE_REQUEST_COMMAND)

				log.Println("Done")

			}

		case "version":
			commandRequestHandle.WriteWithoutResponse(command.Action.GetVersion())

		case "hardware info":
			commandRequestHandle.WriteWithoutResponse(command.Action.GetHardwareInfo())

		case "shutter on":
			commandRequestHandle.WriteWithoutResponse(command.Action.TurnShutterOn())

		case "shutter off":
			commandRequestHandle.WriteWithoutResponse(command.Action.TurnShutterOff())

		case "set date time":
			joe, _ := time.LoadLocation("America/Chicago")
			commandRequestHandle.WriteWithoutResponse(command.Action.SetDateTime(time.Date(2018, 1, 31, 3, 4, 5, 0, joe)))
			commandRequestHandle.WriteWithoutResponse(command.Action.GetHardwareInfo())

		case "analytics":
			commandRequestHandle.WriteWithoutResponse(command.Action.Analytics())

		}

	}

}
