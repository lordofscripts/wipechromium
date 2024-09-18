/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Custom error & error formatting.
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
var (
	ErrUnsupportedBrowser  = errors.New("Unsupported browser")
	ErrNotBrowserCache     = errors.New("Not a browser cache directory")
	ErrNotBrowserProfile   = errors.New("Not a browser user profile directory")
	ErrNoProfile           = errors.New("Browser user profile not given")
	ErrCleanerFailure      = errors.New("Browser cleaner instantiation error")
	ErrProfileDoesNotExist = errors.New("Profile does not exist!")
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
func WrapError(err error, code int, msgformat string, v ...any) error {
	errP := fmt.Errorf(msgformat, v...)
	errC := fmt.Errorf("Error E-%03d: %w...\n\tInner: %w", code, errP, err)
	return errC
}

// ThisLocation() returns a string with source code location information that
// is usable for error reports. Things like package name, filename, line nr., etc.
// The frame parameter should be 1 if the error ocurred in the calling function.
func ThisLocation(frame int) string {
	const (
		Marker     rune = '⛿' // else ⛔
		FrameLevel int  = 1   // 1 for in-file demo, 2 from elsewhere
	)
	type CallerInfo struct {
		packageN  string
		structure string
		function  string
		filename  string
		lineno    int
	}

	// (CallerInfo) Formatted stringify takes any of the following format specifiers:
	// %P package name
	// %S structure name (or empty if none), i.e. Event{}
	// %F or %M function/method name (aliased), i.e. Sum()
	// %L line nr.
	// %A short version (%P)
	// %B median version (%P.%F)
	// %C long version (%P.%S.%M)
	ciStringer := func(format string, c *CallerInfo) string {
		const (
			IFUNC = "()"
			ISTRU = "{}"
		)
		//out := strings.ToUpper(format)
		out := format
		// macro replacements
		out = strings.Replace(out, "%A", "%P", 1)
		out = strings.Replace(out, "%B", "%P.%F#%L", 1)
		out = strings.Replace(out, "%C", "%P.%S.%M#%L", 1)
		out = strings.Replace(out, "%c", "%p.%S.%M#%L", 1)
		var struN string = ""
		if len(c.structure) > 0 {
			struN = c.structure + ISTRU
		}
		// atomic replacements
		out = strings.Replace(out, "%P", c.packageN, 1)
		out = strings.Replace(out, "%S", struN, 1)
		out = strings.Replace(out, "%F", c.function+IFUNC, 1)
		out = strings.Replace(out, "%M", c.function+IFUNC, 1)
		out = strings.Replace(out, "%L", strconv.Itoa(c.lineno), 1)
		if strings.Index(out, "%p") > -1 && strings.LastIndexByte(c.packageN, '/') > 0 {
			pname := c.packageN[strings.LastIndexByte(c.packageN, '/')+1:]
			out = strings.Replace(out, "%p", pname, 1)
		} else {
			pname := c.packageN[strings.LastIndexByte(c.packageN, '/')+1:]
			out = strings.Replace(out, "%p", pname, 1)
		}
		return out
	}

	if frame < FrameLevel {
		frame = FrameLevel
	}
	pc, fileName, lineNo, ok := runtime.Caller(frame)

	var location string = "???"
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		//fmt.Println(funcName)
		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

		//var title string
		var pkg, stru, fun string
		//fmt.Printf("RCI ***%s\n", funcName)
		if idx := strings.IndexByte(funcName[:lastDot], '.'); idx > -1 {
			// Value struct:  Event
			// Pointer to struct: (*Event)
			// Package Init(): init
			stru = funcName[idx+1 : lastDot]
			pkg = funcName[:idx]
			if stru == "init" {
				// A package.init() comes as stru:init func:0
				stru = ""
			}
			//title = "Method "
		} else {
			stru = ""
			pkg = funcName[:lastDot]
			//title = "Func   "
		}

		fun = strings.Trim(funcName[lastDot+1:], " ")
		if fun == "0" {
			fun = "init"
		}
		ci := &CallerInfo{packageN: pkg, structure: stru, function: fun, filename: fileName, lineno: lineNo}
		if len(ci.structure) == 0 {
			location = ciStringer("%B", ci)
		} else {
			location = ciStringer("%c", ci) // or %C
		}
	}

	return fmt.Sprintf("%c %s", Marker, location)
}

// SpitOutError() nicely formats an application error in a uniform way.
// The frame parameter should be 1 if the error ocurred in the calling function.
func SpitOutError(frame int, err error) {
	const (
		Label = "⚓ Application Error: "
	)

	location := ThisLocation(frame + 1)
	fmt.Printf("%s\n⚓\t%s :: %s\n", Label, location, err)
}
