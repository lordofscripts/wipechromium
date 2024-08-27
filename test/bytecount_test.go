/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	"fmt"
	cmn "lordofscripts/wipechromium"
	"math"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var subCases []subCase = []subCase{
	{999, "999 B", "999 B"},
	{1000, "1.0 kB", "1000 B"},
	{1023, "1.0 kB", "1023 B"},
	{1024, "1.0 kB", "1.0 KiB"},
	{987654321, "987.7 MB", "941.9 MiB"},
	{math.MaxInt64, "9.2 EB", "8.0 EiB"},
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/
type subCase struct {
	Input     int64
	OutputSI  string
	OutputIEC string
}

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/
func Test_ByteCountSI(t *testing.T) {
	for i, sub := range subCases {
		result := cmn.ByteCountSI(sub.Input)
		if result != sub.OutputSI {
			t.Errorf("#%d SI %d expected %s but got %s", i, sub.Input, sub.OutputSI, result)
		} else {
			fmt.Println(sub.Input, sub.OutputSI)
		}
	}
}

func Test_ByteCountIEC(t *testing.T) {
	for i, sub := range subCases {
		result := cmn.ByteCountIEC(sub.Input)
		if result != sub.OutputIEC {
			t.Errorf("#%d SI %d expected %s but got %s", i, sub.Input, sub.OutputIEC, result)
		} else {
			fmt.Println(sub.Input, sub.OutputIEC)
		}
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/
