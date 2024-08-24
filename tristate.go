/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium


/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	Undecided TriState = iota
	No
	Yes
)

var (
	TriStateDefaultLabels = []string{"Undecided", "No", "Yes"}
	TriStateDingbatLabels = []string{"\u2753", "\u2718", "\u2714"} // ❓ ✘ ✔
	TriStateSquarishLabels  = []string{"\u2b1a", "\u2b1c", "\u2b1b"} // ⬚ ⬜ ⬛
	TriStateSquareLabels  = []string{"\u2610", "\u2612", "\u2611"} // ☐ ☒ ☑
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type TriState uint

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (t TriState) String() string {
	return t.StringWith(TriStateDefaultLabels)
}

func (t TriState) StringWith(v []string) string {
	if len(v) != 3 && len(v) != 4 {
		panic("TriState.StringWith() must have 3 or 4 strings")
	}

	if t == Undecided {
		return v[0]
	} else if t == No {
		return v[1]
	} else if t == Yes {
		return v[2]
	} else {
		if len(v) == 3 {
			return "!!!Invalid"
		} else {
			return v[3]	// user-defined invalid label
		}
	}
}