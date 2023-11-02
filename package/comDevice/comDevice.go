package comDevice

import (
	"cloudMonitor/package/constant"
	"cloudMonitor/package/utils"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type SerialPortData struct {
	config *serial.Config
	port   *serial.Port
	buf    string
	cache  []byte
	mutex  sync.Mutex
}

type SerialPort interface {
	// initialize and open a serial port
	OpenPort(comPort string) error

	// send string command to the serial port
	SendCommand(command string) error

	// receive data and put to buffer. The data are kept in string format
	Receive() error

	// skip n bytes from received buffer
	SkipBytes(n uint16)

	// wait until a specified string is received, then clear the content if necessary
	WaitForString(str string, clear bool)

	// close the COM connection
	ClosePort()

	// get the buffer
	GetBuffer() *string
}

func (serialPortData *SerialPortData) OpenPort(comPort string) error {
	serialPortData.config = &serial.Config{Name: comPort, Baud: constant.COM_BAUD, ReadTimeout: constant.COM_TIMEOUT, Parity: serial.ParityNone, Size: 8, StopBits: serial.Stop1}
	port, err := serial.OpenPort(serialPortData.config)
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to open serial port ", comPort, ". Reason: ", err.Error()))
	}
	serialPortData.port = port
	serialPortData.cache = make([]byte, constant.COM_BUF_SIZE)
	return nil
}

func (serialPortData *SerialPortData) SendCommand(command string) error {
	_, err := serialPortData.port.Write([]byte(command))
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to send command ", command, ". Reason: ", err.Error()))
	}
	return nil
}

func (serialPortData *SerialPortData) Receive() error {
	num, err := serialPortData.port.Read(serialPortData.cache)
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to read com data", ". Reason: ", err.Error()))
	}
	if 0 == num {
		return nil
	}
	serialPortData.mutex.Lock()
	serialPortData.buf = utils.ConcatStrings(serialPortData.buf, string(serialPortData.cache[:num]))
	serialPortData.mutex.Unlock()
	return nil
}

func (serialPortData *SerialPortData) SkipBytes(n int) {
	if n <= 0 {
		return
	}

	for {
		if n <= len(serialPortData.buf) {
			serialPortData.mutex.Lock()
			serialPortData.buf = serialPortData.buf[n:]
			serialPortData.mutex.Unlock()
			return
		}
		time.Sleep(constant.COM_RECEIVE_INTERVAL_MS)
	}
}

func (serialPortData *SerialPortData) WaitForString(str string, clear bool) {
	if "" == str {
		return
	}
	for {
		index := strings.Index(serialPortData.buf, str)
		if 0 <= index {
			if clear {
				serialPortData.mutex.Lock()
				startPos := index + len(str)
				serialPortData.buf = serialPortData.buf[startPos:]
				serialPortData.mutex.Unlock()
			}
			return
		}
		time.Sleep(constant.COM_RECEIVE_INTERVAL_MS)
	}
}

func (serialPortData *SerialPortData) Close() {
	serialPortData.port.Close()
}

func (serialPortData *SerialPortData) GetBuffer() *string {
	return &serialPortData.buf
}

// this is the seperated GoRoutine to receive data from the buffer
func Receiver(serialPortData *SerialPortData, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for {
		serialPortData.Receive()
		time.Sleep(constant.COM_RECEIVE_INTERVAL_MS)
	}
}
