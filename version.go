/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *	PACKAGE VERSION
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	// Change these values accordingly
	NAME string = "WipeChromium"
	DESC string = "Wipes out all (or selected) Chromium junk"
	// don't change
	statusAlpha    status = "Alpha"
	statusBeta     status = "Beta"
	statusRC       status = "RC" // Release Candidate
	statusReleased status = ""
)

// NOTE: Change these values accordingly
var Version string = version{NAME, "0.1", statusRC, 1}.String()
var Copyright string = "Copyright (C)2024 Didimo Grimaldo"

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/
type status = string

type version struct {
	n  string
	v  string
	s  status
	sv int
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/
func (v version) String() string {
	var ver string
	switch v.s {
	case statusAlpha:
		fallthrough
	case statusBeta:
		fallthrough
	case statusRC:
		ver = fmt.Sprintf("%s v%s-%s-%d", v.n, v.v, v.s, v.sv)
		break
	default:
		ver = fmt.Sprintf("%s v%s", v.n, v.v)
		break
	}
	return ver
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

