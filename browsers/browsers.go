/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Common stuff for ALL browsers supported by the wipe application.
 *-----------------------------------------------------------------*/
package browsers

import (
	"log"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ChromiumBrowser Browser = iota

)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Browser uint

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// In this 1st release, only Chromium browser is supported. However, the
// software is designed as such that, under the browsers package, we can
// implement cleaners a-la-chromium.ChromiumCleaner which do the same for
// those specific browsers. All you have to do is let the cleaner package
// under "lordofscripts/wipechromium/browsers" implement this interface.
type IBrowsers interface {
	Name() Browser
	String() string

	// Clears a user profile and/or cache.
	// Returns: error (or nil) and if error, an error code
	ClearProfile(doCache, doProfile bool) (error, int)
	// Prints out the location of the directories the program
	// thinks (as per configuration) it should use. Should be checked
	// prior to cleaning the first time!
	Tell() bool

	// Browser data for ALL profiles. A user account has ONE browser AppDataRoot,
	// but therein it may have more than one user Profile, each with its
	// different settings, extensions, bookmarks, etc.
	// Returns: true if GetDataDir() is the root of all user account profiles.
	IdentifyAppDataRoot() bool
	// A user profile's cache directory.
	// Returns: true if GetCacheDir()+profile is a valid browser Cache directory.
	IdentifyProfileCache(profile string) bool
	// A user profile specific data (extensions, Bookmarks, etc.).
	// Applies to each of the profile directories under the directory identified
	// by IdentifyAppDataRoot().
	// Returns: true if directory contains browser user profile data & settings.
	IdentifyProfileData(profile string) bool
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (b Browser) String() string {
	id := ""
	switch b {
		case ChromiumBrowser:
			id = "Chromium"
			break
		default:
			log.Print("Unknown browser")
	}
	return id
}

 