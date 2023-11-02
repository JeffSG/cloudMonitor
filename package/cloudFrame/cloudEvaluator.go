package cloudFrame

import (
	"cloudMonitor/package/config"
	"cloudMonitor/package/constant"
	"cloudMonitor/package/utils"
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var centerWeight float32
var rangeSquare float32
var reference []float32

const centX = (constant.FRAME_WIDTH - 1) / 2.0
const centY = (constant.FRAME_HEIGHT - 1) / 2.0

func LoadConfig(config *config.ConfigData) {
	centerWeight = config.CenterWeight
	rangeSquare = config.DetectRange * config.DetectRange
}

func LoadReference(filePath string) error {
	dataBytes, err := os.ReadFile(filePath)
	if nil != err {
		return errors.New(utils.ConcatStrings("Failed to read reference file ", filePath, ". Reason: ", err.Error()))
	}
	re := regexp.MustCompile(" +")
	items := re.Split(string(dataBytes), -1)
	itemCount := len(items)
	if itemCount < constant.FRAME_HEIGHT*constant.FRAME_WIDTH {
		return errors.New("Invalid reference file data. Too few data points.")
	}
	reference = make([]float32, len(items))
	for i := 0; i < itemCount; i++ {
		v, _ := strconv.ParseFloat(items[i], 32)
		reference[i] = float32(v)
	}
	return nil
}

func CloudEvaluator(waitGroup *sync.WaitGroup, dataCh chan []float32, scoreCh chan float32) {
	defer waitGroup.Done()

	for {
		frameData := <-dataCh
		frameScore := calculateFrameScore(frameData)
		scoreCh <- frameScore
	}
}

func calculateFrameScore(frameData []float32) float32 {
	var weight float32 // default value is 0
	var value float32  // default value is 0
	for y := 0; y < constant.FRAME_HEIGHT; y++ {
		for x := 0; x < constant.FRAME_WIDTH; x++ {
			dsquare := (float32(x)-centX)*(float32(x)-centX) + (float32(y)-centY)*(float32(y)-centY)
			if rangeSquare < dsquare {
				continue
			}
			w := 1.0 / math.Pow(float64(dsquare), float64(centerWeight)/2.0)
			weight += float32(w)
			id := y*constant.FRAME_WIDTH + x
			if frameData[id] <= reference[id] {
				continue
			}
			v := frameData[id] - reference[id]
			value += v * float32(w)
		}
	}
	if 0 == weight {
		value = 0
	} else {
		value /= weight
	}

	if value <= 100 {
		return value
	} else {
		return 100.0
	}
}
