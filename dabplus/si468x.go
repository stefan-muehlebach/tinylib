package dabplus

import (
	"embed"
)

//go:embed si468x.rom
var embedFS embed.FS

const (
	SI46XX_RD_REPLY                     = 0x00
	SI46XX_POWER_UP                     = 0x01
	SI46XX_HOST_LOAD                    = 0x04
	SI46XX_FLASH_LOAD                   = 0x05
	SI46XX_LOAD_INIT                    = 0x06
	SI46XX_BOOT                         = 0x07
	SI46XX_GET_PART_INFO                = 0x08
	SI46XX_GET_SYS_STATE                = 0x09
	SI46XX_GET_POWER_UP_ARGS            = 0x0A
	SI46XX_READ_OFFSET                  = 0x10
	SI46XX_GET_FUNC_INFO                = 0x12
	SI46XX_SET_PROPERTY                 = 0x13
	SI46XX_GET_PROPERTY                 = 0x14
	SI46XX_WRITE_STORAGE                = 0x15
	SI46XX_READ_STORAGE                 = 0x16
	SI46XX_GET_DIGITAL_SERVICE_LIST     = 0x80
	SI46XX_START_DIGITAL_SERVICE        = 0x81
	SI46XX_STOP_DIGITAL_SERVICE         = 0x82
	SI46XX_GET_DIGITAL_SERVICE_DATA     = 0x84
	SI46XX_DAB_TUNE_FREQ                = 0xB0
	SI46XX_DAB_DIGRAD_STATUS            = 0xB2
	SI46XX_DAB_GET_EVENT_STATUS         = 0xB3
	SI46XX_DAB_GET_ENSEMBLE_INFO        = 0xB4
	SI46XX_DAB_GET_SERVICE_LINKING_INFO = 0xB7
	SI46XX_DAB_SET_FREQ_LIST            = 0xB8
	SI46XX_DAB_GET_FREQ_LIST            = 0xB9
	SI46XX_DAB_GET_COMPONENT_INFO       = 0xBB
	SI46XX_DAB_GET_TIME                 = 0xBC
	SI46XX_DAB_GET_AUDIO_INFO           = 0xBD
	SI46XX_DAB_GET_SUBCHAN_INFO         = 0xBE
	SI46XX_DAB_GET_FREQ_INFO            = 0xBF
	SI46XX_DAB_GET_SERVICE_INFO         = 0xC0
	SI46XX_TEST_GET_RSSI                = 0xE5
	SI46XX_DAB_TEST_GET_BER_INFO        = 0xE8
)
