package liblorago

func TAKE_N_BITS_FROM(b, p, n byte) byte {
	return (((b) >> (p)) & ((1 << (n)) - 1))
}
