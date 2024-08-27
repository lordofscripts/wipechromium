/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	cmn "github.com/lordofscripts/wipechromium"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	triStateDefault  cmn.TriState
	triStatePositive cmn.TriState = cmn.Yes
	triStateNegative cmn.TriState = cmn.No
	triStateUnset    cmn.TriState = cmn.Undecided
	triStateInvalid  cmn.TriState = cmn.TriState(5)
)

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

// Check values with default labels
func Test_Values(t *testing.T) {
	const (
		FLOATING        = "Undecided"
		HIGH            = "Yes"
		LOW             = "No"
		INVALID_DEFAULT = "!!!Invalid"
	)

	if triStateDefault.String() != FLOATING {
		t.Errorf("Default value is not %q != %q", triStateDefault.String(), FLOATING)
	}

	if triStateUnset.String() != FLOATING {
		t.Errorf("Unset value is not %q != %q", triStateDefault.String(), FLOATING)
	}

	if triStateNegative.String() != LOW {
		t.Errorf("Negated value is not %q", LOW)
	}

	if triStatePositive.String() != HIGH {
		t.Errorf("Asserted value is not %q", HIGH)
	}
}

// Check values with custom labels as in electronic circuits
func Test_ValuesCustomLabel(t *testing.T) {
	const (
		FLOATING        = "Z"
		HIGH            = "H"
		LOW             = "L"
		INVALID_DEFAULT = "!!!Invalid"
		INVALID_CUSTOM  = "???"
	)

	// A. Using the default String()
	labels := []string{FLOATING, LOW, HIGH}

	if triStateDefault.StringWith(labels) != FLOATING {
		t.Errorf("Default value is not %q != %q", triStateDefault.String(), FLOATING)
	}

	if triStateUnset.StringWith(labels) != FLOATING {
		t.Errorf("Unset value is not %q != %q", triStateUnset.String(), FLOATING)
	}

	if triStateNegative.StringWith(labels) != LOW {
		t.Errorf("Negated value is not %q", LOW)
	}

	if triStatePositive.StringWith(labels) != HIGH {
		t.Errorf("Asserted value is not %q != %q", triStatePositive.String(), HIGH)
	}

	if triStateInvalid.StringWith(labels) != INVALID_DEFAULT {
		t.Errorf("Out-of-range value is not %q != %q", triStateInvalid.String(), INVALID_DEFAULT)
	}

	// B. Using the custom StringWith()
	// Now add a custom out-of-range value
	labels = append(labels, INVALID_CUSTOM)
	triS := triStateInvalid.StringWith(labels)

	if triS != INVALID_CUSTOM {
		t.Errorf("Out-of-range value is not %q", INVALID_CUSTOM)
	}
}

// Check StringWith() with valid parameter
func Test_StringWith(t *testing.T) {
	const (
		FLOATING        = "Z"
		HIGH            = "H"
		LOW             = "L"
		INVALID_DEFAULT = "!!!Invalid"
		INVALID_CUSTOM  = "???"
	)

	labelsGood1 := []string{FLOATING, LOW, HIGH}
	labelsGood2 := []string{FLOATING, LOW, HIGH, INVALID_CUSTOM}

	tri := cmn.Yes
	if tri.StringWith(labelsGood1) != HIGH {
		t.Errorf("Bad value")
	}

	if tri.StringWith(labelsGood2) != HIGH {
		t.Errorf("Bad value")
	}
}

// Check StringWith() with invalid parameter
func Test_StringWithPanic1(t *testing.T) {
	const (
		FLOATING        = "Z"
		HIGH            = "H"
		LOW             = "L"
		INVALID_DEFAULT = "!!!Invalid"
		INVALID_CUSTOM  = "???"
	)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TriState.StringWith() did not panic")
		}
	}()

	labelsBad1 := []string{FLOATING, LOW, INVALID_CUSTOM, FLOATING, FLOATING}

	tri := cmn.Yes
	if tri.StringWith(labelsBad1) == "" {
		t.Errorf("This shouldn't happen")
	}
}

// Check StringWith() with invalid parameter
func Test_StringWithPanic2(t *testing.T) {
	const (
		FLOATING        = "Z"
		HIGH            = "H"
		LOW             = "L"
		INVALID_DEFAULT = "!!!Invalid"
		INVALID_CUSTOM  = "???"
	)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TriState.StringWith() did not panic")
		}
	}()

	labelsBad2 := []string{FLOATING, LOW}

	tri := cmn.Yes
	if tri.StringWith(labelsBad2) == "" {
		t.Errorf("This shouldn't happen")
	}
}
