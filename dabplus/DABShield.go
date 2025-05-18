package dabplus

import (
	"errors"
	"io"
	"machine"
	"time"
	"tinylib"
	"tinylib/conf"
)

const (
	dabMaxServices        = 32
	dabMaxServiceDataLen  = 128
	dabMaxEnsembleNameLen = 17
	spiBufSize            = 512
	defDABPollRate        = 100 * time.Millisecond
)

var (
	FreqList = [...]uint32{
		174928, 176640, 178352, 180064, 181936, 183648, 185360, 187072,
		188928, 190640, 192352, 194064, 195936, 197648, 199360, 201072,
		202928, 204640, 206352, 208064, 209936, 211648, 213360, 215072,
		216928, 218640, 220352, 222064, 223936, 225648, 227360, 229072,
		230784, 232496, 234208, 235776, 237488, 239200,
	}
	NumFreqs = len(FreqList)
)

var (
	ProgTypeList = [...]string{
		"None",
		"News",
		"Current affairs",
		"Information",
		"Sport",
		"Education",
		"Drama",
		"Culture",
		"Science",
		"Varied",
		"Pop music",
		"Rock music",
		"Easy listening music",
		"Light classical",
		"Serious classical",
		"Other music",
		"Weather",
		"Finance",
		"Children’s programmes",
		"Social Affairs",
		"Religion",
		"Phone In",
		"Travel",
		"Leisure",
		"Jazz music",
		"Country music",
		"National music",
		"Oldies music",
		"Folk music",
		"Documentary",
		"Alarm test",
		"Alarm",
	}
	NumProgTypes = len(ProgTypeList)
)

var (
	spiBuf []byte = make([]byte, spiBufSize+8)
	rcvBuf []byte = make([]byte, spiBufSize+8)
)

var (
	ErrUnspecified     = errors.New("dabplus: unspecified error")
	ErrReplyOverflow   = errors.New("dabplus: reply overflow")
	ErrNotAvailable    = errors.New("dabplus: not available")
	ErrNotSupported    = errors.New("dabplus: not supported")
	ErrBadFreq         = errors.New("dabplus: bad frequency")
	ErrCmdNotFound     = errors.New("dabplus: command not found")
	ErrBadArg1         = errors.New("dabplus: bad argument 1")
	ErrBadArg2         = errors.New("dabplus: bad argument 2")
	ErrBadArg3         = errors.New("dabplus: bad argument 3")
	ErrBadArg4         = errors.New("dabplus: bad argument 4")
	ErrBadArg5         = errors.New("dabplus: bad argument 5")
	ErrBadArg6         = errors.New("dabplus: bad argument 6")
	ErrBadArg7         = errors.New("dabplus: bad argument 7")
	ErrCmdBusy         = errors.New("dabplus: command busy")
	ErrAtBandLimit     = errors.New("dabplus: at band limit")
	ErrBadNVM          = errors.New("dabplus: bad NVM")
	ErrBadPatch        = errors.New("dabplus: bad patch")
	ErrBadBootmode     = errors.New("dabplus: bad bootmode")
	ErrBadProperty     = errors.New("dabplus: bad property")
	ErrNotAcquired     = errors.New("dabplus: not acquired")
	ErrAPPnotSupported = errors.New("dabplus: APP not supported")

	ErrMap = map[uint8]error{
		0x01: ErrUnspecified,
		0x02: ErrReplyOverflow,
		0x03: ErrNotAvailable,
		0x04: ErrNotSupported,
		0x05: ErrBadFreq,
		0x10: ErrCmdNotFound,
		0x11: ErrBadArg1,
		0x12: ErrBadArg2,
		0x13: ErrBadArg3,
		0x14: ErrBadArg4,
		0x15: ErrBadArg5,
		0x16: ErrBadArg6,
		0x17: ErrBadArg7,
		0x18: ErrCmdBusy,
		0x19: ErrAtBandLimit,
		0x20: ErrBadNVM,
		0x30: ErrBadPatch,
		0x31: ErrBadBootmode,
		0x40: ErrBadProperty,
		0x50: ErrNotAcquired,
		0xFF: ErrAPPnotSupported,
	}
)

type DABSpeaker int

const (
	SpeakerNone DABSpeaker = iota
	SpeakerDiff
	SpeakerStereo
)

type ServiceType int

