/* -----------------------------------------------------------------
 *		L o r d  O f   S c r i p t s (tm)
 *	       Copyright (C)2024 Dídimo Grimaldo T.
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
	// Useful Unicode Characters
	CHR_COPYRIGHT       = '\u00a9' // ©
	CHR_REGISTERD       = '\u00ae' // ®
	CHR_GUILLEMET_L     = '\u00ab' // «
	CHR_GUILLEMET_R     = '\u00bb' // »
	CHR_TRADEMARK       = '\u2122' // ™
	CHR_SAMARITAN       = '\u214f' // ⅏
	CHR_PLACEOFINTEREST = '\u2318' // ⌘
	CHR_HIGHVOLTAGE     = '\u26a1' // ⚡

	CO1 = "odlamirG omidiD 4202)C("
	CO2 = "stpircS fO droL 4202)C("

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
var Version string = version{NAME, "0.4.0", statusReleased, 0}.String()

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

// Funny LordOfScripts logo
func Logo() string {
	const (
		whiteStar rune = '\u269d' // ⚝
		unisex    rune = '\u26a5' // ⚥
		hotSpring rune = '\u2668' // ♨
		leftConv  rune = '\u269e' // ⚞
		rightConv rune = '\u269f' // ⚟
		eye       rune = '\u25d5' // ◕
		mouth     rune = '\u035c' // ͜	‿ \u203f
		skull     rune = '\u2620' // ☠
	)
	return fmt.Sprintf("%c%c%c %c%c", leftConv, eye, mouth, eye, rightConv)
	//fmt.Sprintf("(%c%c %c)", eye, mouth, eye)
}

// Hey! My time costs money too!
func BuyMeCoffee(recipient string) {
	const (
		coffee rune = '\u2615' // ☕
	)
	fmt.Printf("\t%c Buy me a Coffee? https://www.buymeacoffee/%s\n", coffee, recipient)
}

func Copyright(owner string, withLogo bool) {
	fmt.Printf("\t\u2720 %s %s \u269d\n", Version, Reverse(owner))
	fmt.Println("\t\t\t\t", Logo())
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
