package liblorago

import (
	"fmt"
	"log"
	"os"
)

const (
	PAGE_ADDR = 0x00
	PAGE_MASK = 0x03
)

var FPGA_VERSION []byte = []byte{31, 33} /* several versions could be supported */

type Lgw_reg_s struct {
	page int8  /*!< page containing the register (-1 for all pages) */
	addr uint8 /*!< base address of the register (7 bit) */
	offs uint8 /*!< position of the register LSB (between 0 to 7) */
	sign uint8 //bool  /*!< 1 indicates the register is signed (2 complem.) */
	leng uint8 /*!< number of bits in the register */
	rdon uint8 //bool  /*!< 1 indicates a read-only register */
	dflt int32 /*!< register default value */
}

var loregs = [...]Lgw_reg_s{
	{-1, 0, 0, 0, 2, 0, 0},     /* PAGE_REG */
	{-1, 0, 7, 0, 1, 0, 0},     /* SOFT_RESET */
	{-1, 1, 0, 0, 8, 1, 103},   /* VERSION */
	{-1, 2, 0, 0, 16, 0, 0},    /* RX_DATA_BUF_ADDR */
	{-1, 4, 0, 0, 8, 0, 0},     /* RX_DATA_BUF_DATA */
	{-1, 5, 0, 0, 8, 0, 0},     /* TX_DATA_BUF_ADDR */
	{-1, 6, 0, 0, 8, 0, 0},     /* TX_DATA_BUF_DATA */
	{-1, 7, 0, 0, 8, 0, 0},     /* CAPTURE_RAM_ADDR */
	{-1, 8, 0, 0, 8, 1, 0},     /* CAPTURE_RAM_DATA */
	{-1, 9, 0, 0, 8, 0, 0},     /* MCU_PROM_ADDR */
	{-1, 10, 0, 0, 8, 0, 0},    /* MCU_PROM_DATA */
	{-1, 11, 0, 0, 8, 0, 0},    /* RX_PACKET_DATA_FIFO_NUM_STORED */
	{-1, 12, 0, 0, 16, 1, 0},   /* RX_PACKET_DATA_FIFO_ADDR_POINTER */
	{-1, 14, 0, 0, 8, 1, 0},    /* RX_PACKET_DATA_FIFO_STATUS */
	{-1, 15, 0, 0, 8, 1, 0},    /* RX_PACKET_DATA_FIFO_PAYLOAD_SIZE */
	{-1, 16, 0, 0, 1, 0, 0},    /* MBWSSF_MODEM_ENABLE */
	{-1, 16, 1, 0, 1, 0, 0},    /* CONCENTRATOR_MODEM_ENABLE */
	{-1, 16, 2, 0, 1, 0, 0},    /* FSK_MODEM_ENABLE */
	{-1, 16, 3, 0, 1, 0, 0},    /* GLOBAL_EN */
	{-1, 17, 0, 0, 1, 0, 1},    /* CLK32M_EN */
	{-1, 17, 1, 0, 1, 0, 1},    /* CLKHS_EN */
	{-1, 18, 0, 0, 1, 0, 0},    /* START_BIST0 */
	{-1, 18, 1, 0, 1, 0, 0},    /* START_BIST1 */
	{-1, 18, 2, 0, 1, 0, 0},    /* CLEAR_BIST0 */
	{-1, 18, 3, 0, 1, 0, 0},    /* CLEAR_BIST1 */
	{-1, 19, 0, 0, 1, 1, 0},    /* BIST0_FINISHED */
	{-1, 19, 1, 0, 1, 1, 0},    /* BIST1_FINISHED */
	{-1, 20, 0, 0, 1, 1, 0},    /* MCU_AGC_PROG_RAM_BIST_STATUS */
	{-1, 20, 1, 0, 1, 1, 0},    /* MCU_ARB_PROG_RAM_BIST_STATUS */
	{-1, 20, 2, 0, 1, 1, 0},    /* CAPTURE_RAM_BIST_STATUS */
	{-1, 20, 3, 0, 1, 1, 0},    /* CHAN_FIR_RAM0_BIST_STATUS */
	{-1, 20, 4, 0, 1, 1, 0},    /* CHAN_FIR_RAM1_BIST_STATUS */
	{-1, 21, 0, 0, 1, 1, 0},    /* CORR0_RAM_BIST_STATUS */
	{-1, 21, 1, 0, 1, 1, 0},    /* CORR1_RAM_BIST_STATUS */
	{-1, 21, 2, 0, 1, 1, 0},    /* CORR2_RAM_BIST_STATUS */
	{-1, 21, 3, 0, 1, 1, 0},    /* CORR3_RAM_BIST_STATUS */
	{-1, 21, 4, 0, 1, 1, 0},    /* CORR4_RAM_BIST_STATUS */
	{-1, 21, 5, 0, 1, 1, 0},    /* CORR5_RAM_BIST_STATUS */
	{-1, 21, 6, 0, 1, 1, 0},    /* CORR6_RAM_BIST_STATUS */
	{-1, 21, 7, 0, 1, 1, 0},    /* CORR7_RAM_BIST_STATUS */
	{-1, 22, 0, 0, 1, 1, 0},    /* MODEM0_RAM0_BIST_STATUS */
	{-1, 22, 1, 0, 1, 1, 0},    /* MODEM1_RAM0_BIST_STATUS */
	{-1, 22, 2, 0, 1, 1, 0},    /* MODEM2_RAM0_BIST_STATUS */
	{-1, 22, 3, 0, 1, 1, 0},    /* MODEM3_RAM0_BIST_STATUS */
	{-1, 22, 4, 0, 1, 1, 0},    /* MODEM4_RAM0_BIST_STATUS */
	{-1, 22, 5, 0, 1, 1, 0},    /* MODEM5_RAM0_BIST_STATUS */
	{-1, 22, 6, 0, 1, 1, 0},    /* MODEM6_RAM0_BIST_STATUS */
	{-1, 22, 7, 0, 1, 1, 0},    /* MODEM7_RAM0_BIST_STATUS */
	{-1, 23, 0, 0, 1, 1, 0},    /* MODEM0_RAM1_BIST_STATUS */
	{-1, 23, 1, 0, 1, 1, 0},    /* MODEM1_RAM1_BIST_STATUS */
	{-1, 23, 2, 0, 1, 1, 0},    /* MODEM2_RAM1_BIST_STATUS */
	{-1, 23, 3, 0, 1, 1, 0},    /* MODEM3_RAM1_BIST_STATUS */
	{-1, 23, 4, 0, 1, 1, 0},    /* MODEM4_RAM1_BIST_STATUS */
	{-1, 23, 5, 0, 1, 1, 0},    /* MODEM5_RAM1_BIST_STATUS */
	{-1, 23, 6, 0, 1, 1, 0},    /* MODEM6_RAM1_BIST_STATUS */
	{-1, 23, 7, 0, 1, 1, 0},    /* MODEM7_RAM1_BIST_STATUS */
	{-1, 24, 0, 0, 1, 1, 0},    /* MODEM0_RAM2_BIST_STATUS */
	{-1, 24, 1, 0, 1, 1, 0},    /* MODEM1_RAM2_BIST_STATUS */
	{-1, 24, 2, 0, 1, 1, 0},    /* MODEM2_RAM2_BIST_STATUS */
	{-1, 24, 3, 0, 1, 1, 0},    /* MODEM3_RAM2_BIST_STATUS */
	{-1, 24, 4, 0, 1, 1, 0},    /* MODEM4_RAM2_BIST_STATUS */
	{-1, 24, 5, 0, 1, 1, 0},    /* MODEM5_RAM2_BIST_STATUS */
	{-1, 24, 6, 0, 1, 1, 0},    /* MODEM6_RAM2_BIST_STATUS */
	{-1, 24, 7, 0, 1, 1, 0},    /* MODEM7_RAM2_BIST_STATUS */
	{-1, 25, 0, 0, 1, 1, 0},    /* MODEM_MBWSSF_RAM0_BIST_STATUS */
	{-1, 25, 1, 0, 1, 1, 0},    /* MODEM_MBWSSF_RAM1_BIST_STATUS */
	{-1, 25, 2, 0, 1, 1, 0},    /* MODEM_MBWSSF_RAM2_BIST_STATUS */
	{-1, 26, 0, 0, 1, 1, 0},    /* MCU_AGC_DATA_RAM_BIST0_STATUS */
	{-1, 26, 1, 0, 1, 1, 0},    /* MCU_AGC_DATA_RAM_BIST1_STATUS */
	{-1, 26, 2, 0, 1, 1, 0},    /* MCU_ARB_DATA_RAM_BIST0_STATUS */
	{-1, 26, 3, 0, 1, 1, 0},    /* MCU_ARB_DATA_RAM_BIST1_STATUS */
	{-1, 26, 4, 0, 1, 1, 0},    /* TX_TOP_RAM_BIST0_STATUS */
	{-1, 26, 5, 0, 1, 1, 0},    /* TX_TOP_RAM_BIST1_STATUS */
	{-1, 26, 6, 0, 1, 1, 0},    /* DATA_MNGT_RAM_BIST0_STATUS */
	{-1, 26, 7, 0, 1, 1, 0},    /* DATA_MNGT_RAM_BIST1_STATUS */
	{-1, 27, 0, 0, 4, 0, 0},    /* GPIO_SELECT_INPUT */
	{-1, 28, 0, 0, 4, 0, 0},    /* GPIO_SELECT_OUTPUT */
	{-1, 29, 0, 0, 5, 0, 0},    /* GPIO_MODE */
	{-1, 30, 0, 0, 5, 1, 0},    /* GPIO_PIN_REG_IN */
	{-1, 31, 0, 0, 5, 0, 0},    /* GPIO_PIN_REG_OUT */
	{-1, 32, 0, 0, 8, 1, 0},    /* MCU_AGC_STATUS */
	{-1, 125, 0, 0, 8, 1, 0},   /* MCU_ARB_STATUS */
	{-1, 126, 0, 0, 8, 1, 1},   /* CHIP_ID */
	{-1, 127, 0, 0, 1, 0, 1},   /* EMERGENCY_FORCE_HOST_CTRL */
	{0, 33, 0, 0, 1, 0, 0},     /* RX_INVERT_IQ */
	{0, 33, 1, 0, 1, 0, 1},     /* MODEM_INVERT_IQ */
	{0, 33, 2, 0, 1, 0, 0},     /* MBWSSF_MODEM_INVERT_IQ */
	{0, 33, 3, 0, 1, 0, 0},     /* RX_EDGE_SELECT */
	{0, 33, 4, 0, 1, 0, 0},     /* MISC_RADIO_EN */
	{0, 33, 5, 0, 1, 0, 0},     /* FSK_MODEM_INVERT_IQ */
	{0, 34, 0, 0, 4, 0, 7},     /* FILTER_GAIN */
	{0, 35, 0, 0, 8, 0, 240},   /* RADIO_SELECT */
	{0, 36, 0, 1, 13, 0, -384}, /* IF_FREQ_0 */
	{0, 38, 0, 1, 13, 0, -128}, /* IF_FREQ_1 */
	{0, 40, 0, 1, 13, 0, 128},  /* IF_FREQ_2 */
	{0, 42, 0, 1, 13, 0, 384},  /* IF_FREQ_3 */
	{0, 44, 0, 1, 13, 0, -384}, /* IF_FREQ_4 */
	{0, 46, 0, 1, 13, 0, -128}, /* IF_FREQ_5 */
	{0, 48, 0, 1, 13, 0, 128},  /* IF_FREQ_6 */
	{0, 50, 0, 1, 13, 0, 384},  /* IF_FREQ_7 */
	{0, 52, 0, 1, 13, 0, 0},    /* IF_FREQ_8 */
	{0, 54, 0, 1, 13, 0, 0},    /* IF_FREQ_9 */
	{0, 64, 0, 0, 1, 0, 0},     /* CHANN_OVERRIDE_AGC_GAIN */
	{0, 64, 1, 0, 4, 0, 7},     /* CHANN_AGC_GAIN */
	{0, 65, 0, 0, 7, 0, 0},     /* CORR0_DETECT_EN */
	{0, 66, 0, 0, 7, 0, 0},     /* CORR1_DETECT_EN */
	{0, 67, 0, 0, 7, 0, 0},     /* CORR2_DETECT_EN */
	{0, 68, 0, 0, 7, 0, 0},     /* CORR3_DETECT_EN */
	{0, 69, 0, 0, 7, 0, 0},     /* CORR4_DETECT_EN */
	{0, 70, 0, 0, 7, 0, 0},     /* CORR5_DETECT_EN */
	{0, 71, 0, 0, 7, 0, 0},     /* CORR6_DETECT_EN */
	{0, 72, 0, 0, 7, 0, 0},     /* CORR7_DETECT_EN */
	{0, 73, 0, 0, 1, 0, 0},     /* CORR_SAME_PEAKS_OPTION_SF6 */
	{0, 73, 1, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF7 */
	{0, 73, 2, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF8 */
	{0, 73, 3, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF9 */
	{0, 73, 4, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF10 */
	{0, 73, 5, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF11 */
	{0, 73, 6, 0, 1, 0, 1},     /* CORR_SAME_PEAKS_OPTION_SF12 */
	{0, 74, 0, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF6 */
	{0, 74, 4, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF7 */
	{0, 75, 0, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF8 */
	{0, 75, 4, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF9 */
	{0, 76, 0, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF10 */
	{0, 76, 4, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF11 */
	{0, 77, 0, 0, 4, 0, 4},     /* CORR_SIG_NOISE_RATIO_SF12 */
	{0, 78, 0, 0, 4, 0, 4},     /* CORR_NUM_SAME_PEAK */
	{0, 78, 4, 0, 3, 0, 5},     /* CORR_MAC_GAIN */
	{0, 81, 0, 0, 12, 0, 0},    /* ADJUST_MODEM_START_OFFSET_RDX4 */
	{0, 83, 0, 0, 12, 0, 4092}, /* ADJUST_MODEM_START_OFFSET_SF12_RDX4 */
	{0, 85, 0, 0, 8, 0, 7},     /* DBG_CORR_SELECT_SF */
	{0, 86, 0, 0, 8, 0, 0},     /* DBG_CORR_SELECT_CHANNEL */
	{0, 87, 0, 0, 8, 1, 0},     /* DBG_DETECT_CPT */
	{0, 88, 0, 0, 8, 1, 0},     /* DBG_SYMB_CPT */
	{0, 89, 0, 0, 1, 0, 1},     /* CHIRP_INVERT_RX */
	{0, 89, 1, 0, 1, 0, 1},     /* DC_NOTCH_EN */
	{0, 90, 0, 0, 1, 0, 0},     /* IMPLICIT_CRC_EN */
	{0, 90, 1, 0, 3, 0, 0},     /* IMPLICIT_CODING_RATE */
	{0, 91, 0, 0, 8, 0, 0},     /* IMPLICIT_PAYLOAD_LENGHT */
	{0, 92, 0, 0, 8, 0, 29},    /* FREQ_TO_TIME_INVERT */
	{0, 93, 0, 0, 6, 0, 9},     /* FREQ_TO_TIME_DRIFT */
	{0, 94, 0, 0, 2, 0, 2},     /* PAYLOAD_FINE_TIMING_GAIN */
	{0, 94, 2, 0, 2, 0, 1},     /* PREAMBLE_FINE_TIMING_GAIN */
	{0, 94, 4, 0, 2, 0, 0},     /* TRACKING_INTEGRAL */
	{0, 95, 0, 0, 4, 0, 1},     /* FRAME_SYNCH_PEAK1_POS */
	{0, 95, 4, 0, 4, 0, 2},     /* FRAME_SYNCH_PEAK2_POS */
	{0, 96, 0, 0, 16, 0, 10},   /* PREAMBLE_SYMB1_NB */
	{0, 98, 0, 0, 1, 0, 1},     /* FRAME_SYNCH_GAIN */
	{0, 98, 1, 0, 1, 0, 1},     /* SYNCH_DETECT_TH */
	{0, 99, 0, 0, 4, 0, 8},     /* LLR_SCALE */
	{0, 99, 4, 0, 2, 0, 2},     /* SNR_AVG_CST */
	{0, 100, 0, 0, 7, 0, 0},    /* PPM_OFFSET */
	{0, 101, 0, 0, 8, 0, 255},  /* MAX_PAYLOAD_LEN */
	{0, 102, 0, 0, 1, 0, 1},    /* ONLY_CRC_EN */
	{0, 103, 0, 0, 8, 0, 0},    /* ZERO_PAD */
	{0, 104, 0, 0, 4, 0, 8},    /* DEC_GAIN_OFFSET */
	{0, 104, 4, 0, 4, 0, 7},    /* CHAN_GAIN_OFFSET */
	{0, 105, 1, 0, 1, 0, 1},    /* FORCE_HOST_RADIO_CTRL */
	{0, 105, 2, 0, 1, 0, 1},    /* FORCE_HOST_FE_CTRL */
	{0, 105, 3, 0, 1, 0, 1},    /* FORCE_DEC_FILTER_GAIN */
	{0, 106, 0, 0, 1, 0, 1},    /* MCU_RST_0 */
	{0, 106, 1, 0, 1, 0, 1},    /* MCU_RST_1 */
	{0, 106, 2, 0, 1, 0, 0},    /* MCU_SELECT_MUX_0 */
	{0, 106, 3, 0, 1, 0, 0},    /* MCU_SELECT_MUX_1 */
	{0, 106, 4, 0, 1, 1, 0},    /* MCU_CORRUPTION_DETECTED_0 */
	{0, 106, 5, 0, 1, 1, 0},    /* MCU_CORRUPTION_DETECTED_1 */
	{0, 106, 6, 0, 1, 0, 0},    /* MCU_SELECT_EDGE_0 */
	{0, 106, 7, 0, 1, 0, 0},    /* MCU_SELECT_EDGE_1 */
	{0, 107, 0, 0, 8, 0, 1},    /* CHANN_SELECT_RSSI */
	{0, 108, 0, 0, 8, 0, 32},   /* RSSI_BB_DEFAULT_VALUE */
	{0, 109, 0, 0, 8, 0, 100},  /* RSSI_DEC_DEFAULT_VALUE */
	{0, 110, 0, 0, 8, 0, 100},  /* RSSI_CHANN_DEFAULT_VALUE */
	{0, 111, 0, 0, 5, 0, 7},    /* RSSI_BB_FILTER_ALPHA */
	{0, 112, 0, 0, 5, 0, 5},    /* RSSI_DEC_FILTER_ALPHA */
	{0, 113, 0, 0, 5, 0, 8},    /* RSSI_CHANN_FILTER_ALPHA */
	{0, 114, 0, 0, 6, 0, 0},    /* IQ_MISMATCH_A_AMP_COEFF */
	{0, 115, 0, 0, 6, 0, 0},    /* IQ_MISMATCH_A_PHI_COEFF */
	{0, 116, 0, 0, 6, 0, 0},    /* IQ_MISMATCH_B_AMP_COEFF */
	{0, 116, 6, 0, 1, 0, 0},    /* IQ_MISMATCH_B_SEL_I */
	{0, 117, 0, 0, 6, 0, 0},    /* IQ_MISMATCH_B_PHI_COEFF */
	{1, 33, 0, 0, 1, 0, 0},     /* TX_TRIG_IMMEDIATE */
	{1, 33, 1, 0, 1, 0, 0},     /* TX_TRIG_DELAYED */
	{1, 33, 2, 0, 1, 0, 0},     /* TX_TRIG_GPS */
	{1, 34, 0, 0, 16, 0, 0},    /* TX_START_DELAY */
	{1, 36, 0, 0, 4, 0, 1},     /* TX_FRAME_SYNCH_PEAK1_POS */
	{1, 36, 4, 0, 4, 0, 2},     /* TX_FRAME_SYNCH_PEAK2_POS */
	{1, 37, 0, 0, 3, 0, 0},     /* TX_RAMP_DURATION */
	{1, 39, 0, 1, 8, 0, 0},     /* TX_OFFSET_I */
	{1, 40, 0, 1, 8, 0, 0},     /* TX_OFFSET_Q */
	{1, 41, 0, 0, 1, 0, 0},     /* TX_MODE */
	{1, 41, 1, 0, 4, 0, 0},     /* TX_ZERO_PAD */
	{1, 41, 5, 0, 1, 0, 0},     /* TX_EDGE_SELECT */
	{1, 41, 6, 0, 1, 0, 0},     /* TX_EDGE_SELECT_TOP */
	{1, 42, 0, 0, 2, 0, 0},     /* TX_GAIN */
	{1, 42, 2, 0, 3, 0, 5},     /* TX_CHIRP_LOW_PASS */
	{1, 42, 5, 0, 2, 0, 0},     /* TX_FCC_WIDEBAND */
	{1, 42, 7, 0, 1, 0, 1},     /* TX_SWAP_IQ */
	{1, 43, 0, 0, 1, 0, 0},     /* MBWSSF_IMPLICIT_HEADER */
	{1, 43, 1, 0, 1, 0, 0},     /* MBWSSF_IMPLICIT_CRC_EN */
	{1, 43, 2, 0, 3, 0, 0},     /* MBWSSF_IMPLICIT_CODING_RATE */
	{1, 44, 0, 0, 8, 0, 0},     /* MBWSSF_IMPLICIT_PAYLOAD_LENGHT */
	{1, 45, 0, 0, 1, 0, 1},     /* MBWSSF_AGC_FREEZE_ON_DETECT */
	{1, 46, 0, 0, 4, 0, 1},     /* MBWSSF_FRAME_SYNCH_PEAK1_POS */
	{1, 46, 4, 0, 4, 0, 2},     /* MBWSSF_FRAME_SYNCH_PEAK2_POS */
	{1, 47, 0, 0, 16, 0, 10},   /* MBWSSF_PREAMBLE_SYMB1_NB */
	{1, 49, 0, 0, 1, 0, 1},     /* MBWSSF_FRAME_SYNCH_GAIN */
	{1, 49, 1, 0, 1, 0, 1},     /* MBWSSF_SYNCH_DETECT_TH */
	{1, 50, 0, 0, 8, 0, 10},    /* MBWSSF_DETECT_MIN_SINGLE_PEAK */
	{1, 51, 0, 0, 3, 0, 3},     /* MBWSSF_DETECT_TRIG_SAME_PEAK_NB */
	{1, 52, 0, 0, 8, 0, 29},    /* MBWSSF_FREQ_TO_TIME_INVERT */
	{1, 53, 0, 0, 6, 0, 36},    /* MBWSSF_FREQ_TO_TIME_DRIFT */
	{1, 54, 0, 0, 12, 0, 0},    /* MBWSSF_PPM_CORRECTION */
	{1, 56, 0, 0, 2, 0, 2},     /* MBWSSF_PAYLOAD_FINE_TIMING_GAIN */
	{1, 56, 2, 0, 2, 0, 1},     /* MBWSSF_PREAMBLE_FINE_TIMING_GAIN */
	{1, 56, 4, 0, 2, 0, 0},     /* MBWSSF_TRACKING_INTEGRAL */
	{1, 57, 0, 0, 8, 0, 0},     /* MBWSSF_ZERO_PAD */
	{1, 58, 0, 0, 2, 0, 0},     /* MBWSSF_MODEM_BW */
	{1, 58, 2, 0, 1, 0, 0},     /* MBWSSF_RADIO_SELECT */
	{1, 58, 3, 0, 1, 0, 1},     /* MBWSSF_RX_CHIRP_INVERT */
	{1, 59, 0, 0, 4, 0, 8},     /* MBWSSF_LLR_SCALE */
	{1, 59, 4, 0, 2, 0, 3},     /* MBWSSF_SNR_AVG_CST */
	{1, 59, 6, 0, 1, 0, 0},     /* MBWSSF_PPM_OFFSET */
	{1, 60, 0, 0, 4, 0, 7},     /* MBWSSF_RATE_SF */
	{1, 60, 4, 0, 1, 0, 1},     /* MBWSSF_ONLY_CRC_EN */
	{1, 61, 0, 0, 8, 0, 255},   /* MBWSSF_MAX_PAYLOAD_LEN */
	{1, 62, 0, 0, 8, 1, 128},   /* TX_STATUS */
	{1, 63, 0, 0, 3, 0, 0},     /* FSK_CH_BW_EXPO */
	{1, 63, 3, 0, 3, 0, 0},     /* FSK_RSSI_LENGTH */
	{1, 63, 6, 0, 1, 0, 0},     /* FSK_RX_INVERT */
	{1, 63, 7, 0, 1, 0, 0},     /* FSK_PKT_MODE */
	{1, 64, 0, 0, 3, 0, 0},     /* FSK_PSIZE */
	{1, 64, 3, 0, 1, 0, 0},     /* FSK_CRC_EN */
	{1, 64, 4, 0, 2, 0, 0},     /* FSK_DCFREE_ENC */
	{1, 64, 6, 0, 1, 0, 0},     /* FSK_CRC_IBM */
	{1, 65, 0, 0, 5, 0, 0},     /* FSK_ERROR_OSR_TOL */
	{1, 65, 7, 0, 1, 0, 0},     /* FSK_RADIO_SELECT */
	{1, 66, 0, 0, 16, 0, 0},    /* FSK_BR_RATIO */
	{1, 68, 0, 0, 32, 0, 0},    /* FSK_REF_PATTERN_LSB */
	{1, 72, 0, 0, 32, 0, 0},    /* FSK_REF_PATTERN_MSB */
	{1, 76, 0, 0, 8, 0, 0},     /* FSK_PKT_LENGTH */
	{1, 77, 0, 0, 1, 0, 1},     /* FSK_TX_GAUSSIAN_EN */
	{1, 77, 1, 0, 2, 0, 0},     /* FSK_TX_GAUSSIAN_SELECT_BT */
	{1, 77, 3, 0, 1, 0, 1},     /* FSK_TX_PATTERN_EN */
	{1, 77, 4, 0, 1, 0, 0},     /* FSK_TX_PREAMBLE_SEQ */
	{1, 77, 5, 0, 3, 0, 0},     /* FSK_TX_PSIZE */
	{1, 80, 0, 0, 8, 0, 0},     /* FSK_NODE_ADRS */
	{1, 81, 0, 0, 8, 0, 0},     /* FSK_BROADCAST */
	{1, 82, 0, 0, 1, 0, 1},     /* FSK_AUTO_AFC_ON */
	{1, 83, 0, 0, 10, 0, 0},    /* FSK_PATTERN_TIMEOUT_CFG */
	{2, 33, 0, 0, 8, 0, 0},     /* SPI_RADIO_A__DATA */
	{2, 34, 0, 0, 8, 1, 0},     /* SPI_RADIO_A__DATA_READBACK */
	{2, 35, 0, 0, 8, 0, 0},     /* SPI_RADIO_A__ADDR */
	{2, 37, 0, 0, 1, 0, 0},     /* SPI_RADIO_A__CS */
	{2, 38, 0, 0, 8, 0, 0},     /* SPI_RADIO_B__DATA */
	{2, 39, 0, 0, 8, 1, 0},     /* SPI_RADIO_B__DATA_READBACK */
	{2, 40, 0, 0, 8, 0, 0},     /* SPI_RADIO_B__ADDR */
	{2, 42, 0, 0, 1, 0, 0},     /* SPI_RADIO_B__CS */
	{2, 43, 0, 0, 1, 0, 0},     /* RADIO_A_EN */
	{2, 43, 1, 0, 1, 0, 0},     /* RADIO_B_EN */
	{2, 43, 2, 0, 1, 0, 1},     /* RADIO_RST */
	{2, 43, 3, 0, 1, 0, 0},     /* LNA_A_EN */
	{2, 43, 4, 0, 1, 0, 0},     /* PA_A_EN */
	{2, 43, 5, 0, 1, 0, 0},     /* LNA_B_EN */
	{2, 43, 6, 0, 1, 0, 0},     /* PA_B_EN */
	{2, 44, 0, 0, 2, 0, 0},     /* PA_GAIN */
	{2, 45, 0, 0, 4, 0, 2},     /* LNA_A_CTRL_LUT */
	{2, 45, 4, 0, 4, 0, 4},     /* PA_A_CTRL_LUT */
	{2, 46, 0, 0, 4, 0, 2},     /* LNA_B_CTRL_LUT */
	{2, 46, 4, 0, 4, 0, 4},     /* PA_B_CTRL_LUT */
	{2, 47, 0, 0, 5, 0, 0},     /* CAPTURE_SOURCE */
	{2, 47, 5, 0, 1, 0, 0},     /* CAPTURE_START */
	{2, 47, 6, 0, 1, 0, 0},     /* CAPTURE_FORCE_TRIGGER */
	{2, 47, 7, 0, 1, 0, 0},     /* CAPTURE_WRAP */
	{2, 48, 0, 0, 16, 0, 0},    /* CAPTURE_PERIOD */
	{2, 51, 0, 0, 8, 1, 0},     /* MODEM_STATUS */
	{2, 52, 0, 0, 8, 1, 0},     /* VALID_HEADER_COUNTER_0 */
	{2, 54, 0, 0, 8, 1, 0},     /* VALID_PACKET_COUNTER_0 */
	{2, 56, 0, 0, 8, 1, 0},     /* VALID_HEADER_COUNTER_MBWSSF */
	{2, 57, 0, 0, 8, 1, 0},     /* VALID_HEADER_COUNTER_FSK */
	{2, 58, 0, 0, 8, 1, 0},     /* VALID_PACKET_COUNTER_MBWSSF */
	{2, 59, 0, 0, 8, 1, 0},     /* VALID_PACKET_COUNTER_FSK */
	{2, 60, 0, 0, 8, 1, 0},     /* CHANN_RSSI */
	{2, 61, 0, 0, 8, 1, 0},     /* BB_RSSI */
	{2, 62, 0, 0, 8, 1, 0},     /* DEC_RSSI */
	{2, 63, 0, 0, 8, 1, 0},     /* DBG_MCU_DATA */
	{2, 64, 0, 0, 8, 1, 0},     /* DBG_ARB_MCU_RAM_DATA */
	{2, 65, 0, 0, 8, 1, 0},     /* DBG_AGC_MCU_RAM_DATA */
	{2, 66, 0, 0, 16, 1, 0},    /* NEXT_PACKET_CNT */
	{2, 68, 0, 0, 16, 1, 0},    /* ADDR_CAPTURE_COUNT */
	{2, 70, 0, 0, 32, 1, 0},    /* TIMESTAMP */
	{2, 74, 0, 0, 4, 1, 0},     /* DBG_CHANN0_GAIN */
	{2, 74, 4, 0, 4, 1, 0},     /* DBG_CHANN1_GAIN */
	{2, 75, 0, 0, 4, 1, 0},     /* DBG_CHANN2_GAIN */
	{2, 75, 4, 0, 4, 1, 0},     /* DBG_CHANN3_GAIN */
	{2, 76, 0, 0, 4, 1, 0},     /* DBG_CHANN4_GAIN */
	{2, 76, 4, 0, 4, 1, 0},     /* DBG_CHANN5_GAIN */
	{2, 77, 0, 0, 4, 1, 0},     /* DBG_CHANN6_GAIN */
	{2, 77, 4, 0, 4, 1, 0},     /* DBG_CHANN7_GAIN */
	{2, 78, 0, 0, 4, 1, 0},     /* DBG_DEC_FILT_GAIN */
	{2, 79, 0, 0, 3, 1, 0},     /* SPI_DATA_FIFO_PTR */
	{2, 79, 3, 0, 3, 1, 0},     /* PACKET_DATA_FIFO_PTR */
	{2, 80, 0, 0, 8, 0, 0},     /* DBG_ARB_MCU_RAM_ADDR */
	{2, 81, 0, 0, 8, 0, 0},     /* DBG_AGC_MCU_RAM_ADDR */
	{2, 82, 0, 0, 1, 0, 0},     /* SPI_MASTER_CHIP_SELECT_POLARITY */
	{2, 82, 1, 0, 1, 0, 0},     /* SPI_MASTER_CPOL */
	{2, 82, 2, 0, 1, 0, 0},     /* SPI_MASTER_CPHA */
	{2, 83, 0, 0, 1, 0, 0},     /* SIG_GEN_ANALYSER_MUX_SEL */
	{2, 84, 0, 0, 1, 0, 0},     /* SIG_GEN_EN */
	{2, 84, 1, 0, 1, 0, 0},     /* SIG_ANALYSER_EN */
	{2, 84, 2, 0, 2, 0, 0},     /* SIG_ANALYSER_AVG_LEN */
	{2, 84, 4, 0, 3, 0, 0},     /* SIG_ANALYSER_PRECISION */
	{2, 84, 7, 0, 1, 1, 0},     /* SIG_ANALYSER_VALID_OUT */
	{2, 85, 0, 0, 8, 0, 0},     /* SIG_GEN_FREQ */
	{2, 86, 0, 0, 8, 0, 0},     /* SIG_ANALYSER_FREQ */
	{2, 87, 0, 0, 8, 1, 0},     /* SIG_ANALYSER_I_OUT */
	{2, 88, 0, 0, 8, 1, 0},     /* SIG_ANALYSER_Q_OUT */
	{2, 89, 0, 0, 1, 0, 0},     /* GPS_EN */
	{2, 89, 1, 0, 1, 0, 1},     /* GPS_POL */
	{2, 90, 0, 1, 8, 0, 0},     /* SW_TEST_REG1 */
	{2, 91, 2, 1, 6, 0, 0},     /* SW_TEST_REG2 */
	{2, 92, 0, 1, 16, 0, 0},    /* SW_TEST_REG3 */
	{2, 94, 0, 0, 4, 1, 0},     /* DATA_MNGT_STATUS */
	{2, 95, 0, 0, 5, 1, 0},     /* DATA_MNGT_CPT_FRAME_ALLOCATED */
	{2, 96, 0, 0, 5, 1, 0},     /* DATA_MNGT_CPT_FRAME_FINISHED */
	{2, 97, 0, 0, 5, 1, 0},     /* DATA_MNGT_CPT_FRAME_READEN */
	{1, 33, 0, 0, 8, 0, 0},     /* TX_TRIG_ALL (alias) */
}

const (
	LGW_PAGE_REG                            uint16 = 0
	LGW_SOFT_RESET                          uint16 = 1
	LGW_VERSION                             uint16 = 2
	LGW_RX_DATA_BUF_ADDR                    uint16 = 3
	LGW_RX_DATA_BUF_DATA                    uint16 = 4
	LGW_TX_DATA_BUF_ADDR                    uint16 = 5
	LGW_TX_DATA_BUF_DATA                    uint16 = 6
	LGW_CAPTURE_RAM_ADDR                    uint16 = 7
	LGW_CAPTURE_RAM_DATA                    uint16 = 8
	LGW_MCU_PROM_ADDR                       uint16 = 9
	LGW_MCU_PROM_DATA                       uint16 = 10
	LGW_RX_PACKET_DATA_FIFO_NUM_STORED      uint16 = 11
	LGW_RX_PACKET_DATA_FIFO_ADDR_POINTER    uint16 = 12
	LGW_RX_PACKET_DATA_FIFO_STATUS          uint16 = 13
	LGW_RX_PACKET_DATA_FIFO_PAYLOAD_SIZE    uint16 = 14
	LGW_MBWSSF_MODEM_ENABLE                 uint16 = 15
	LGW_CONCENTRATOR_MODEM_ENABLE           uint16 = 16
	LGW_FSK_MODEM_ENABLE                    uint16 = 17
	LGW_GLOBAL_EN                           uint16 = 18
	LGW_CLK32M_EN                           uint16 = 19
	LGW_CLKHS_EN                            uint16 = 20
	LGW_START_BIST0                         uint16 = 21
	LGW_START_BIST1                         uint16 = 22
	LGW_CLEAR_BIST0                         uint16 = 23
	LGW_CLEAR_BIST1                         uint16 = 24
	LGW_BIST0_FINISHED                      uint16 = 25
	LGW_BIST1_FINISHED                      uint16 = 26
	LGW_MCU_AGC_PROG_RAM_BIST_STATUS        uint16 = 27
	LGW_MCU_ARB_PROG_RAM_BIST_STATUS        uint16 = 28
	LGW_CAPTURE_RAM_BIST_STATUS             uint16 = 29
	LGW_CHAN_FIR_RAM0_BIST_STATUS           uint16 = 30
	LGW_CHAN_FIR_RAM1_BIST_STATUS           uint16 = 31
	LGW_CORR0_RAM_BIST_STATUS               uint16 = 32
	LGW_CORR1_RAM_BIST_STATUS               uint16 = 33
	LGW_CORR2_RAM_BIST_STATUS               uint16 = 34
	LGW_CORR3_RAM_BIST_STATUS               uint16 = 35
	LGW_CORR4_RAM_BIST_STATUS               uint16 = 36
	LGW_CORR5_RAM_BIST_STATUS               uint16 = 37
	LGW_CORR6_RAM_BIST_STATUS               uint16 = 38
	LGW_CORR7_RAM_BIST_STATUS               uint16 = 39
	LGW_MODEM0_RAM0_BIST_STATUS             uint16 = 40
	LGW_MODEM1_RAM0_BIST_STATUS             uint16 = 41
	LGW_MODEM2_RAM0_BIST_STATUS             uint16 = 42
	LGW_MODEM3_RAM0_BIST_STATUS             uint16 = 43
	LGW_MODEM4_RAM0_BIST_STATUS             uint16 = 44
	LGW_MODEM5_RAM0_BIST_STATUS             uint16 = 45
	LGW_MODEM6_RAM0_BIST_STATUS             uint16 = 46
	LGW_MODEM7_RAM0_BIST_STATUS             uint16 = 47
	LGW_MODEM0_RAM1_BIST_STATUS             uint16 = 48
	LGW_MODEM1_RAM1_BIST_STATUS             uint16 = 49
	LGW_MODEM2_RAM1_BIST_STATUS             uint16 = 50
	LGW_MODEM3_RAM1_BIST_STATUS             uint16 = 51
	LGW_MODEM4_RAM1_BIST_STATUS             uint16 = 52
	LGW_MODEM5_RAM1_BIST_STATUS             uint16 = 53
	LGW_MODEM6_RAM1_BIST_STATUS             uint16 = 54
	LGW_MODEM7_RAM1_BIST_STATUS             uint16 = 55
	LGW_MODEM0_RAM2_BIST_STATUS             uint16 = 56
	LGW_MODEM1_RAM2_BIST_STATUS             uint16 = 57
	LGW_MODEM2_RAM2_BIST_STATUS             uint16 = 58
	LGW_MODEM3_RAM2_BIST_STATUS             uint16 = 59
	LGW_MODEM4_RAM2_BIST_STATUS             uint16 = 60
	LGW_MODEM5_RAM2_BIST_STATUS             uint16 = 61
	LGW_MODEM6_RAM2_BIST_STATUS             uint16 = 62
	LGW_MODEM7_RAM2_BIST_STATUS             uint16 = 63
	LGW_MODEM_MBWSSF_RAM0_BIST_STATUS       uint16 = 64
	LGW_MODEM_MBWSSF_RAM1_BIST_STATUS       uint16 = 65
	LGW_MODEM_MBWSSF_RAM2_BIST_STATUS       uint16 = 66
	LGW_MCU_AGC_DATA_RAM_BIST0_STATUS       uint16 = 67
	LGW_MCU_AGC_DATA_RAM_BIST1_STATUS       uint16 = 68
	LGW_MCU_ARB_DATA_RAM_BIST0_STATUS       uint16 = 69
	LGW_MCU_ARB_DATA_RAM_BIST1_STATUS       uint16 = 70
	LGW_TX_TOP_RAM_BIST0_STATUS             uint16 = 71
	LGW_TX_TOP_RAM_BIST1_STATUS             uint16 = 72
	LGW_DATA_MNGT_RAM_BIST0_STATUS          uint16 = 73
	LGW_DATA_MNGT_RAM_BIST1_STATUS          uint16 = 74
	LGW_GPIO_SELECT_INPUT                   uint16 = 75
	LGW_GPIO_SELECT_OUTPUT                  uint16 = 76
	LGW_GPIO_MODE                           uint16 = 77
	LGW_GPIO_PIN_REG_IN                     uint16 = 78
	LGW_GPIO_PIN_REG_OUT                    uint16 = 79
	LGW_MCU_AGC_STATUS                      uint16 = 80
	LGW_MCU_ARB_STATUS                      uint16 = 81
	LGW_CHIP_ID                             uint16 = 82
	LGW_EMERGENCY_FORCE_HOST_CTRL           uint16 = 83
	LGW_RX_INVERT_IQ                        uint16 = 84
	LGW_MODEM_INVERT_IQ                     uint16 = 85
	LGW_MBWSSF_MODEM_INVERT_IQ              uint16 = 86
	LGW_RX_EDGE_SELECT                      uint16 = 87
	LGW_MISC_RADIO_EN                       uint16 = 88
	LGW_FSK_MODEM_INVERT_IQ                 uint16 = 89
	LGW_FILTER_GAIN                         uint16 = 90
	LGW_RADIO_SELECT                        uint16 = 91
	LGW_IF_FREQ_0                           uint16 = 92
	LGW_IF_FREQ_1                           uint16 = 93
	LGW_IF_FREQ_2                           uint16 = 94
	LGW_IF_FREQ_3                           uint16 = 95
	LGW_IF_FREQ_4                           uint16 = 96
	LGW_IF_FREQ_5                           uint16 = 97
	LGW_IF_FREQ_6                           uint16 = 98
	LGW_IF_FREQ_7                           uint16 = 99
	LGW_IF_FREQ_8                           uint16 = 100
	LGW_IF_FREQ_9                           uint16 = 101
	LGW_CHANN_OVERRIDE_AGC_GAIN             uint16 = 102
	LGW_CHANN_AGC_GAIN                      uint16 = 103
	LGW_CORR0_DETECT_EN                     uint16 = 104
	LGW_CORR1_DETECT_EN                     uint16 = 105
	LGW_CORR2_DETECT_EN                     uint16 = 106
	LGW_CORR3_DETECT_EN                     uint16 = 107
	LGW_CORR4_DETECT_EN                     uint16 = 108
	LGW_CORR5_DETECT_EN                     uint16 = 109
	LGW_CORR6_DETECT_EN                     uint16 = 110
	LGW_CORR7_DETECT_EN                     uint16 = 111
	LGW_CORR_SAME_PEAKS_OPTION_SF6          uint16 = 112
	LGW_CORR_SAME_PEAKS_OPTION_SF7          uint16 = 113
	LGW_CORR_SAME_PEAKS_OPTION_SF8          uint16 = 114
	LGW_CORR_SAME_PEAKS_OPTION_SF9          uint16 = 115
	LGW_CORR_SAME_PEAKS_OPTION_SF10         uint16 = 116
	LGW_CORR_SAME_PEAKS_OPTION_SF11         uint16 = 117
	LGW_CORR_SAME_PEAKS_OPTION_SF12         uint16 = 118
	LGW_CORR_SIG_NOISE_RATIO_SF6            uint16 = 119
	LGW_CORR_SIG_NOISE_RATIO_SF7            uint16 = 120
	LGW_CORR_SIG_NOISE_RATIO_SF8            uint16 = 121
	LGW_CORR_SIG_NOISE_RATIO_SF9            uint16 = 122
	LGW_CORR_SIG_NOISE_RATIO_SF10           uint16 = 123
	LGW_CORR_SIG_NOISE_RATIO_SF11           uint16 = 124
	LGW_CORR_SIG_NOISE_RATIO_SF12           uint16 = 125
	LGW_CORR_NUM_SAME_PEAK                  uint16 = 126
	LGW_CORR_MAC_GAIN                       uint16 = 127
	LGW_ADJUST_MODEM_START_OFFSET_RDX4      uint16 = 128
	LGW_ADJUST_MODEM_START_OFFSET_SF12_RDX4 uint16 = 129
	LGW_DBG_CORR_SELECT_SF                  uint16 = 130
	LGW_DBG_CORR_SELECT_CHANNEL             uint16 = 131
	LGW_DBG_DETECT_CPT                      uint16 = 132
	LGW_DBG_SYMB_CPT                        uint16 = 133
	LGW_CHIRP_INVERT_RX                     uint16 = 134
	LGW_DC_NOTCH_EN                         uint16 = 135
	LGW_IMPLICIT_CRC_EN                     uint16 = 136
	LGW_IMPLICIT_CODING_RATE                uint16 = 137
	LGW_IMPLICIT_PAYLOAD_LENGHT             uint16 = 138
	LGW_FREQ_TO_TIME_INVERT                 uint16 = 139
	LGW_FREQ_TO_TIME_DRIFT                  uint16 = 140
	LGW_PAYLOAD_FINE_TIMING_GAIN            uint16 = 141
	LGW_PREAMBLE_FINE_TIMING_GAIN           uint16 = 142
	LGW_TRACKING_INTEGRAL                   uint16 = 143
	LGW_FRAME_SYNCH_PEAK1_POS               uint16 = 144
	LGW_FRAME_SYNCH_PEAK2_POS               uint16 = 145
	LGW_PREAMBLE_SYMB1_NB                   uint16 = 146
	LGW_FRAME_SYNCH_GAIN                    uint16 = 147
	LGW_SYNCH_DETECT_TH                     uint16 = 148
	LGW_LLR_SCALE                           uint16 = 149
	LGW_SNR_AVG_CST                         uint16 = 150
	LGW_PPM_OFFSET                          uint16 = 151
	LGW_MAX_PAYLOAD_LEN                     uint16 = 152
	LGW_ONLY_CRC_EN                         uint16 = 153
	LGW_ZERO_PAD                            uint16 = 154
	LGW_DEC_GAIN_OFFSET                     uint16 = 155
	LGW_CHAN_GAIN_OFFSET                    uint16 = 156
	LGW_FORCE_HOST_RADIO_CTRL               uint16 = 157
	LGW_FORCE_HOST_FE_CTRL                  uint16 = 158
	LGW_FORCE_DEC_FILTER_GAIN               uint16 = 159
	LGW_MCU_RST_0                           uint16 = 160
	LGW_MCU_RST_1                           uint16 = 161
	LGW_MCU_SELECT_MUX_0                    uint16 = 162
	LGW_MCU_SELECT_MUX_1                    uint16 = 163
	LGW_MCU_CORRUPTION_DETECTED_0           uint16 = 164
	LGW_MCU_CORRUPTION_DETECTED_1           uint16 = 165
	LGW_MCU_SELECT_EDGE_0                   uint16 = 166
	LGW_MCU_SELECT_EDGE_1                   uint16 = 167
	LGW_CHANN_SELECT_RSSI                   uint16 = 168
	LGW_RSSI_BB_DEFAULT_VALUE               uint16 = 169
	LGW_RSSI_DEC_DEFAULT_VALUE              uint16 = 170
	LGW_RSSI_CHANN_DEFAULT_VALUE            uint16 = 171
	LGW_RSSI_BB_FILTER_ALPHA                uint16 = 172
	LGW_RSSI_DEC_FILTER_ALPHA               uint16 = 173
	LGW_RSSI_CHANN_FILTER_ALPHA             uint16 = 174
	LGW_IQ_MISMATCH_A_AMP_COEFF             uint16 = 175
	LGW_IQ_MISMATCH_A_PHI_COEFF             uint16 = 176
	LGW_IQ_MISMATCH_B_AMP_COEFF             uint16 = 177
	LGW_IQ_MISMATCH_B_SEL_I                 uint16 = 178
	LGW_IQ_MISMATCH_B_PHI_COEFF             uint16 = 179
	LGW_TX_TRIG_IMMEDIATE                   uint16 = 180
	LGW_TX_TRIG_DELAYED                     uint16 = 181
	LGW_TX_TRIG_GPS                         uint16 = 182
	LGW_TX_START_DELAY                      uint16 = 183
	LGW_TX_FRAME_SYNCH_PEAK1_POS            uint16 = 184
	LGW_TX_FRAME_SYNCH_PEAK2_POS            uint16 = 185
	LGW_TX_RAMP_DURATION                    uint16 = 186
	LGW_TX_OFFSET_I                         uint16 = 187
	LGW_TX_OFFSET_Q                         uint16 = 188
	LGW_TX_MODE                             uint16 = 189
	LGW_TX_ZERO_PAD                         uint16 = 190
	LGW_TX_EDGE_SELECT                      uint16 = 191
	LGW_TX_EDGE_SELECT_TOP                  uint16 = 192
	LGW_TX_GAIN                             uint16 = 193
	LGW_TX_CHIRP_LOW_PASS                   uint16 = 194
	LGW_TX_FCC_WIDEBAND                     uint16 = 195
	LGW_TX_SWAP_IQ                          uint16 = 196
	LGW_MBWSSF_IMPLICIT_HEADER              uint16 = 197
	LGW_MBWSSF_IMPLICIT_CRC_EN              uint16 = 198
	LGW_MBWSSF_IMPLICIT_CODING_RATE         uint16 = 199
	LGW_MBWSSF_IMPLICIT_PAYLOAD_LENGHT      uint16 = 200
	LGW_MBWSSF_AGC_FREEZE_ON_DETECT         uint16 = 201
	LGW_MBWSSF_FRAME_SYNCH_PEAK1_POS        uint16 = 202
	LGW_MBWSSF_FRAME_SYNCH_PEAK2_POS        uint16 = 203
	LGW_MBWSSF_PREAMBLE_SYMB1_NB            uint16 = 204
	LGW_MBWSSF_FRAME_SYNCH_GAIN             uint16 = 205
	LGW_MBWSSF_SYNCH_DETECT_TH              uint16 = 206
	LGW_MBWSSF_DETECT_MIN_SINGLE_PEAK       uint16 = 207
	LGW_MBWSSF_DETECT_TRIG_SAME_PEAK_NB     uint16 = 208
	LGW_MBWSSF_FREQ_TO_TIME_INVERT          uint16 = 209
	LGW_MBWSSF_FREQ_TO_TIME_DRIFT           uint16 = 210
	LGW_MBWSSF_PPM_CORRECTION               uint16 = 211
	LGW_MBWSSF_PAYLOAD_FINE_TIMING_GAIN     uint16 = 212
	LGW_MBWSSF_PREAMBLE_FINE_TIMING_GAIN    uint16 = 213
	LGW_MBWSSF_TRACKING_INTEGRAL            uint16 = 214
	LGW_MBWSSF_ZERO_PAD                     uint16 = 215
	LGW_MBWSSF_MODEM_BW                     uint16 = 216
	LGW_MBWSSF_RADIO_SELECT                 uint16 = 217
	LGW_MBWSSF_RX_CHIRP_INVERT              uint16 = 218
	LGW_MBWSSF_LLR_SCALE                    uint16 = 219
	LGW_MBWSSF_SNR_AVG_CST                  uint16 = 220
	LGW_MBWSSF_PPM_OFFSET                   uint16 = 221
	LGW_MBWSSF_RATE_SF                      uint16 = 222
	LGW_MBWSSF_ONLY_CRC_EN                  uint16 = 223
	LGW_MBWSSF_MAX_PAYLOAD_LEN              uint16 = 224
	LGW_TX_STATUS                           uint16 = 225
	LGW_FSK_CH_BW_EXPO                      uint16 = 226
	LGW_FSK_RSSI_LENGTH                     uint16 = 227
	LGW_FSK_RX_INVERT                       uint16 = 228
	LGW_FSK_PKT_MODE                        uint16 = 229
	LGW_FSK_PSIZE                           uint16 = 230
	LGW_FSK_CRC_EN                          uint16 = 231
	LGW_FSK_DCFREE_ENC                      uint16 = 232
	LGW_FSK_CRC_IBM                         uint16 = 233
	LGW_FSK_ERROR_OSR_TOL                   uint16 = 234
	LGW_FSK_RADIO_SELECT                    uint16 = 235
	LGW_FSK_BR_RATIO                        uint16 = 236
	LGW_FSK_REF_PATTERN_LSB                 uint16 = 237
	LGW_FSK_REF_PATTERN_MSB                 uint16 = 238
	LGW_FSK_PKT_LENGTH                      uint16 = 239
	LGW_FSK_TX_GAUSSIAN_EN                  uint16 = 240
	LGW_FSK_TX_GAUSSIAN_SELECT_BT           uint16 = 241
	LGW_FSK_TX_PATTERN_EN                   uint16 = 242
	LGW_FSK_TX_PREAMBLE_SEQ                 uint16 = 243
	LGW_FSK_TX_PSIZE                        uint16 = 244
	LGW_FSK_NODE_ADRS                       uint16 = 245
	LGW_FSK_BROADCAST                       uint16 = 246
	LGW_FSK_AUTO_AFC_ON                     uint16 = 247
	LGW_FSK_PATTERN_TIMEOUT_CFG             uint16 = 248
	LGW_SPI_RADIO_A__DATA                   uint16 = 249
	LGW_SPI_RADIO_A__DATA_READBACK          uint16 = 250
	LGW_SPI_RADIO_A__ADDR                   uint16 = 251
	LGW_SPI_RADIO_A__CS                     uint16 = 252
	LGW_SPI_RADIO_B__DATA                   uint16 = 253
	LGW_SPI_RADIO_B__DATA_READBACK          uint16 = 254
	LGW_SPI_RADIO_B__ADDR                   uint16 = 255
	LGW_SPI_RADIO_B__CS                     uint16 = 256
	LGW_RADIO_A_EN                          uint16 = 257
	LGW_RADIO_B_EN                          uint16 = 258
	LGW_RADIO_RST                           uint16 = 259
	LGW_LNA_A_EN                            uint16 = 260
	LGW_PA_A_EN                             uint16 = 261
	LGW_LNA_B_EN                            uint16 = 262
	LGW_PA_B_EN                             uint16 = 263
	LGW_PA_GAIN                             uint16 = 264
	LGW_LNA_A_CTRL_LUT                      uint16 = 265
	LGW_PA_A_CTRL_LUT                       uint16 = 266
	LGW_LNA_B_CTRL_LUT                      uint16 = 267
	LGW_PA_B_CTRL_LUT                       uint16 = 268
	LGW_CAPTURE_SOURCE                      uint16 = 269
	LGW_CAPTURE_START                       uint16 = 270
	LGW_CAPTURE_FORCE_TRIGGER               uint16 = 271
	LGW_CAPTURE_WRAP                        uint16 = 272
	LGW_CAPTURE_PERIOD                      uint16 = 273
	LGW_MODEM_STATUS                        uint16 = 274
	LGW_VALID_HEADER_COUNTER_0              uint16 = 275
	LGW_VALID_PACKET_COUNTER_0              uint16 = 276
	LGW_VALID_HEADER_COUNTER_MBWSSF         uint16 = 277
	LGW_VALID_HEADER_COUNTER_FSK            uint16 = 278
	LGW_VALID_PACKET_COUNTER_MBWSSF         uint16 = 279
	LGW_VALID_PACKET_COUNTER_FSK            uint16 = 280
	LGW_CHANN_RSSI                          uint16 = 281
	LGW_BB_RSSI                             uint16 = 282
	LGW_DEC_RSSI                            uint16 = 283
	LGW_DBG_MCU_DATA                        uint16 = 284
	LGW_DBG_ARB_MCU_RAM_DATA                uint16 = 285
	LGW_DBG_AGC_MCU_RAM_DATA                uint16 = 286
	LGW_NEXT_PACKET_CNT                     uint16 = 287
	LGW_ADDR_CAPTURE_COUNT                  uint16 = 288
	LGW_TIMESTAMP                           uint16 = 289
	LGW_DBG_CHANN0_GAIN                     uint16 = 290
	LGW_DBG_CHANN1_GAIN                     uint16 = 291
	LGW_DBG_CHANN2_GAIN                     uint16 = 292
	LGW_DBG_CHANN3_GAIN                     uint16 = 293
	LGW_DBG_CHANN4_GAIN                     uint16 = 294
	LGW_DBG_CHANN5_GAIN                     uint16 = 295
	LGW_DBG_CHANN6_GAIN                     uint16 = 296
	LGW_DBG_CHANN7_GAIN                     uint16 = 297
	LGW_DBG_DEC_FILT_GAIN                   uint16 = 298
	LGW_SPI_DATA_FIFO_PTR                   uint16 = 299
	LGW_PACKET_DATA_FIFO_PTR                uint16 = 300
	LGW_DBG_ARB_MCU_RAM_ADDR                uint16 = 301
	LGW_DBG_AGC_MCU_RAM_ADDR                uint16 = 302
	LGW_SPI_MASTER_CHIP_SELECT_POLARITY     uint16 = 303
	LGW_SPI_MASTER_CPOL                     uint16 = 304
	LGW_SPI_MASTER_CPHA                     uint16 = 305
	LGW_SIG_GEN_ANALYSER_MUX_SEL            uint16 = 306
	LGW_SIG_GEN_EN                          uint16 = 307
	LGW_SIG_ANALYSER_EN                     uint16 = 308
	LGW_SIG_ANALYSER_AVG_LEN                uint16 = 309
	LGW_SIG_ANALYSER_PRECISION              uint16 = 310
	LGW_SIG_ANALYSER_VALID_OUT              uint16 = 311
	LGW_SIG_GEN_FREQ                        uint16 = 312
	LGW_SIG_ANALYSER_FREQ                   uint16 = 313
	LGW_SIG_ANALYSER_I_OUT                  uint16 = 314
	LGW_SIG_ANALYSER_Q_OUT                  uint16 = 315
	LGW_GPS_EN                              uint16 = 316
	LGW_GPS_POL                             uint16 = 317
	LGW_SW_TEST_REG1                        uint16 = 318
	LGW_SW_TEST_REG2                        uint16 = 319
	LGW_SW_TEST_REG3                        uint16 = 320
	LGW_DATA_MNGT_STATUS                    uint16 = 321
	LGW_DATA_MNGT_CPT_FRAME_ALLOCATED       uint16 = 322
	LGW_DATA_MNGT_CPT_FRAME_FINISHED        uint16 = 323
	LGW_DATA_MNGT_CPT_FRAME_READEN          uint16 = 324
	LGW_TX_TRIG_ALL                         uint16 = 325
	LGW_TOTALREGS                           uint16 = 326
)

func Page_switch(c *os.File, lgw_spi_mux_mode, target byte) error {
	lgw_regpage := PAGE_MASK & target
	err := Lgw_spi_w(c, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, PAGE_ADDR, lgw_regpage)
	if err != nil {
		return err
	}
	return nil
}

func Check_fpga_version(version byte) bool {
	for i := 0; i < len(FPGA_VERSION); i++ {
		if FPGA_VERSION[i] == version {
			return true
		}
	}
	return false
}

func reg_w_align32(c *os.File, spi_mux_mode, spi_mux_target byte, r Lgw_reg_s, reg_value int32) error {

	buf := make([]byte, 4)
	if (r.leng == 8) && (r.offs == 0) {
		/* direct write */
		err := Lgw_spi_w(c, spi_mux_mode, spi_mux_target, r.addr, byte(reg_value))
		if err != nil {
			return err
		}
	} else if (r.offs + r.leng) <= 8 {
		/* single-byte read-modify-write, offs:[0-7], leng:[1-7] */
		b, err := Lgw_spi_r(c, spi_mux_mode, spi_mux_target, r.addr)
		if err != nil {
			return err
		}
		buf[0] = b
		buf[1] = ((1 << r.leng) - 1) << r.offs            /* bit mask */
		buf[2] = (byte(reg_value)) << r.offs              /* new data offsetted */
		buf[3] = (^(buf[1]) & buf[0]) | (buf[1] & buf[2]) /* mixing old & new data */
		err = Lgw_spi_w(c, spi_mux_mode, spi_mux_target, r.addr, buf[3])
		if err != nil {
			return err
		}
	} else if (r.offs == 0) && (r.leng > 0) && (r.leng <= 32) {
		/* multi-byte direct write routine */
		size_byte := (r.leng + 7) / 8 /* add a byte if it's not an exact multiple of 8 */

		for i := 0; i < int(size_byte); i++ {
			/* big endian register file for a file on N bytes
			   Least significant byte is stored in buf[0], most one in buf[N-1] */
			buf[i] = byte(0x000000FF & reg_value)
			reg_value = (reg_value >> 8)
		}
		err := Lgw_spi_wb(c, spi_mux_mode, spi_mux_target, r.addr, buf[0:size_byte]) /* write the register in one burst */
		if err != nil {
			return err
		}
	} else {
		/* register spanning multiple memory bytes but with an offset */
		return fmt.Errorf("ERROR: REGISTER SIZE AND OFFSET ARE NOT SUPPORTED\n")
	}
	return nil
}

func reg_r_align32(c *os.File, spi_mux_mode, spi_mux_target byte, r Lgw_reg_s) (int32, error) {
	bufu := make([]byte, 4)

	if (r.offs + r.leng) <= 8 {
		/* read one byte, then shift and mask bits to get reg value with sign extension if needed */

		b, err := Lgw_spi_r(c, spi_mux_mode, spi_mux_target, r.addr)
		if err != nil {
			return 0, err
		}
		bufu[0] = b
		bufu[1] = bufu[0] << (8 - r.leng - r.offs) /* left-align the data */
		if r.sign == 1 {
			bs := int8(bufu[1]) >> (8 - r.leng) /* right align the data with sign extension (ARITHMETIC right shift) */
			return int32(bs), nil               /* signed pointer -> 32b sign extension */
		} else {
			bufu[2] = bufu[1] >> (8 - r.leng) /* right align the data, no sign extension */
			return int32(bufu[2]), nil        /* unsigned pointer -> no sign extension */
		}
	} else if (r.offs == 0) && (r.leng > 0) && (r.leng <= 32) {
		size_byte := int((r.leng)+7) / 8 /* add a byte if it's not an exact multiple of 8 */
		bufu, err := Lgw_spi_rb(c, spi_mux_mode, spi_mux_target, r.addr, uint16(size_byte))
		if err != nil {
			return 0, err
		}
		var u uint32
		for i := (size_byte - 1); i >= 0; i-- {
			u = uint32(bufu[i]) + (u << 8) /* transform a 4-byte array into a 32 bit word */
		}
		if r.sign == 1 {
			u = u << (32 - r.leng)                /* left-align the data */
			return int32(u) >> (32 - r.leng), nil /* right-align the data with sign extension (ARITHMETIC right shift) */
		} else {
			return int32(u), nil /* unsigned value -> return 'as is' */
		}
	} else {
		/* register spanning multiple memory bytes but with an offset */
		return 0, fmt.Errorf("ERROR: REGISTER SIZE AND OFFSET ARE NOT SUPPORTED\n")
	}
}

func Lgw_connect(path string, spi_only bool, tx_notch_freq uint32) (*os.File, byte, byte, error) {

	var lgw_spi_mux_mode, lgw_spi_mux_target byte
	/* open the SPI link */
	c, err := Lgw_spi_open(path)
	if err != nil {
		return nil, lgw_spi_mux_mode, lgw_spi_mux_target, err
	}

	if spi_only == false {
		/* Detect if the gateway has an FPGA with SPI mux header support */
		/* First, we assume there is an FPGA, and try to read its version */
		u, err := Lgw_spi_r(c, LGW_SPI_MUX_MODE1, LGW_SPI_MUX_TARGET_FPGA, loregs[LGW_VERSION].addr)
		if err != nil {
			return nil, lgw_spi_mux_mode, lgw_spi_mux_target, err
		}

		if Check_fpga_version(u) != true {
			/* We failed to read expected FPGA version, so let's assume there is no FPGA */
			log.Printf("INFO: no FPGA detected or version not supported (v%d)\n", u)
			lgw_spi_mux_mode = LGW_SPI_MUX_MODE0
			lgw_spi_mux_target = LGW_SPI_MUX_TARGET_SX1301
		} else {
			fmt.Printf("INFO: detected FPGA with SPI mux header (v%d)\n", u)
			lgw_spi_mux_mode = LGW_SPI_MUX_MODE1
			lgw_spi_mux_target = LGW_SPI_MUX_TARGET_FPGA
			/* FPGA Soft Reset */
			Lgw_spi_w(c, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_FPGA, 0, 1)
			Lgw_spi_w(c, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_FPGA, 0, 0)
			/* FPGA configure */
			err := Lgw_fpga_configure(c, tx_notch_freq)
			if err != nil {
				return nil, lgw_spi_mux_mode, lgw_spi_mux_target, err
			}
			return c, lgw_spi_mux_mode, lgw_spi_mux_target, nil
		}

		/* check SX1301 version */
		u, err = Lgw_spi_r(c, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, loregs[LGW_VERSION].addr)
		if err != nil {
			return nil, lgw_spi_mux_mode, lgw_spi_mux_target, fmt.Errorf("ERROR READING CHIP VERSION REGISTER\n")
		}
		if u != byte(loregs[LGW_VERSION].dflt) {
			return nil, lgw_spi_mux_mode, lgw_spi_mux_target, fmt.Errorf("ERROR: NOT EXPECTED CHIP VERSION (v%d)\n", u)
		}

		/* write 0 to the page/reset register */
		err = Lgw_spi_w(c, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, loregs[LGW_PAGE_REG].addr, 0)
		if err != nil {
			return nil, lgw_spi_mux_mode, lgw_spi_mux_target, fmt.Errorf("ERROR WRITING PAGE REGISTER\n")
		}
	}

	fmt.Printf("Note: success connecting the concentrator\n")
	return c, lgw_spi_mux_mode, lgw_spi_mux_target, nil
}

func Lgw_disconnect(f *os.File) error {
	err := Lgw_spi_close(f)
	if err != nil {
		return err
	}
	return nil
}

func Lgw_soft_reset(f *os.File, lgw_spi_mux_mode byte) error {
	/* check if SPI is initialised */
	err := Lgw_spi_w(f, lgw_spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, 0, 0x80) /* 1 -> SOFT_RESET bit */
	if err != nil {
		return err
	}
	return nil
}

/* register verification */
func Lgw_reg_check(c *os.File, spi_mux_mode, spi_mux_target byte) error {
	r := Lgw_reg_s{}
	ok_msg := "+++MATCH+++"
	notok_msg := "###MISMATCH###"

	fmt.Print("Start of register verification\n")
	for i := uint16(0); i < LGW_TOTALREGS; i++ {
		r = loregs[i]
		read_value, err := Lgw_reg_r(c, spi_mux_mode, spi_mux_target, uint16(i))
		if err != nil {
			return err
		}
		ptr := ""
		if read_value == r.dflt {
			ptr = ok_msg
		} else {
			ptr = notok_msg
		}
		if r.sign == 1 {
			fmt.Printf("%s reg number %d read: %d (%x) default: %d (%X)\n", ptr, i, read_value, read_value, r.dflt, r.dflt)
		} else {
			fmt.Printf("%s reg number %d read: %d (%x) default: %d (%X)\n", ptr, i, read_value, read_value, r.dflt, r.dflt)
		}
	}
	fmt.Print("End of register verification\n")

	return nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Write to a register addressed by name */
func Lgw_reg_w(c *os.File, spi_mux_mode, spi_mux_target byte, register_id uint16, reg_value int32) error {
	r := Lgw_reg_s{}

	/* check input parameters */
	if register_id >= LGW_TOTALREGS {
		return fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* intercept direct access to PAGE_REG & SOFT_RESET */
	if register_id == LGW_PAGE_REG {
		Page_switch(c, spi_mux_mode, byte(reg_value))
		return nil
	} else if register_id == LGW_SOFT_RESET {
		/* only reset if lsb is 1 */
		if reg_value&0x01 != 0 {
			err := Lgw_soft_reset(c, spi_mux_mode)
			if err != nil {
				return err
			}
		}
		return nil
	}

	/* get register struct from the struct array */
	r = loregs[register_id]

	/* reject write to read-only registers */
	if r.rdon == 1 {
		fmt.Errorf("ERROR: TRYING TO WRITE A READ-ONLY REGISTER\n")
	}

	/* select proper register page if needed */
	if r.page != -1 {
		Page_switch(c, spi_mux_mode, byte(r.page))
	}

	err := reg_w_align32(c, spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, r, reg_value)
	if err != nil {
		return err
	}
	return nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Read to a register addressed by name */
func Lgw_reg_r(c *os.File, spi_mux_mode, spi_mux_target byte, register_id uint16) (int32, error) {
	r := Lgw_reg_s{}

	/* check input parameters */
	if register_id >= LGW_TOTALREGS {
		return 0, fmt.Errorf("ERROR: REGISTER NUMBER OUT OF DEFINED RANGE\n")
	}

	/* get register struct from the struct array */
	r = loregs[register_id]

	/* select proper register page if needed */
	if r.page != -1 {
		Page_switch(c, spi_mux_mode, byte(r.page))
	}

	val, err := reg_r_align32(c, spi_mux_mode, spi_mux_target, r)
	if err != nil {
		return 0, err
	}
	return val, nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Point to a register by name and do a burst write */
func Lgw_reg_wb(fd *os.File, spi_mux_mode, spi_mux_target byte, register_id uint16, data []byte) error {
	/* get register struct from the struct array */
	r := loregs[register_id]

	/* reject write to read-only registers */
	if r.rdon == 1 {
		fmt.Errorf("ERROR: TRYING TO BURST WRITE A READ-ONLY REGISTER\n")
	}

	/* select proper register page if needed */
	if r.page != -1 {
		Page_switch(fd, spi_mux_mode, byte(r.page))
	}

	/* do the burst write */
	err := Lgw_spi_wb(fd, spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, r.addr, data)
	if err != nil {
		return err
	}
	return nil
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

/* Point to a register by name and do a burst read */
func Lgw_reg_rb(c *os.File, spi_mux_mode, spi_mux_target byte, register_id uint16, size uint16) ([]byte, error) {
	/* get register struct from the struct array */
	r := loregs[register_id]

	/* select proper register page if needed */
	if r.page != -1 {
		Page_switch(c, spi_mux_mode, byte(r.page))
	}

	/* do the burst read */
	val, err := Lgw_spi_rb(c, spi_mux_mode, LGW_SPI_MUX_TARGET_SX1301, r.addr, size)
	if err != nil {
		return nil, err
	}
	return val, nil
}
