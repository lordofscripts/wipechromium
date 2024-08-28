/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	"testing"

	cmn "github.com/lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

func Test_DefaultValue(t *testing.T) {
	var sm cmn.SizeMode

	if sm != cmn.SizeModeStd {
		t.Errorf("Default value expected Std but got %s", sm.ShortString())
	}
}

func Test_ValidValues(t *testing.T) {
	sm := cmn.SizeModeStd
	if sm.ShortString() != "Std" {
		t.Errorf("Expected Std but got %s", sm.ShortString())
	}

	sm = cmn.SizeModeSI
	if sm.ShortString() != "SI" {
		t.Errorf("Expected SI but got %s", sm.ShortString())
	}

	sm = cmn.SizeModeIEC
	if sm.ShortString() != "IEC" {
		t.Errorf("Expected Std but got %s", sm.ShortString())
	}
}

func Test_InvalidValue(t *testing.T) {
	var sm cmn.SizeMode = 5 // invalid

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SizeMode.String() did not panic")
		}
	}()

	// String() should panic
	if sm.String() == "" {
		t.Errorf("This shouldn't happen")
	}

	// ShortString() should panic
	if sm.ShortString() == "" {
		t.Errorf("This shouldn't happen")
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/
