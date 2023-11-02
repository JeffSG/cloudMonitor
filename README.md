# Cloud Monitor App
Jeff Shangguan, shangguan@hotmail.com

## Description

CloudMonitorGo is an application to monitor the cloud coverage of the night sky. It uses an infrared camera to capture the temperature of sky. Then an algorithm is applied to the sky temperature to evaluate the cloud coverage % and thickness.

## Application

### Features
1. Configuration driven. All the behaviors are conducted from a Yaml configuration file.
2. Communication with the infrared camera with COM port.
3. Flex and robust algorithm to nspect the cloud coverage.
4. Communication with the Weather server

### Syntax
```
go run .\mainProc.go -config config.yaml -reference reference.txt
```

### Modules
1. main @cloudMonitor/mainProc.go - the main process.
2. config @cloudMonitor/package/config/config.go - load and manage the configuration items.
3. utils @cloudMonitor/package/utils/utils.go - the global common functions.
4. constant @cloudMonitor/package/constant/constant.go - the global constant values.

### Pre-request
1. Yaml package
```
go get gopkg.in/yaml.v3
```
2. Serial package
```
go get github.com/tarm/serial
```

### Dataflow
1. Load configuration items from input argument.
2. Connect to the COM port device.
3. Send request to the COM port to retrieve the current temperature data.
4. Process the data frame to evaluate the cloud coverage.
5. Upload the result and the temperature image to the server.

### COM Port Communication
Tested with Virtual Serial Port Driver V6.9

#### Commands
| ID | Command | Send | Received | In init | Description |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | Disable auto output | $SETP=7,0\r\n | OK\r\n | Yes | Disable, otherwise the COM device will send the endless sample frames |
| 2 | Turn off the LED | $SETP=9,0\r\n | OK\r\n | Yes | Must turn off the LED before the normal process |
| 3 | Capture one frame | $SETP=7,1\r\n | OK - 4 bytes, MLX_TMP - 9 bytes, then data segment | No | Data frames are received continuously, starting with MLX_TMP |
| 4 | Stop capture | $SETP=7,0\r\n | OK\r\n | No | The same command as 1, but in data processing stage |
