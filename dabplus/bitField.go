package dabplus

type BitField uint8
type BitPos uint8

func (bf BitField) Has(bp BitPos) bool {
	if bf&(0x01<<bp) == 0x00 {
		return false
	} else {
		return true
	}
}

func (bf *BitField) Set(bp BitPos) {
	*bf |= (0x01 << bp)
}

func (bf *BitField) Clear(bp BitPos) {
	*bf &^= (0x01 << bp)
}


var (
    // Bit positions in status[0]
    ClearToSend BitPos = 7
    CmdError BitPos = 6
    LinkChange BitPos = 5
    ServiceChange BitPos = 4
    TuneDone BitPos = 0

    // Bit positions in status[0]
    EventChange BitPos = 5
)
