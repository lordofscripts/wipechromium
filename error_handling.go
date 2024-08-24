/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Custom error.
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
	"errors"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
var (
	ErrUnsupportedBrowser = errors.New("Unsupported browser")
	ErrNotBrowserCache = errors.New("Not a browser cache directory")
	ErrNotBrowserProfile = errors.New("Not a browser user profile directory")
	ErrNoProfile = errors.New("Browser user profile not given")
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/


/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// wrap two errors into one (as of GO v1.13)
// Example: wrapError(err, "Unable to remove file %s", filename)
func WrapError(err error, code int,  msgformat string, v ...any) error {
	errP := fmt.Errorf(msgformat, v...)
	errC := fmt.Errorf("Error E-%03d: %w...\n\tInner: %w", code, errP, err)
	return errC
}