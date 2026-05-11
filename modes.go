package rcd

import "os"

const execBits os.FileMode = 0o111

func maskedMode(mode os.FileMode) os.FileMode {
	return mode &^ execBits
}

func unmaskedMode(mode os.FileMode) os.FileMode {
	return mode | (mode&0o444)>>2
}