const (
	ServiceNone ServiceType = iota
	ServiceAudio
	ServiceData
)

type AudioMode int

const (
	Dual AudioMode = iota
	Mono
	Stereo
	JointStereo
)

func (m AudioMode) String() string {
	switch m {
	case Dual:
		return "Dual"
	case Mono:
		return "Mono"
	case Stereo:
		return "Stereo"
	case JointStereo:
		return "JointStereo"
	default:
		return "(unspec. audio mode)"
	}
}

type ServiceDataCallback func(data string)

type DABService struct {
	FreqId    uint8
	ServiceId uint32
	CompId    uint32
	Label     string
}

type DABConfig struct {
	CSPin, ResetPin, PowerEnablePin, IntPin machine.Pin
}

type DAB struct {
	spi                                     *machine.SPI
	i2c                                     *machine.I2C
	uart                                    *machine.UART
	csPin, intPin, resetPin, powerEnablePin machine.Pin

	serviceList []DABService
	numServices int
	pro         bool
	ServiceData []byte

	ChipRevision, RomID, VerMajor, VerMinor, VerBuild uint8
	PartNo                                            uint16

	ensembleId        uint32
	ensembleName      []byte
	ecc               uint16
	freqId            uint8
	serviceId, compId uint32
	freq              uint16
	// signalStrength    int8
	// snr               int8
	// quality           uint8
	valid    bool
	CmdError uint8

	// bitRate, sampleRate uint16
	// audioMode           AudioMode
	typ     ServiceType
	dabPlus bool
	pty     uint8

	callback ServiceDataCallback
}

// ---------------------------------------------------------------------------
// Public Methoden
func (d *DAB) Configure(cfg DABConfig) {
	d.csPin = cfg.CSPin
	d.intPin = cfg.IntPin
	d.powerEnablePin = cfg.PowerEnablePin
	d.resetPin = cfg.ResetPin

	d.spi = machine.SPI0
	d.i2c = machine.I2C0
	d.uart = machine.UART0

	// println("DAB.Configure(): configure SPI bus...")
	if err := d.spi.Configure(machine.SPIConfig{
		Frequency: 2_000_000,
		SCK:       conf.PinSCK,
		SDO:       conf.PinMOSI,
		SDI:       conf.PinMISO,
	}); err != nil {
		println("  error:", err.Error())
	}

	// println("DAB.Configure(): configure I2C bus...")
	if err := d.i2c.Configure(machine.I2CConfig{
		Frequency: 100_000,
		SDA:       conf.PinSDA,
		SCL:       conf.PinSCL,
	}); err != nil {
		println("  error:", err.Error())
	}
	// println("DAB.Configure(): configure serial lines RX/TX...")
	if err := d.uart.Configure(machine.UARTConfig{
		RX: conf.PinRX,
		TX: conf.PinTX,
	}); err != nil {
		println("  error:", err.Error())
	}

	// println("DAB.Configure(): configure control lines...")
	d.resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.powerEnablePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.intPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	d.freqId = 0xFF

	d.serviceList = make([]DABService, 0, dabMaxServices)
	d.ServiceData = make([]byte, dabMaxServiceDataLen)
	d.ensembleName = make([]byte, dabMaxEnsembleNameLen)

	d.Reset()
}

func (d *DAB) Reset() {
	// println("DAB.Reset(): si468x_reset()")
	d.si468x_reset()
	// println("  error:", d.CmdError)
	// println("DAB.Reset(): si468x_init_dab()")
	d.si468x_init_dab()
	// println("  error:", d.CmdError)
	// println("DAB.Reset(): si468x_get_part_info()")
	d.si468x_get_part_info()
	// println("  error:", d.CmdError)
	// println("DAB.Reset(): si468x_get_func_info()")
	d.si468x_get_func_info()
	// println("  error:", d.CmdError)

	// println("DAB.Reset(): check if we are on a PRO version")
	if err := d.i2c.Tx(0x1A, nil, nil); err != nil {
		println("Error connecting over I2C:", err.Error())
	} else {
		d.pro = true
	}

	if d.pro {
		// println("DAB.Reset(): set up digital audio slave...")

		d.si468x_set_property(0x0200, 0x0000)
		d.si468x_set_property(0x0201, 48000)
		d.si468x_set_property(0x0800, 0x0002)

		d.nau8822_write_reg(NAU8822_REG_CLOCKING, 0x009)
		d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_1, 0x00D)
		d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_2, 0x180)
		d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_3, 0x00F)
	}
}

