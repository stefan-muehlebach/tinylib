package dabplus

const (
	NAU8822_REG_RESET                    = 0x00
	NAU8822_REG_POWER_MANAGEMENT_1       = 0x01
	NAU8822_REG_POWER_MANAGEMENT_2       = 0x02
	NAU8822_REG_POWER_MANAGEMENT_3       = 0x03
	NAU8822_REG_AUDIO_INTERFACE          = 0x04
	NAU8822_REG_COMPANDING_CONTROL       = 0x05
	NAU8822_REG_CLOCKING                 = 0x06
	NAU8822_REG_ADDITIONAL_CONTROL       = 0x07
	NAU8822_REG_GPIO_CONTROL             = 0x08
	NAU8822_REG_JACK_DETECT_CONTROL_1    = 0x09
	NAU8822_REG_DAC_CONTROL              = 0x0A
	NAU8822_REG_LEFT_DAC_DIGITAL_VOLUME  = 0x0B
	NAU8822_REG_RIGHT_DAC_DIGITAL_VOLUME = 0x0C
	NAU8822_REG_JACK_DETECT_CONTROL_2    = 0x0D
	NAU8822_REG_ADC_CONTROL              = 0x0E
	NAU8822_REG_LEFT_ADC_DIGITAL_VOLUME  = 0x0F
	NAU8822_REG_RIGHT_ADC_DIGITAL_VOLUME = 0x10
	NAU8822_REG_EQ1                      = 0x12
	NAU8822_REG_EQ2                      = 0x13
	NAU8822_REG_EQ3                      = 0x14
	NAU8822_REG_EQ4                      = 0x15
	NAU8822_REG_EQ5                      = 0x16
	NAU8822_REG_DAC_LIMITER_1            = 0x18
	NAU8822_REG_DAC_LIMITER_2            = 0x19
	NAU8822_REG_NOTCH_FILTER_1           = 0x1B
	NAU8822_REG_NOTCH_FILTER_2           = 0x1C
	NAU8822_REG_NOTCH_FILTER_3           = 0x1D
	NAU8822_REG_NOTCH_FILTER_4           = 0x1E
	NAU8822_REG_ALC_CONTROL_1            = 0x20
	NAU8822_REG_ALC_CONTROL_2            = 0x21
	NAU8822_REG_ALC_CONTROL_3            = 0x22
	NAU8822_REG_NOISE_GATE               = 0x23
	NAU8822_REG_PLL_N                    = 0x24
	NAU8822_REG_PLL_K1                   = 0x25
	NAU8822_REG_PLL_K2                   = 0x26
	NAU8822_REG_PLL_K3                   = 0x27
	NAU8822_REG_3D_CONTROL               = 0x29
	NAU8822_REG_RIGHT_SPEAKER_CONTROL    = 0x2B
	NAU8822_REG_INPUT_CONTROL            = 0x2C
	NAU8822_REG_LEFT_INP_PGA_CONTROL     = 0x2D
	NAU8822_REG_RIGHT_INP_PGA_CONTROL    = 0x2E
	NAU8822_REG_LEFT_ADC_BOOST_CONTROL   = 0x2F
	NAU8822_REG_RIGHT_ADC_BOOST_CONTROL  = 0x30
	NAU8822_REG_OUTPUT_CONTROL           = 0x31
	NAU8822_REG_LEFT_MIXER_CONTROL       = 0x32
	NAU8822_REG_RIGHT_MIXER_CONTROL      = 0x33
	NAU8822_REG_LHP_VOLUME               = 0x34
	NAU8822_REG_RHP_VOLUME               = 0x35
	NAU8822_REG_LSPKOUT_VOLUME           = 0x36
	NAU8822_REG_RSPKOUT_VOLUME           = 0x37
	NAU8822_REG_AUX2_MIXER               = 0x38
	NAU8822_REG_AUX1_MIXER               = 0x39
	NAU8822_REG_POWER_MANAGEMENT_4       = 0x3A
	NAU8822_REG_LEFT_TIME_SLOT           = 0x3B
	NAU8822_REG_MISC                     = 0x3C
	NAU8822_REG_RIGHT_TIME_SLOT          = 0x3D
	NAU8822_REG_DEVICE_REVISION          = 0x3E
	NAU8822_REG_DEVICE_ID                = 0x3F
	NAU8822_REG_DAC_DITHER               = 0x41
	NAU8822_REG_ALC_ENHANCE_1            = 0x46
	NAU8822_REG_ALC_ENHANCE_2            = 0x47
	NAU8822_REG_192KHZ_SAMPLING          = 0x48
	NAU8822_REG_MISC_CONTROL             = 0x49
	NAU8822_REG_INPUT_TIEOFF             = 0x4A
	NAU8822_REG_POWER_REDUCTION          = 0x4B
	NAU8822_REG_AGC_PEAK2PEAK            = 0x4C
	NAU8822_REG_AGC_PEAK_DETECT          = 0x4D
	NAU8822_REG_AUTOMUTE_CONTROL         = 0x4E
	NAU8822_REG_OUTPUT_TIEOFF            = 0x4F
	NAU8822_REG_MAX_REGISTER             = NAU8822_REG_OUTPUT_TIEOFF

	/* NAU8822_REG_POWER_MANAGEMENT_1 (0x1) */
	NAU8822_REFIMP_MASK = 0x3
	NAU8822_REFIMP_80K  = 0x1
	NAU8822_REFIMP_300K = 0x2
	NAU8822_REFIMP_3K   = 0x3
	NAU8822_IOBUF_EN    = (0x1 << 2)
	NAU8822_ABIAS_EN    = (0x1 << 3)
	NAU8822_PLL_EN_MASK = (0x1 << 5)
	NAU8822_PLL_ON      = (0x1 << 5)
	NAU8822_PLL_OFF     = (0x0 << 5)

	/* NAU8822_REG_AUDIO_INTERFACE (0x4) */
	NAU8822_AIFMT_MASK = (0x3 << 3)
	NAU8822_WLEN_MASK  = (0x3 << 5)
	NAU8822_WLEN_20    = (0x1 << 5)
	NAU8822_WLEN_24    = (0x2 << 5)
	NAU8822_WLEN_32    = (0x3 << 5)
	NAU8822_LRP_MASK   = (0x1 << 7)
	NAU8822_BCLKP_MASK = (0x1 << 8)

	/* NAU8822_REG_COMPANDING_CONTROL (0x5) */
	NAU8822_ADDAP_SFT = 0
	NAU8822_ADCCM_SFT = 1
	NAU8822_DACCM_SFT = 3

	/* NAU8822_REG_CLOCKING (0x6) */
	NAU8822_CLKIOEN_MASK = 0x1
	NAU8822_CLK_MASTER   = 0x1
	NAU8822_CLK_SLAVE    = 0x0
	NAU8822_MCLKSEL_SFT  = 5
	NAU8822_MCLKSEL_MASK = (0x7 << 5)
	NAU8822_BCLKSEL_SFT  = 2
	NAU8822_BCLKSEL_MASK = (0x7 << 2)
	NAU8822_BCLKDIV_1    = (0x0 << 2)
	NAU8822_BCLKDIV_2    = (0x1 << 2)
	NAU8822_BCLKDIV_4    = (0x2 << 2)
	NAU8822_BCLKDIV_8    = (0x3 << 2)
	NAU8822_BCLKDIV_16   = (0x4 << 2)
	NAU8822_CLKM_MASK    = (0x1 << 8)
	NAU8822_CLKM_MCLK    = (0x0 << 8)
	NAU8822_CLKM_PLL     = (0x1 << 8)

	/* NAU8822_REG_ADDITIONAL_CONTROL (0x08) */
	NAU8822_SMPLR_SFT  = 1
	NAU8822_SMPLR_MASK = (0x7 << 1)
	NAU8822_SMPLR_48K  = (0x0 << 1)
	NAU8822_SMPLR_32K  = (0x1 << 1)
	NAU8822_SMPLR_24K  = (0x2 << 1)
	NAU8822_SMPLR_16K  = (0x3 << 1)
	NAU8822_SMPLR_12K  = (0x4 << 1)
	NAU8822_SMPLR_8K   = (0x5 << 1)

	/* NAU8822_REG_EQ1 (0x12) */
	NAU8822_EQ1GC_SFT = 0
	NAU8822_EQ1CF_SFT = 5
	NAU8822_EQM_SFT   = 8

	/* NAU8822_REG_EQ2 (0x13) */
	NAU8822_EQ2GC_SFT = 0
	NAU8822_EQ2CF_SFT = 5
	NAU8822_EQ2BW_SFT = 8

	/* NAU8822_REG_EQ3 (0x14) */
	NAU8822_EQ3GC_SFT = 0
	NAU8822_EQ3CF_SFT = 5
	NAU8822_EQ3BW_SFT = 8

	/* NAU8822_REG_EQ4 (0x15) */
	NAU8822_EQ4GC_SFT = 0
	NAU8822_EQ4CF_SFT = 5
	NAU8822_EQ4BW_SFT = 8

	/* NAU8822_REG_EQ5 (0x16) */
	NAU8822_EQ5GC_SFT = 0
	NAU8822_EQ5CF_SFT = 5

	/* NAU8822_REG_ALC_CONTROL_1 (0x20) */
	NAU8822_ALCMINGAIN_SFT = 0
	NAU8822_ALCMXGAIN_SFT  = 3
	NAU8822_ALCEN_SFT      = 7

	/* NAU8822_REG_ALC_CONTROL_2 (0x21) */
	NAU8822_ALCSL_SFT = 0
	NAU8822_ALCHT_SFT = 4

	/* NAU8822_REG_ALC_CONTROL_3 (0x22) */
	NAU8822_ALCATK_SFT = 0
	NAU8822_ALCDCY_SFT = 4
	NAU8822_ALCM_SFT   = 8

	/* NAU8822_REG_PLL_N (0x24) */
	NAU8822_PLLMCLK_DIV2 = (0x1 << 4)
	NAU8822_PLLN_MASK    = 0xF

	NAU8822_PLLK1_SFT  = 18
	NAU8822_PLLK1_MASK = 0x3F

	/* NAU8822_REG_PLL_K2 (0x26) */
	NAU8822_PLLK2_SFT  = 9
	NAU8822_PLLK2_MASK = 0x1FF

	/* NAU8822_REG_PLL_K3 (0x27) */
	NAU8822_PLLK3_MASK = 0x1FF

	/* NAU8822_REG_RIGHT_SPEAKER_CONTROL (0x2B) */
	NAU8822_RMIXMUT = 0x20
	NAU8822_RSUBBYP = 0x10

	NAU8822_RAUXRSUBG_SFT  = 1
	NAU8822_RAUXRSUBG_MASK = 0x0E

	NAU8822_RAUXSMUT = 0x01
)
