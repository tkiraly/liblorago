package liblorago

import (
	"fmt"
	"os"
	"time"
)

var SX125x_TX_DAC_CLK_SEL = 1   /* 0:int, 1:ext */
var SX125x_TX_DAC_GAIN = 2      /* 3:0, 2:-3, 1:-6, 0:-9 dBFS (default 2) */
var SX125x_TX_MIX_GAIN = 14     /* -38 + 2*TxMixGain dB (default 14) */
var SX125x_TX_PLL_BW = 1        /* 0:75, 1:150, 2:225, 3:300 kHz (default 3) */
var SX125x_TX_ANA_BW = 0        /* 17.5 / 2*(41-TxAnaBw) MHz (default 0) */
var SX125x_TX_DAC_BW = 5        /* 24 + 8*TxDacBw Nb FIR taps (default 2) */
var SX125x_RX_LNA_GAIN = 1      /* 1 to 6, 1 highest gain */
var SX125x_RX_BB_GAIN = 12      /* 0 to 15 , 15 highest gain */
var SX125x_LNA_ZIN = 1          /* 0:50, 1:200 Ohms (default 1) */
var SX125x_RX_ADC_BW = 7        /* 0 to 7, 2:100<BW<200, 5:200<BW<400,7:400<BW kHz SSB (default 7) */
var SX125x_RX_ADC_TRIM = 6      /* 0 to 7, 6 for 32MHz ref, 5 for 36MHz ref */
var SX125x_RX_BB_BW = 0         /* 0:750, 1:500, 2:375; 3:250 kHz SSB (default 1, max 3) */
var SX125x_RX_PLL_BW = 0        /* 0:75, 1:150, 2:225, 3:300 kHz (default 3, max 3) */
var SX125x_ADC_TEMP = 0         /* ADC temperature measurement mode (default 0) */
var SX125x_XOSC_GM_STARTUP = 13 /* (default 13) */
var SX125x_XOSC_DISABLE = 2     /* Disable of Xtal Oscillator blocks bit0:regulator, bit1:core(gm), bit2:amplifier */
var SX125x_32MHz_FRAC = uint32(15625)
var PLL_LOCK_MAX_ATTEMPTS = 5

func Lgw_setup_sx125x(c *os.File, lgw_spi_mux_mode, spi_mux_target, rf_chain, rf_clkout byte, rf_enable bool, rf_radio_type lgw_radio_type_e, freq_hz uint32) error {
	if rf_chain >= LGW_RF_CHAIN_NB {
		return fmt.Errorf("ERROR: INVALID RF_CHAIN\n")
	}

	/* Get version to identify SX1255/57 silicon revision */
	b, err := Sx125x_read(c, lgw_spi_mux_mode, spi_mux_target, rf_chain, 0x07)
	if err != nil {
		return err
	}
	fmt.Print("Note: SX125x #%d version register returned 0x%02X\n", rf_chain, b)

	/* General radio setup */
	if rf_clkout == rf_chain {
		err := Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x10, uint8(SX125x_TX_DAC_CLK_SEL+2))
		if err != nil {
			return err
		}
	} else {
		err := Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x10, uint8(SX125x_TX_DAC_CLK_SEL))
		if err != nil {
			return err
		}
	}

	switch rf_radio_type {
	case LGW_RADIO_TYPE_SX1255:
		err := Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x28, uint8(SX125x_XOSC_GM_STARTUP+SX125x_XOSC_DISABLE*16))
		if err != nil {
			return err
		}
	case LGW_RADIO_TYPE_SX1257:
		err := Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x26, uint8(SX125x_XOSC_GM_STARTUP+SX125x_XOSC_DISABLE*16))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("ERROR: UNEXPECTED VALUE %d FOR RADIO TYPE\n", rf_radio_type)
	}

	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x08, uint8(SX125x_TX_MIX_GAIN+SX125x_TX_DAC_GAIN*16))
	if err != nil {
		return err
	}
	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x0A, uint8(SX125x_TX_ANA_BW+SX125x_TX_PLL_BW*32))
	if err != nil {
		return err
	}
	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x0B, uint8(SX125x_TX_DAC_BW))
	if err != nil {
		return err
	}

	/* Rx gain and trim */
	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x0C, uint8(SX125x_LNA_ZIN+SX125x_RX_BB_GAIN*2+SX125x_RX_LNA_GAIN*32))
	if err != nil {
		return err
	}
	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x0D, uint8(SX125x_RX_BB_BW+SX125x_RX_ADC_TRIM*4+SX125x_RX_ADC_BW*32))
	if err != nil {
		return err
	}
	err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x0E, uint8(SX125x_ADC_TEMP+SX125x_RX_PLL_BW*2))
	if err != nil {
		return err
	}

	/* set RX PLL frequency */
	switch rf_radio_type {
	case LGW_RADIO_TYPE_SX1255:
		part_int := freq_hz / (SX125x_32MHz_FRAC << 7)                               /* integer part, gives the MSB */
		part_frac := ((freq_hz % (SX125x_32MHz_FRAC << 7)) << 9) / SX125x_32MHz_FRAC /* fractional part, gives middle part and LSB */

		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x01, 0xFF&uint8(part_int)) /* Most Significant Byte */
		if err != nil {
			return err
		}
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x02, 0xFF&uint8(part_frac>>8)) /* middle byte */
		if err != nil {
			return err
		}
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x03, 0xFF&uint8(part_frac)) /* Least Significant Byte */
		if err != nil {
			return err
		}
	case LGW_RADIO_TYPE_SX1257:
		part_int := freq_hz / (SX125x_32MHz_FRAC << 8)                                                /* integer part, gives the MSB */
		part_frac := ((freq_hz % (SX125x_32MHz_FRAC << 8)) << 8) / SX125x_32MHz_FRAC                  /* fractional part, gives middle part and LSB */
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x01, 0xFF&uint8(part_int)) /* Most Significant Byte */
		if err != nil {
			return err
		}
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x02, 0xFF&uint8(part_frac>>8)) /* middle byte */
		if err != nil {
			return err
		}
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x03, 0xFF&uint8(part_frac)) /* Least Significant Byte */
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("ERROR: UNEXPECTED VALUE %d FOR RADIO TYPE\n", rf_radio_type)
	}
	/* start and PLL lock */
	for cpt_attempts := 0; cpt_attempts < PLL_LOCK_MAX_ATTEMPTS; cpt_attempts++ {
		if cpt_attempts >= PLL_LOCK_MAX_ATTEMPTS {
			return fmt.Errorf("ERROR: FAIL TO LOCK PLL\n")
		}
		err := Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x00, 1) /* enable Xtal oscillator */
		if err != nil {
			return err
		}
		err = Sx125x_write(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x00, 3) /* Enable RX (PLL+FE) */
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
		val, err := Sx125x_read(c, rf_chain, lgw_spi_mux_mode, spi_mux_target, 0x11)
		if err != nil {
			return err
		}
		if (val & 0x02) != 0 {
			return err
		}
	}

	return nil
}