func (d *DAB) Speaker(value DABSpeaker) {
	if d.pro {
		switch value {
		case SpeakerNone:
			d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_3, 0x00F)
			d.nau8822_write_reg(NAU8822_REG_RIGHT_SPEAKER_CONTROL, 0x000)
		case SpeakerStereo:
			d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_3, 0x06F)
			d.nau8822_write_reg(NAU8822_REG_RIGHT_SPEAKER_CONTROL, 0x000)
		case SpeakerDiff:
			d.nau8822_write_reg(NAU8822_REG_POWER_MANAGEMENT_3, 0x06F)
			d.nau8822_write_reg(NAU8822_REG_RIGHT_SPEAKER_CONTROL, 0x010)
		}
	}
}

func (d *DAB) Volume(value uint8) {
	if d.pro {
		d.nau8822_write_reg(NAU8822_REG_LHP_VOLUME, 0x000|uint16(value&0x3F))
		d.nau8822_write_reg(NAU8822_REG_RHP_VOLUME, 0x100|uint16(value&0x3F))
		d.nau8822_write_reg(NAU8822_REG_LSPKOUT_VOLUME, 0x000|uint16(value&0x3F))
		d.nau8822_write_reg(NAU8822_REG_RSPKOUT_VOLUME, 0x100|uint16(value&0x3F))
	} else {
		d.si468x_set_property(0x0300, uint16(value&0x3F))
	}
}

func (d *DAB) Bass(level int8) {
	if d.pro {
		d.nau8822_write_reg(NAU8822_REG_EQ1, 0x120|uint16((12-level)&0x1F))
	}
}

func (d *DAB) Middle(level int8) {
	if d.pro {
		d.nau8822_write_reg(NAU8822_REG_EQ3, 0x120|uint16((12-level)&0x1F))
	}
}

func (d *DAB) Treble(level int8) {
	if d.pro {
		d.nau8822_write_reg(NAU8822_REG_EQ5, 0x000|uint16((12-level)&0x1F))
	}
}

func (d *DAB) Tune(freqId uint8) {
	d.freqId = freqId
	d.si468x_dab_tune_freq(d.freqId)
	if d.CmdError != 0 {
		d.freqId = 0xFF
		return
	}
	d.getEnsembleInfo()
}

func (d *DAB) TuneService(freqId uint8, serviceId, compId uint32) error {
	if d.freqId != freqId {
		d.freqId = freqId
		d.si468x_dab_tune_freq(d.freqId)
	}
	d.serviceId = serviceId
	d.compId = compId

	timeout := 1000
	for {
		time.Sleep(4 * time.Millisecond)
		d.si468x_start_digital_service(d.serviceId, d.compId)
		if d.CmdError == 0 {
			break
		}
		timeout--
		if timeout == 0 {
			break
		}
	}
	return nil
}

func (d *DAB) AudioInfo() (bitRate, sampleRate uint16, audioMode AudioMode) {
	d.si468x_get_audio_info()
	d.si468x_responseN(19)
	bitRate = uint16(spiBuf[5]) + (uint16(spiBuf[6]) << 8)
	sampleRate = uint16(spiBuf[7]) + (uint16(spiBuf[8]) << 8)
	audioMode = AudioMode(spiBuf[9] & 0x03)
	return
}

func (d *DAB) RadioStatus() (signalStrength, snr, quality uint8) {
	d.si468x_dab_digrad_status()
	d.si468x_responseN(22)
	signalStrength = spiBuf[7]
	snr = spiBuf[8]
	quality = spiBuf[9]
	return
}

func (d *DAB) ComponentInfo() {
	d.si468x_get_subchan_info(d.serviceId, d.compId)
}

func (d *DAB) Time() time.Time {
	spiBuf[0] = SI46XX_DAB_GET_TIME
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
	d.si468x_responseN(12)
	Y := int(spiBuf[5]) + (int(spiBuf[6]) << 8)
	M := time.Month(spiBuf[7])
	D := int(spiBuf[8])
	hh := int(spiBuf[9])
	mm := int(spiBuf[10])
	ss := int(spiBuf[11])
	loc, _ := time.LoadLocation("Local")
	return time.Date(Y, M, D, hh, mm, ss, 0, loc)
}

