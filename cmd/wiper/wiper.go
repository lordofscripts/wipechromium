/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Utility to Clear Chromium Cache & Profile Data for PRIVACY.
 * Settings, Bookmarks, FileSystem & Extensions are cleaned but not deleted.
 * Supports Chromium & Firefox ESR browsers.
 * Created: 21 Aug 2024
 * Updated: 18 Sep 2024
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	cmn "github.com/lordofscripts/wipechromium"
	// Here one package for each supported browser
	"github.com/lordofscripts/wipechromium/browsers"
	"github.com/lordofscripts/wipechromium/browsers/chromium"
	"github.com/lordofscripts/wipechromium/browsers/firefox"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	FLAG_HELP_BROWSER string = "Browser name"
	FLAG_HELP_SCAN    string = "Scan for browsers"
	FLAG_HELP_NAME    string = "Profile name"
	FLAG_HELP_ME      string = "This help"
	FLAG_HELP_CACHE   string = "Erase cache only"
	FLAG_HELP_PROFILE string = "Erase profile junk only"
	FLAG_HELP_SIZE    string = "Select size reporting mode (Std, SI, IEC)"
	FLAG_HELP_LOG     string = "Enable log output"
	FLAG_HELP_DRYRUN  string = "Enable dry-run"
)

var (
	// A superbly simple conditional logger
	logx cmn.ILogger = nil
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type BrowserWipe struct {
	cleaner  browsers.IBrowsers
	SizeMode cmn.SizeMode
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Scan the system for all supported browsers and indicate whether the
// data/cache directories exist
func (b *BrowserWipe) Scan() {
	for _, br := range browsers.SupportedBrowsers {
		err := b.GetCleaner(br, "", true, b.SizeMode, false)
		if err != nil {
			fmt.Printf("\tBad Thing: %s %s\n", br, err)
		} else {
			b.cleaner.Tell()
			installed := b.cleaner.IdentifyAppDataRoot()
			fmt.Printf("\tInstalled: %t\n", installed)
			if profiles, err := b.cleaner.FindProfileNames(); err == nil {
				fmt.Println("\tProfiles :")
				for _, p := range profiles {
					fmt.Println("\t\t- ", p)
				}
			}
		}
	}
}

// Browser Cleaner factory method
func (b *BrowserWipe) GetCleaner(which browsers.Browser, profile string, scanning bool, mode cmn.SizeMode, dryRun bool) error {
	if which == browsers.ChromiumBrowser {
		b.cleaner = chromium.NewChromiumCleaner(profile, mode, dryRun, logx)
	} else if which == browsers.FirefoxBrowser {
		b.cleaner = firefox.NewFirefoxCleaner(profile, scanning, mode, dryRun, logx)
	} else {
		return cmn.ErrUnsupportedBrowser
	}

	if b.cleaner == nil {
		return cmn.ErrCleanerFailure
	}
	return nil
}

func (b *BrowserWipe) Run(cacheOnly, profileOnly bool) (int, error) {
	if err, code := b.cleaner.ClearProfile(cacheOnly, profileOnly); err != nil {
		return code, err
	}

	fmt.Println(b.cleaner.String())
	return 0, nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func help() {
	const (
		RECIPIENT = "lostinwriting"
		PACKAGE   = "WipeChromium"
	)

	cmn.Copyright(cmn.CO1, true)
	fmt.Println("Usage:")
	fmt.Println("\tScan the system for browser data presence.")
	fmt.Println("\t\twipechromium -scan")

	fmt.Println("\tErase both profile junk & cache")
	fmt.Println("\t\twipechromium -b Chromium  -n 'Profile 1'")

	fmt.Println("\tErase profile junk only")
	fmt.Println("\t\twipechromium -b Chromium -n 'Profile 1' -p")

	fmt.Println("\tErase profile cache only")
	fmt.Println("\t\twipechromium -b Chromium -n 'Profile 1' -c")

	fmt.Println("Options:")
	const HELP_TEMPLATE string = "\t%2s %-8s %10s %s\n"
	fmt.Printf(HELP_TEMPLATE, "Op", "Long", "Parameter", "Description")
	fmt.Printf(HELP_TEMPLATE, "-n", "-name", "PROFILE", FLAG_HELP_NAME)
	fmt.Printf(HELP_TEMPLATE, "-c", "-cache", "", FLAG_HELP_CACHE)
	fmt.Printf(HELP_TEMPLATE, "-p", "-profile", "", FLAG_HELP_PROFILE)
	fmt.Printf(HELP_TEMPLATE, "-b", "-browser", "BROWSER", FLAG_HELP_BROWSER)
	fmt.Printf(HELP_TEMPLATE, "-z", "-size", "Std", FLAG_HELP_SIZE)
	fmt.Printf(HELP_TEMPLATE, "-s", "-scan", "", FLAG_HELP_SCAN)
	//fmt.Printf(HELP_TEMPLATE, "", "-log", "", FLAG_HELP_LOG)
	fmt.Printf(HELP_TEMPLATE, "", "-dry", "", FLAG_HELP_DRYRUN) // hidden option

	cmn.BuyMeCoffee(RECIPIENT)
}

// Show a message and die with exit code
func die(exitCode int, msgformat string, v ...any) {
	const CR = "\n"
	if !strings.HasSuffix(msgformat, CR) {
		msgformat += CR
	}

	fmt.Printf("⚓ Bad Thing Happened: exit code %d\n", exitCode)
	fmt.Println("⚓ Message:")
	fmt.Printf("⚓ \t"+msgformat, v...)

	os.Exit(exitCode)
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

// Usage: wipechromium -p 'Profile 1'
func main() {
	// A. Command-line options
	var profile, browserName, szmodeS string
	var cacheOnly, profileOnly, logging, scanOnly, dryRun, helpme bool
	flag.StringVar(&browserName, "b", browsers.ChromiumBrowser.String(), FLAG_HELP_BROWSER)
	flag.StringVar(&browserName, "browser", browsers.ChromiumBrowser.String(), FLAG_HELP_BROWSER)
	flag.BoolVar(&scanOnly, "s", false, FLAG_HELP_SCAN)
	flag.BoolVar(&scanOnly, "scan", false, FLAG_HELP_SCAN)
	flag.StringVar(&profile, "n", "", FLAG_HELP_NAME)
	flag.StringVar(&profile, "name", "", FLAG_HELP_NAME)
	flag.BoolVar(&helpme, "h", false, FLAG_HELP_ME)
	flag.BoolVar(&helpme, "help", false, FLAG_HELP_ME)
	flag.BoolVar(&cacheOnly, "c", false, FLAG_HELP_CACHE)
	flag.BoolVar(&cacheOnly, "cache", false, FLAG_HELP_CACHE)
	flag.BoolVar(&profileOnly, "p", false, FLAG_HELP_PROFILE)
	flag.BoolVar(&profileOnly, "profile", false, FLAG_HELP_PROFILE)
	flag.StringVar(&szmodeS, "size", "Std", FLAG_HELP_SIZE)
	flag.StringVar(&szmodeS, "z", "Std", FLAG_HELP_SIZE)
	flag.BoolVar(&logging, "log", false, FLAG_HELP_LOG)
	flag.BoolVar(&dryRun, "dry", false, FLAG_HELP_DRYRUN)
	flag.Parse()

	// B. Validation

	// (b.1) Help!
	if helpme {
		help()
		os.Exit(0)
	}

	// (b.2) Target Profile name (-name)
	if !scanOnly && len(profile) == 0 {
		help()
		die(1, "Need profile directory base name")
	}

	// (b.3) No -cache nor -profile is same as ALL
	if !cacheOnly && !profileOnly {
		cacheOnly = true
		profileOnly = true
	}

	// (b.4) Browser capabilities
	var browser browsers.Browser
	switch strings.ToLower(browserName) {
	case "chromium": // default
		browser = browsers.ChromiumBrowser
		break
	case "firefox":
		browser = browsers.FirefoxBrowser
		break
	default:
		die(2, "Not a supported browser %q", browserName)
	}

	// (b.5) Size reporting mode
	var sizeMode cmn.SizeMode
	switch strings.ToLower(szmodeS) {
	case "si":
		sizeMode = cmn.SizeModeSI
		break
	case "iec":
		sizeMode = cmn.SizeModeIEC
		break
	case "std":
		sizeMode = cmn.SizeModeStd
		break
	default:
		die(3, "Invalid size mode (SI|IEC|STD) %q", szmodeS)
	}

	// (b.6) Conditional Logging
	logx = cmn.NewConditionalLogger(logging, "Main")

	// (b.7) Prologue
	if !scanOnly {
		fmt.Printf("Browser name  : %s\n", browser)
		fmt.Printf("Profile name  : %s\n", profile)
		fmt.Printf("Erase cache   : %t\n", cacheOnly)
		fmt.Printf("Erase profile : %t\n", profileOnly)
		fmt.Printf("Size mode     : %s\n", sizeMode)
		fmt.Printf("Logging enable: %t\n", logging)
		if dryRun {
			fmt.Printf("Dry Run enable: %t\n", dryRun)
		}
	}
	//os.Exit(100)
	// C. Execute
	runner := &BrowserWipe{}
	runner.SizeMode = sizeMode

	if scanOnly {
		runner.Scan()
	} else {
		if err := runner.GetCleaner(browser, profile, scanOnly, sizeMode, dryRun); err == nil {
			if code, err := runner.Run(cacheOnly, profileOnly); err != nil {
				die(code, err.Error())
			}
		} else {
			die(4, err.Error())
		}
	}

	// D. Report
	fmt.Println("DONE!!!")
}
