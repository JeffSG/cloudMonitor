package cloudFrame

import (
	"cloudMonitor/package/comDevice"
	"cloudMonitor/package/constant"
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

func FrameFetcher(infraredCam *comDevice.SerialPortData, intervalSec uint8, waitGroup *sync.WaitGroup, ch chan []float32) {
	defer waitGroup.Done()

	for {
		err := infraredCam.SendCommand(constant.COMCMD_START_CAPTURE)
		if nil != err {
			fmt.Println("Failed to send start capture command to the COM port. Reason", err.Error())
			time.Sleep(constant.COM_FAILED_INTERVAL_SEC)
			continue
		}

		infraredCam.SkipBytes(4) // skip the init OK part
		infraredCam.SkipBytes(9) // skip the frame header MLX_TMP part
		infraredCam.WaitForString(constant.COM_FRAME_HEADER, false)
		rawData := getDataFromFrame(infraredCam.GetBuffer())
		calculableData := generateCalculableData(&rawData)

		infraredCam.SendCommand(constant.COMCMD_STOP_OUTPUT)
		infraredCam.WaitForString(constant.COMCMD_OK, true)

		if len(calculableData) < constant.FRAME_HEIGHT*constant.FRAME_WIDTH {
			fmt.Println("Received an invalid frame. Just skip.")
			time.Sleep(constant.COM_FAILED_INTERVAL_SEC)
			continue
		}

		ch <- calculableData

		time.Sleep(time.Second * time.Duration(intervalSec))
	}
}

func getDataFromFrame(buf *string) string {
	mlxPos := strings.Index(*buf, constant.COM_FRAME_HEADER)
	return (*buf)[1 : mlxPos-3]
}

func generateCalculableData(rawData *string) []float32 {
	re := regexp.MustCompile(" +")
	items := re.Split(*rawData, -1)
	itemCount := len(items)
	calcData := make([]float32, itemCount)
	for i := 0; i < itemCount; i++ {
		bits := binary.LittleEndian.Uint32([]byte(items[i]))
		calcData[i] = math.Float32frombits(bits)
	}
	return calcData
}