func (d *DAB) SetServiceCallback(cb ServiceDataCallback) {
	d.callback = cb
}

func (d *DAB) Task() *tinylib.Task {
	return tinylib.NewTask(d.Tick, tinylib.TaskConfig{Interval: defDABPollRate})
}

func (d *DAB) Tick() {
	d.si468x_response()
	status0 := BitField(spiBuf[1])
	status1 := BitField(spiBuf[2])
	if status0.Has(LinkChange) {
		// Gem. Doku mit DAB_DIGRAD_STATUS behandeln
		// println(">>> LinkChange on DABShield")
	}
	if status0.Has(ServiceChange) {
		// Gem. Doku mit GET_DIGITAL_SERVICE_DATA behandeln
		// wie hier gemacht.
		// println(">>> ServiceChange on DABShield <<<")
		d.si468x_get_digital_service_data()
		d.si468x_responseN(20)
		len := int(spiBuf[19]) | (int(spiBuf[20]) << 8)
		if len < spiBufSize-24 {
			d.si468x_responseN(len + 24)
		} else {
			d.si468x_responseN(spiBufSize - 1)
		}
		data := d.parseServiceData()
		if d.callback != nil && data != nil {
			d.callback(string(data))
		}
	}
	if status0.Has(TuneDone) {
		// Zeigt an, wann ein Sendersuchlauf abgeschlossen oder noch am Laufen
		// ist.
		// println(">>> Tune is done on DABShield <<<")

	}
	if status1.Has(EventChange) {
		// Gem. Doku ebenfalls mit DAB_DIGRAD_STATUS behandeln - ich bin
		// jedoch unsicher - eher noch mit DAB_GET_EVENT_STATUS.
		// println(">>> EventChange on DABShield <<<")
	}
}

// ---------------------------------------------------------------------------
// Private Methoden
func (d *DAB) getEnsembleInfo() {
	d.si468x_dab_digrad_status()
	d.si468x_responseN(23)
	if d.isServiceValid() {
		d.waitServiceList()
		d.si468x_dab_get_ensemble_info()
		d.si468x_responseN(29)
		d.ensembleId = uint32(spiBuf[5]) + (uint32(spiBuf[6]) << 8)
		copy(d.ensembleName, spiBuf[7:23])
		d.ecc = uint16(spiBuf[23])
		d.si468x_get_digital_service_list()
		d.si468x_responseN(6)
		len := uint16(spiBuf[5]) + (uint16(spiBuf[6]) << 8) + 2
		if len < spiBufSize {
			d.si468x_responseN(int(len) + 4)
			d.parseServiceList()
		} else {
			offset := uint16(0)
			first := true
			for len >= spiBufSize {
				d.si468x_readoffset(offset)
				d.si468x_responseN(spiBufSize + 4)
				if d.parseServiceListPart(first, spiBuf[5:], spiBufSize) {
					len = 0
				} else {
					len -= spiBufSize
				}
				first = false
				offset += spiBufSize
			}
			if len > 0 {
				d.si468x_readoffset(offset)
				d.si468x_responseN(int(len) + 4)
				d.parseServiceListPart(first, spiBuf[5:], len)
			}
		}
	} else {
		// No services
	}
}

func (d *DAB) isServiceValid() bool {
	d.si468x_dab_digrad_status()
	d.si468x_responseN(23)

	if ((spiBuf[6] & 0x01) == 0x01) && (spiBuf[7] > 0x20) && (spiBuf[9] > 25) {
		return true
	} else {
		return false
	}
}

func (d *DAB) waitServiceList() {
	timeout := 1000
	for {
		time.Sleep(4 * time.Millisecond)
		d.si468x_dab_get_event_status()
		d.si468x_responseN(8)
		timeout--
		if timeout == 0 {
			d.CmdError |= 0x80
			break
		}
		if spiBuf[6]&0x01 != 0x00 {
			break
		}
	}
}

func (d *DAB) parseServiceList() {
	// TO DO
}

func (d *DAB) parseServiceListPart(first bool, data []byte, len uint16) bool {
	// TO DO
	return true
}

