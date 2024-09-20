/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Size Reporting Mode enumeration
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"errors"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	SizeModeStd SizeMode = iota // full numeric size
	SizeModeSI                  // Sistema Internacional: 1K = 1000
	SizeModeIEC                 // Binary System: 1KB = 1024
)

var (
	ErrUnknownSizeMode = errors.New("Unknown size mode (Std|SI|IEC")
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type SizeMode uint

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Stringer interface
func (s SizeMode) String() string {
	switch s {
	case SizeModeStd:
		return "Standard"
	case SizeModeSI:
		return "International"
	case SizeModeIEC:
		return "Binary"
	default:
		panic(ErrUnknownSizeMode)
	}
}

func (s SizeMode) ShortString() string {
	switch s {
	case SizeModeStd:
		return "Std"
	case SizeModeSI:
		return "SI"
	case SizeModeIEC:
		return "IEC"
	default:
		panic(ErrUnknownSizeMode)
	}
}
