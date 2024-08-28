/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Given a total size return a string with that value formatted in
// either of the size formats (Standard, International, Binary)
func ReportByteCount(count int64, mode SizeMode) string {
	var output string
	switch mode {
	case SizeModeSI:
		output = ByteCountSI(count)
		break
	case SizeModeIEC:
		output = ByteCountIEC(count)
		break
	case SizeModeStd:
		fallthrough
	default:
		output = AddThousands(count, ',')
	}
	return output
}

// Byte count formatted using International System (1K = 1000)
func ByteCountSI(b int64) string {
	const UNIT = 1000
	if b < UNIT {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(UNIT), 0
	for n := b / UNIT; n >= UNIT; n /= UNIT {
		div *= UNIT
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// Byte count formatted using IEC (Binary) system (1K = 1024)
func ByteCountIEC(b int64) string {
	const UNIT = 1024
	if b < UNIT {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(UNIT), 0
	for n := b / UNIT; n >= UNIT; n /= UNIT {
		div *= UNIT
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func AddThousands(nr int64, sep rune) string {
	var result string = ""
	nrS := Reverse(fmt.Sprintf("%d", nr))

	for start := 0; start < len(nrS); start += 3 {
		if start+3 < len(nrS) {
			group := nrS[start : start+3]
			if start != 0 {
				result = result + string(sep) + group
			} else {
				result = group
			}
		} else {
			group := nrS[start:]
			result = result + string(sep) + group
			break
		}
	}
	return Reverse(result)
}