func (d *DAB) parseServiceData() []byte {
	byteCount := int(spiBuf[19]) + (int(spiBuf[20]) << 8)
	dataSrc := (spiBuf[8] >> 6) & 0x03

	if dataSrc == 0x02 {
		header1 := spiBuf[25]
		if (header1 & 0x10) != 0x10 {
			if byteCount > dabMaxServiceDataLen {
				byteCount = dabMaxServiceDataLen
			}
			return spiBuf[27 : 27+byteCount-3]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Private NAU8822 Methoden
func (d *DAB) nau8822_write_reg(reg uint8, data uint16) {
	var i2cdata [2]byte

	i2cdata[0] = (reg << 1) & 0x7E
	if data&0x100 != 0 {
		i2cdata[0] |= 0x01
	}
	i2cdata[1] = uint8(data & 0xFF)
	if err := d.i2c.Tx(0x1A, i2cdata[:], nil); err != nil {
		println("Couldn't write NAU882 register:", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Private Methoden fuer den Aufruf der SI468-Funktionen.
func (d *DAB) si468x_reset() {
	d.powerEnablePin.High()
	time.Sleep(100 * time.Millisecond)

	d.resetPin.Low()
	time.Sleep(100 * time.Millisecond)
	d.resetPin.High()
	time.Sleep(100 * time.Millisecond)

	// println("si468x_reset(): si468x_power_up()")
	d.si468x_power_up()
	// println("si468x_reset(): si468x_load_init()")
	d.si468x_load_init()
	// println("si468x_reset(): si468x_host_load()")
	d.si468x_host_load()
	// println("si468x_reset(): si468x_load_init()")
	d.si468x_load_init()
	// println("si468x_reset(): si468x_flash_set_property()")
	d.si468x_flash_set_property(0x0001, 10000)
}

func (d *DAB) si468x_init_dab() {

	d.si468x_flash_load(0x6000)
	d.si468x_boot()
	d.si468x_set_freq_list()

	d.si468x_set_property(0x0000, 0x0010)

	d.si468x_set_property(0x1710, 0xF83E)
	d.si468x_set_property(0x1711, 0x01A4)
	d.si468x_set_property(0x1712, 0x0001)

	d.si468x_set_property(0x8100, 0x0001) //enable DSRVPCKTINT
	d.si468x_set_property(0xb400, 0x0007) //enable XPAD data
}

// Dies ist ein Art Warte-Methode, welche periodisch (4ms!) beim si468x
// das sog. CTS (clear to send) Bit abruft.
func (d *DAB) si468x_cts() {
	d.CmdError = 0x00
	timeout := 1_000
	for {
		time.Sleep(4 * time.Millisecond)
		d.si468x_response()
		timeout--
		if timeout == 0 {
			d.CmdError = 0x80
			break
		}
		if BitField(spiBuf[1]).Has(ClearToSend) {
			break
		}
	}
	if BitField(spiBuf[1]).Has(CmdError) {
		d.si468x_responseN(5)
		d.CmdError = 0x80 | spiBuf[5]
	}
}

func (d *DAB) si468x_response() {
	d.si468x_responseN(4)
}

func (d *DAB) si468x_responseN(len int) {
	for i := range len + 1 {
		spiBuf[i] = 0x00
	}
	d.spiRecvData(spiBuf[:len+1])
}

func (d *DAB) si468x_flash_load(flashAddr uint32) {
	spiBuf[0] = SI46XX_FLASH_LOAD
	spiBuf[1] = 0x00
	spiBuf[2] = 0x00
	spiBuf[3] = 0x00

	spiBuf[4] = byte(flashAddr & 0xff)
	spiBuf[5] = byte((flashAddr >> 8) & 0xff)
	spiBuf[6] = byte((flashAddr >> 16) & 0xff)
	spiBuf[7] = byte((flashAddr >> 24) & 0xff)

	spiBuf[8] = 0x00
	spiBuf[9] = 0x00
	spiBuf[10] = 0x00
	spiBuf[11] = 0x00
	d.spiSendData(spiBuf[:12])
	d.si468x_cts()
}

func (d *DAB) si468x_boot() {
	spiBuf[0] = SI46XX_BOOT
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_power_up() {
	spiBuf[0] = SI46XX_POWER_UP
	spiBuf[1] = 0x00
	spiBuf[2] = 0x17
	spiBuf[3] = 0x48
	spiBuf[4] = 0x00
	spiBuf[5] = 0xf8
	spiBuf[6] = 0x24
	spiBuf[7] = 0x01
	spiBuf[8] = 0x1F
	spiBuf[9] = 0x10
	spiBuf[10] = 0x00
	spiBuf[11] = 0x00
	spiBuf[12] = 0x00
	spiBuf[13] = 0x18
	spiBuf[14] = 0x00
	spiBuf[15] = 0x00
	// println("spiSendData()...")
	d.spiSendData(spiBuf[:16])
	// println("si468x_cts()...")
	d.si468x_cts()
}

func (d *DAB) si468x_load_init() {
	spiBuf[0] = SI46XX_LOAD_INIT
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_host_load() {
	fh, err := embedFS.Open("si468x.rom")
	if err != nil {
		panic("os.Open: " + err.Error())
	}
	defer fh.Close()

	for {
		spiBuf[0] = SI46XX_HOST_LOAD
		spiBuf[1] = 0x00
		spiBuf[2] = 0x00
		spiBuf[3] = 0x00

		n, err := fh.Read(spiBuf[4:])
		if n == 0 && err == io.EOF {
			break
		}
		d.spiSendData(spiBuf[:n+4])
		d.si468x_cts()
	}
}

func (d *DAB) si468x_readoffset(offset uint16) {
	spiBuf[0] = SI46XX_READ_OFFSET
	spiBuf[1] = 0x00
	spiBuf[2] = byte(offset & 0xFF)
	spiBuf[3] = byte((offset >> 8) & 0xFF)
	d.spiSendData(spiBuf[:4])
	d.si468x_cts()
}

func (d *DAB) si468x_flash_set_property(property, value uint16) {
	spiBuf[0] = SI46XX_FLASH_LOAD
	spiBuf[1] = 0x10
	spiBuf[2] = 0x0
	spiBuf[3] = 0x0

	spiBuf[4] = byte(property & 0xFF)
	spiBuf[5] = byte((property >> 8) & 0xFF)
	spiBuf[6] = byte(value & 0xFF)
	spiBuf[7] = byte((value >> 8) & 0xFF)
	d.spiSendData(spiBuf[:8])
	d.si468x_cts()
}

func (d *DAB) si468x_set_property(property, value uint16) {
	spiBuf[0] = SI46XX_SET_PROPERTY
	spiBuf[1] = 0x00

	spiBuf[2] = byte(property & 0xFF)
	spiBuf[3] = byte((property >> 8) & 0xFF)
	spiBuf[4] = byte(value & 0xFF)
	spiBuf[5] = byte((value >> 8) & 0xFF)
	d.spiSendData(spiBuf[:6])
	d.si468x_cts()
}

func (d *DAB) si468x_set_freq_list() {
	spiBuf[0] = SI46XX_DAB_SET_FREQ_LIST
	spiBuf[1] = byte(NumFreqs)
	spiBuf[2] = 0x00
	spiBuf[3] = 0x00

	for i, freq := range FreqList {
		spiBuf[4+(i*4)] = byte(freq & 0xFF)
		spiBuf[5+(i*4)] = byte((freq >> 8) & 0xFF)
		spiBuf[6+(i*4)] = byte((freq >> 16) & 0xFF)
		spiBuf[7+(i*4)] = byte((freq >> 24) & 0xFF)
	}
	d.spiSendData(spiBuf[:4+4*NumFreqs])
	d.si468x_cts()
}

func (d *DAB) si468x_dab_tune_freq(freqIndex uint8) {
	spiBuf[0] = SI46XX_DAB_TUNE_FREQ
	spiBuf[1] = 0x00
	spiBuf[2] = freqIndex
	spiBuf[3] = 0x00
	spiBuf[4] = 0x00
	spiBuf[5] = 0x00
	d.spiSendData(spiBuf[:6])
	d.si468x_cts()

	timeout := 1000
	for {
		time.Sleep(4 * time.Millisecond)
		d.si468x_response()
		timeout--
		if timeout == 0 {
			d.CmdError |= 0x80
			break
		}
		if spiBuf[1]&0x01 != 0 {
			break
		}
	}
}

func (d *DAB) si468x_start_digital_service(serviceId, compId uint32) {
	spiBuf[0] = SI46XX_START_DIGITAL_SERVICE
	spiBuf[1] = 0x00
	spiBuf[2] = 0x00
	spiBuf[3] = 0x00
	spiBuf[4] = byte(serviceId & 0xff)
	spiBuf[5] = byte((serviceId >> 8) & 0xff)
	spiBuf[6] = byte((serviceId >> 16) & 0xff)
	spiBuf[7] = byte((serviceId >> 24) & 0xff)
	spiBuf[8] = byte(compId & 0xff)
	spiBuf[9] = byte((compId >> 8) & 0xff)
	spiBuf[10] = byte((compId >> 16) & 0xff)
	spiBuf[11] = byte((compId >> 24) & 0xff)
	d.spiSendData(spiBuf[:12])
	d.si468x_cts()
}

func (d *DAB) si468x_dab_get_ensemble_info() {
	spiBuf[0] = SI46XX_DAB_GET_ENSEMBLE_INFO
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_dab_digrad_status() {
	spiBuf[0] = SI46XX_DAB_DIGRAD_STATUS
	spiBuf[1] = 0x09
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_dab_get_event_status() {
	spiBuf[0] = SI46XX_DAB_GET_EVENT_STATUS
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_get_digital_service_list() {
	spiBuf[0] = SI46XX_GET_DIGITAL_SERVICE_LIST
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_get_digital_service_data() {
	spiBuf[0] = SI46XX_GET_DIGITAL_SERVICE_DATA
	spiBuf[1] = 0x01
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_get_audio_info() {
	spiBuf[0] = SI46XX_DAB_GET_AUDIO_INFO
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
}

func (d *DAB) si468x_get_part_info() {
	spiBuf[0] = SI46XX_GET_PART_INFO
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
	d.si468x_responseN(10)
	d.ChipRevision = spiBuf[5]
	d.RomID = spiBuf[6]
	d.PartNo = (uint16(spiBuf[10]) << 8) | uint16(spiBuf[9])
}

func (d *DAB) si468x_get_func_info() {
	spiBuf[0] = SI46XX_GET_FUNC_INFO
	spiBuf[1] = 0x00
	d.spiSendData(spiBuf[:2])
	d.si468x_cts()
	d.si468x_responseN(8)
	d.VerMajor = spiBuf[5]
	d.VerMinor = spiBuf[6]
	d.VerBuild = spiBuf[7]
}

func (d *DAB) si468x_get_service_info(serviceId uint32) {
	spiBuf[0] = SI46XX_DAB_GET_SERVICE_INFO
	spiBuf[1] = 0x00
	spiBuf[2] = 0x00
	spiBuf[3] = 0x00
	spiBuf[4] = byte(serviceId & 0xFF)
	spiBuf[5] = byte((serviceId >> 8) & 0xFF)
	spiBuf[6] = byte((serviceId >> 16) & 0xFF)
	spiBuf[7] = byte((serviceId >> 24) & 0xFF)
	d.spiSendData(spiBuf[:8])
	d.si468x_cts()
}

func (d *DAB) si468x_get_subchan_info(serviceId, compId uint32) {
	spiBuf[0] = SI46XX_DAB_GET_SUBCHAN_INFO
	spiBuf[1] = 0x00
	spiBuf[2] = 0x00
	spiBuf[3] = 0x00
	spiBuf[4] = byte(serviceId & 0xFF)
	spiBuf[5] = byte((serviceId >> 8) & 0xFF)
	spiBuf[6] = byte((serviceId >> 16) & 0xFF)
	spiBuf[7] = byte((serviceId >> 24) & 0xFF)
	spiBuf[8] = byte(compId & 0xFF)
	spiBuf[9] = byte((compId >> 8) & 0xFF)
	spiBuf[10] = byte((compId >> 16) & 0xFF)
	spiBuf[11] = byte((compId >> 24) & 0xFF)
	d.spiSendData(spiBuf[:12])
	d.si468x_cts()
}

func (d *DAB) spiSendData(data []byte) {
	d.csPin.Low()
	for _, ch := range data {
		if _, err := d.spi.Transfer(ch); err != nil {
			println("spi.Transfer() in sendData failed:", err.Error())
		}
	}
	d.csPin.High()
}

func (d *DAB) spiRecvData(data []byte) {
	d.csPin.Low()
	for i, _ := range data {
		if ch, err := d.spi.Transfer(0x00); err != nil {
			println("spi.Transfer() in receiveData failed:", err.Error())
		} else {
			data[i] = ch
		}
	}
	d.csPin.High()
}
