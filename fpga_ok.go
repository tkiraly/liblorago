package liblorago

import (
	"fmt"
	"os"
)

const (
	LGW_MIN_NOTCH_FREQ               = 126000 /* 126 KHz */
	LGW_MAX_NOTCH_FREQ               = 250000 /* 250 KHz */
	LGW_DEFAULT_NOTCH_FREQ           = 129000 /* 129 KHz */
	LGW_FPGA_SOFT_RESET              = 0
	LGW_FPGA_FEATURE                 = 1
	LGW_FPGA_LBT_INITIAL_FREQ        = 2
	LGW_FPGA_VERSION                 = 3
	LGW_FPGA_STATUS                  = 4
	LGW_FPGA_CTRL_FEATURE_START      = 5
	LGW_FPGA_CTRL_RADIO_RESET        = 6
	LGW_FPGA_CTRL_INPUT_SYNC_I       = 7
	LGW_FPGA_CTRL_INPUT_SYNC_Q       = 8
	LGW_FPGA_CTRL_OUTPUT_SYNC        = 9
	LGW_FPGA_CTRL_INVERT_IQ          = 10
	LGW_FPGA_CTRL_ACCESS_HISTO_MEM   = 11
	LGW_FPGA_CTRL_CLEAR_HISTO_MEM    = 12
	LGW_FPGA_HISTO_RAM_ADDR          = 13
	LGW_FPGA_HISTO_RAM_DATA          = 14
	LGW_FPGA_HISTO_NB_READ           = 15
	LGW_FPGA_LBT_TIMESTAMP_CH        = 16
	LGW_FPGA_LBT_TIMESTAMP_SELECT_CH = 17
	LGW_FPGA_LBT_CH0_FREQ_OFFSET     = 18
	LGW_FPGA_LBT_CH1_FREQ_OFFSET     = 19
	LGW_FPGA_LBT_CH2_FREQ_OFFSET     = 20
	LGW_FPGA_LBT_CH3_FREQ_OFFSET     = 21
	LGW_FPGA_LBT_CH4_FREQ_OFFSET     = 22
	LGW_FPGA_LBT_CH5_FREQ_OFFSET     = 23
	LGW_FPGA_LBT_CH6_FREQ_OFFSET     = 24
	LGW_FPGA_LBT_CH7_FREQ_OFFSET     = 25
	LGW_FPGA_SCAN_FREQ_OFFSET        = 26
	LGW_FPGA_LBT_SCAN_TIME_CH0       = 27
	LGW_FPGA_LBT_SCAN_TIME_CH1       = 28
	LGW_FPGA_LBT_SCAN_TIME_CH2       = 29
	LGW_FPGA_LBT_SCAN_TIME_CH3       = 30
	LGW_FPGA_LBT_SCAN_TIME_CH4       = 31
	LGW_FPGA_LBT_SCAN_TIME_CH5       = 32
	LGW_FPGA_LBT_SCAN_TIME_CH6       = 33
	LGW_FPGA_LBT_SCAN_TIME_CH7       = 34
	LGW_FPGA_RSSI_TARGET             = 35
	LGW_FPGA_HISTO_SCAN_FREQ         = 36
	LGW_FPGA_NOTCH_FREQ_OFFSET       = 37
	LGW_FPGA_TOTALREGS               = 38
)

