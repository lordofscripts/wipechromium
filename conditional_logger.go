/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"log"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var _ ILogger = (*ConditionalLogger)(nil)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ILogger interface {
	Printf(template string, v ...any)
	Print(v ...any)

	IsEnabled() bool
	InheritAs(string) *ConditionalLogger
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type ConditionalLogger struct {
	enabled bool
	prefix	string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewConditionalLogger(enabled bool, prefix string) *ConditionalLogger {
	return &ConditionalLogger{enabled, prefix}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Whether log is enabled or not
func (l *ConditionalLogger) IsEnabled() bool {
	return l.enabled
}

// A new logger inheriting configuration from the parent but different prefix.
func (l *ConditionalLogger) InheritAs(prefix string) *ConditionalLogger {
	return NewConditionalLogger(l.IsEnabled(), prefix)
}

// Formatted logging
func (l *ConditionalLogger) Printf(template string, v ...any) {
	if l.enabled {
		log.Printf(l.prefix + " " +template, v...)
	}
}

// Print items on log stream
func (l *ConditionalLogger) Print(v ...any) {
	if l.enabled {
		v = append([]any{l.prefix + " "}, v...)
		log.Print(v...)
	}
}




