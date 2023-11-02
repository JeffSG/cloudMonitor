package constant

import "time"

// Serial port management
const COM_BUF_SIZE = 8000
const COM_RECEIVE_INTERVAL_MS = time.Millisecond * 100
const COM_BAUD = 256000
const COM_TIMEOUT = 5 * time.Second
const COM_FAILED_INTERVAL_SEC = 10 * time.Second

// Serial port command
const COMCMD_STOP_OUTPUT = "$SETP=7,0\r\n"
const COMCMD_TURN_OFF_LED = "$SETP=9,0\r\n"
const COMCMD_START_CAPTURE = "$SETP=7,1\r\n"
const COMCMD_OK = "OK\r\n"
const COM_FRAME_HEADER = "MLX_TMP"

// Cloud frame management
const FRAME_HEIGHT = 24
const FRAME_WIDTH = 32