var fpga_regs = [...]Lgw_reg_s{
	{-1, 0, 0, 0, 1, 0, 0},     /* SOFT_RESET */
	{-1, 0, 1, 0, 4, 1, 0},     /* FPGA_FEATURE */
	{-1, 0, 5, 0, 3, 1, 0},     /* LBT_INITIAL_FREQ */
	{-1, 1, 0, 0, 8, 1, 0},     /* VERSION */
	{-1, 2, 0, 0, 8, 1, 0},     /* FPGA_STATUS */
	{-1, 3, 0, 0, 1, 0, 0},     /* FPGA_CTRL_FEATURE_START */
	{-1, 3, 1, 0, 1, 0, 0},     /* FPGA_CTRL_RADIO_RESET */
	{-1, 3, 2, 0, 1, 0, 0},     /* FPGA_CTRL_INPUT_SYNC_I */
	{-1, 3, 3, 0, 1, 0, 0},     /* FPGA_CTRL_INPUT_SYNC_Q */
	{-1, 3, 4, 0, 1, 0, 0},     /* FPGA_CTRL_OUTPUT_SYNC */
	{-1, 3, 5, 0, 1, 0, 0},     /* FPGA_CTRL_INVERT_IQ */
	{-1, 3, 6, 0, 1, 0, 0},     /* FPGA_CTRL_ACCESS_HISTO_MEM */
	{-1, 3, 7, 0, 1, 0, 0},     /* FPGA_CTRL_CLEAR_HISTO_MEM */
	{-1, 4, 0, 0, 8, 0, 0},     /* HISTO_RAM_ADDR */
	{-1, 5, 0, 0, 8, 1, 0},     /* HISTO_RAM_DATA */
	{-1, 8, 0, 0, 16, 0, 1000}, /* HISTO_NB_READ */
	{-1, 14, 0, 0, 16, 1, 0},   /* LBT_TIMESTAMP_CH */
	{-1, 17, 0, 0, 4, 0, 0},    /* LBT_TIMESTAMP_SELECT_CH */
	{-1, 18, 0, 0, 8, 0, 0},    /* LBT_CH0_FREQ_OFFSET */
	{-1, 19, 0, 0, 8, 0, 0},    /* LBT_CH1_FREQ_OFFSET */
	{-1, 20, 0, 0, 8, 0, 0},    /* LBT_CH2_FREQ_OFFSET */
	{-1, 21, 0, 0, 8, 0, 0},    /* LBT_CH3_FREQ_OFFSET */
	{-1, 22, 0, 0, 8, 0, 0},    /* LBT_CH4_FREQ_OFFSET */
	{-1, 23, 0, 0, 8, 0, 0},    /* LBT_CH5_FREQ_OFFSET */
	{-1, 24, 0, 0, 8, 0, 0},    /* LBT_CH6_FREQ_OFFSET */
	{-1, 25, 0, 0, 8, 0, 0},    /* LBT_CH7_FREQ_OFFSET */
	{-1, 26, 0, 0, 8, 0, 0},    /* SCAN_FREQ_OFFSET */
	{-1, 28, 0, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH0 */
	{-1, 28, 1, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH1 */
	{-1, 28, 2, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH2 */
	{-1, 28, 3, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH3 */
	{-1, 28, 4, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH4 */
	{-1, 28, 5, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH5 */
	{-1, 28, 6, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH6 */
	{-1, 28, 7, 0, 1, 0, 0},    /* LBT_SCAN_TIME_CH7 */
	{-1, 30, 0, 0, 8, 0, 160},  /* RSSI_TARGET */
	{-1, 31, 0, 0, 24, 0, 0},   /* HISTO_SCAN_FREQ */
	{-1, 34, 0, 0, 6, 0, 0},    /* NOTCH_FREQ_OFFSET */
}

func lgw_fpga_get_tx_notch_delay(tx_notch_support, tx_notch_offset byte) float64 {
	if tx_notch_support != 1 {
		return 0
	}

	/* Notch filtering performed by FPGA adds a constant delay (group delay) that we need to compensate */
	tx_notch_delay := (31.25 * ((64 + float64(tx_notch_offset)) / 2)) / 1E3 /* 32MHz => 31.25ns */

	return tx_notch_delay
}

func Lgw_fpga_configure(f *os.File, tx_notch_freq uint32) error {

	/* Check input parameters */
	if (tx_notch_freq < LGW_MIN_NOTCH_FREQ) || (tx_notch_freq > LGW_MAX_NOTCH_FREQ) {
		return fmt.Errorf("WARNING: FPGA TX notch frequency is out of range (%u - [%u..%u]), setting it to default (%u)\n", tx_notch_freq, LGW_MIN_NOTCH_FREQ, LGW_MAX_NOTCH_FREQ, LGW_DEFAULT_NOTCH_FREQ)
		tx_notch_freq = LGW_DEFAULT_NOTCH_FREQ
	}

	/* Get supported FPGA features */
	fmt.Printf("INFO: FPGA supported features:")
	val, err := Lgw_fpga_reg_r(f, LGW_FPGA_FEATURE)
	if err != nil {
		return err
	}
	tx_notch_support := TAKE_N_BITS_FROM(byte(val), 0, 1)
	if tx_notch_support == 1 {
		fmt.Printf(" [TX filter] ")
	}
	spectral_scan_support := TAKE_N_BITS_FROM(byte(val), 1, 1)
	if spectral_scan_support == 1 {
		fmt.Printf(" [Spectral Scan] ")
	}
	lbt_support := TAKE_N_BITS_FROM(byte(val), 2, 1)
	if lbt_support == 1 {
		fmt.Printf(" [LBT] ")
	}
	fmt.Printf("\n")

	err = Lgw_fpga_reg_w(f, LGW_FPGA_CTRL_INPUT_SYNC_I, 1)
	if err != nil {
		return err
	}
	err = Lgw_fpga_reg_w(f, LGW_FPGA_CTRL_INPUT_SYNC_Q, 1)
	if err != nil {
		return err
	}
	err = Lgw_fpga_reg_w(f, LGW_FPGA_CTRL_OUTPUT_SYNC, 0)
	if err != nil {
		return err
	}
	/* Required for Semtech AP2 reference design */
	err = Lgw_fpga_reg_w(f, LGW_FPGA_CTRL_INVERT_IQ, 1)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to configure FPGA polarity\n")
	}

	/* Configure TX notch filter */
	if tx_notch_support == 1 {
		tx_notch_offset := int32((32E6 / (2 * tx_notch_freq)) - 64)
		err := Lgw_fpga_reg_w(f, LGW_FPGA_NOTCH_FREQ_OFFSET, tx_notch_offset)
		if err != nil {
			fmt.Errorf("ERROR: Failed to configure FPGA TX notch filter\n")
		}

		/* Readback to check that notch frequency is programmable */
		val, err := Lgw_fpga_reg_r(f, LGW_FPGA_NOTCH_FREQ_OFFSET)
		if err != nil {
			return fmt.Errorf("ERROR: Failed to read FPGA TX notch frequency\n")
		}
		if val != tx_notch_offset {
			return fmt.Errorf("WARNING: TX notch filter frequency is not programmable (check your FPGA image)\n")
		} else {
			fmt.Printf("INFO: TX notch filter frequency set to %d (%d)\n", tx_notch_freq, tx_notch_offset)
		}
	}

	return nil
}

/* Write to a register addressed by name */
func Lgw_fpga_reg_w(f *os.File, register_id uint16, reg_value int32) error {
	/* check input parameters */
	if register_id >= LGW_FPGA_TOTALREGS {
		return fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* get register struct from the struct array */
	r := fpga_regs[register_id]

	/* reject write to read-only registers */
	if r.rdon == 1 {
		return fmt.Errorf("ERROR: TRYING TO WRITE A READ-ONLY REGISTER\n")
	}

	err := reg_w_align32(f, LGW_SPI_MUX_MODE1, LGW_SPI_MUX_TARGET_FPGA, r, reg_value)
	if err != nil {
		return err
	}
	return nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Read to a register addressed by name */
func Lgw_fpga_reg_r(f *os.File, register_id uint16) (int32, error) {
	/* check input parameters */
	if register_id >= LGW_FPGA_TOTALREGS {
		return 0, fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* get register struct from the struct array */
	r := fpga_regs[register_id]

	b, err := reg_r_align32(f, LGW_SPI_MUX_MODE1, LGW_SPI_MUX_TARGET_FPGA, r)
	if err != nil {
		return 0, err
	}
	return b, err

}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Point to a register by name and do a burst write */
func Lgw_fpga_reg_wb(f *os.File, register_id uint16, data []byte) error {
	/* check input parameters */
	if len(data) == 0 {
		return fmt.Errorf("ERROR: BURST OF NULL LENGTH\n")
	}
	if register_id >= LGW_FPGA_TOTALREGS {
		return fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* get register struct from the struct array */
	r := fpga_regs[register_id]

	/* reject write to read-only registers */
	if r.rdon == 1 {
		return fmt.Errorf("ERROR: TRYING TO BURST WRITE A READ-ONLY REGISTER\n")
	}

	/* do the burst write */
	err := Lgw_spi_wb(f, LGW_SPI_MUX_MODE1, LGW_SPI_MUX_TARGET_FPGA, r.addr, data)
	if err != nil {
		return err
	}
	return nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Point to a register by name and do a burst read */
func Lgw_fpga_reg_rb(f *os.File, register_id, size uint16) ([]byte, error) {
	/* check input parameters */
	if size == 0 {
		return nil, fmt.Errorf("ERROR: BURST OF NULL LENGTH\n")
	}
	if register_id >= LGW_FPGA_TOTALREGS {
		return nil, fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* get register struct from the struct array */
	r := fpga_regs[register_id]

	/* do the burst read */
	b, err := Lgw_spi_rb(f, LGW_SPI_MUX_MODE1, LGW_SPI_MUX_TARGET_FPGA, r.addr, size)
	if err != nil {
		return nil, err
	}
	return b, nil
}