func Sx125x_write(c *os.File, channel, spi_mux_mode, spi_mux_target byte, addr, data uint8) error {
	var reg_add, reg_dat, reg_cs uint16

	/* checking input parameters */
	if channel >= LGW_RF_CHAIN_NB {
		return fmt.Errorf("ERROR: INVALID RF_CHAIN\n")
	}
	if addr >= 0x7F {
		return fmt.Errorf("ERROR: ADDRESS OUT OF RANGE\n")
	}

	/* selecting the target radio */
	switch channel {
	case 0:
		reg_add = LGW_SPI_RADIO_A__ADDR
		reg_dat = LGW_SPI_RADIO_A__DATA
		reg_cs = LGW_SPI_RADIO_A__CS
		break

	case 1:
		reg_add = LGW_SPI_RADIO_B__ADDR
		reg_dat = LGW_SPI_RADIO_B__DATA
		reg_cs = LGW_SPI_RADIO_B__CS
		break

	default:
		return fmt.Errorf("ERROR: UNEXPECTED VALUE %d IN SWITCH STATEMENT\n", channel)
	}

	/* SPI master data write procedure */
	err := Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 0)
	if err != nil {
		return err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_add, int32(0x80|addr)) /* MSB at 1 for write operation */
	if err != nil {
		return err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_dat, int32(data))
	if err != nil {
		return err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 1)
	if err != nil {
		return err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 0)
	if err != nil {
		return err
	}

	return nil
}

func Sx125x_read(c *os.File, spi_mux_mode, spi_mux_target byte, channel, addr byte) (byte, error) {
	var reg_add, reg_dat, reg_cs, reg_rb uint16

	/* checking input parameters */
	if channel >= LGW_RF_CHAIN_NB {
		return 0, fmt.Errorf("ERROR: INVALID RF_CHAIN\n")
	}
	if addr >= 0x7F {
		return 0, fmt.Errorf("ERROR: ADDRESS OUT OF RANGE\n")
	}

	/* selecting the target radio */
	switch channel {
	case 0:
		reg_add = LGW_SPI_RADIO_A__ADDR
		reg_dat = LGW_SPI_RADIO_A__DATA
		reg_cs = LGW_SPI_RADIO_A__CS
		reg_rb = LGW_SPI_RADIO_A__DATA_READBACK
		break

	case 1:
		reg_add = LGW_SPI_RADIO_B__ADDR
		reg_dat = LGW_SPI_RADIO_B__DATA
		reg_cs = LGW_SPI_RADIO_B__CS
		reg_rb = LGW_SPI_RADIO_B__DATA_READBACK
		break

	default:
		return 0, fmt.Errorf("ERROR: UNEXPECTED VALUE %d IN SWITCH STATEMENT\n", channel)
	}

	/* SPI master data read procedure */
	err := Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 0)
	if err != nil {
		return 0, err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_add, int32(addr)) /* MSB at 0 for read operation */
	if err != nil {
		return 0, err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_dat, 0)
	if err != nil {
		return 0, err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 1)
	if err != nil {
		return 0, err
	}
	err = Lgw_reg_w(c, spi_mux_mode, spi_mux_target, reg_cs, 0)
	if err != nil {
		return 0, err
	}
	read_value, err := Lgw_reg_r(c, spi_mux_mode, spi_mux_target, reg_rb)
	if err != nil {
		return 0, err
	}

	return byte(read_value), nil
}
