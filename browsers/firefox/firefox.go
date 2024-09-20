/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package firefox

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmn "github.com/lordofscripts/wipechromium"
	"github.com/lordofscripts/wipechromium/browsers"

	"github.com/go-ini/ini"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	CODENAME              = "Alpha"
	PERMS                 = 0700
	RecreateCacheDir bool = true
)

var (
	_ browsers.IBrowsers = (*FirefoxCleaner)(nil)

	ErrFirefoxCleaner = errors.New("FirefoxCleaner error :(")

	// Don't delete these on Firefox root, i.e. ~/.mozilla/firefox/
	FirefoxExceptions []string = []string{ // top-level exceptions (add user profiles)
		"installs.ini",
		"profiles.ini",
	}
	// Don't delete these on Profile root, i.e. ~/.mozilla/firefox/{profile name}/
	FirefoxProfileExceptions []string = []string{ // user profile exceptions
		"extension-preferences.json",
		"extensions.json",
		"lock",
		"places.sqlite",
		// {profile_name}/directories
		"bookmarkbackups",
		"extensions",
		"security_state",
		"settings",
		"features",
	}
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type FirefoxCleaner struct {
	Class       browsers.Browser
	ProfileName string
	CacheRoot   string
	ProfileRoot string
	Profiles    map[string]firefoxProfile
	cleanedSize int64
	sizeMode    cmn.SizeMode
	doDryRun    bool
	scanOnly    bool
	logx        cmn.ILogger
}

// A [Profile*] section in Firefox's profiles.ini
type firefoxProfile struct {
	Name      string
	SubPath   string
	IsDefault bool
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewFirefoxCleaner(profile string, scanOnly bool, smode cmn.SizeMode, dry bool, logger ...cmn.ILogger) *FirefoxCleaner {
	const cName = "FirefoxCleaner"
	var logCtx cmn.ILogger
	if len(logger) == 0 {
		logCtx = cmn.NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}

	// find out which Firefox user profiles are defined
	err, mapping := getProfiles()
	if err != nil {
		return nil
	}

	var subPath string = "" // empty if not profile-specific
	if !scanOnly {
		// translate profile name to profile sub-path
		profile = strings.ToLower(profile)
		pinfo, ok := mapping[profile]
		if !ok {
			cmn.SpitOutError(1, cmn.ErrProfileDoesNotExist)
			return nil
		}
		subPath = pinfo.SubPath
	}

	err, dataDir, cachesDir := GetFirefoxDirs(subPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &FirefoxCleaner{
		browsers.FirefoxBrowser,
		strings.Trim(profile, " \t"),
		cachesDir, //filepath.Join(cachesDir, subPath),
		dataDir,   //filepath.Join(dataDir, subPath),
		mapping,
		0,
		smode,
		dry,
		scanOnly,
		logCtx,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (f firefoxProfile) String() string {
	return fmt.Sprintf("- %15q %-30s Default:%-5t", f.Name, f.SubPath, f.IsDefault)
}

func (c *FirefoxCleaner) String() string {
	reportedSize := cmn.ReportByteCount(c.cleanedSize, c.sizeMode)

	c.logx.Print("Cleaned ", cmn.AddThousands(c.cleanedSize, ','))
	c.logx.Print("Cleaned ", cmn.ByteCountSI(c.cleanedSize))
	c.logx.Print("Cleaned ", cmn.ByteCountIEC(c.cleanedSize))
	return fmt.Sprintf("%sCleaner %q cleaned %s aka %q", c.Class, c.ProfileName, reportedSize, CODENAME)
}

func (c *FirefoxCleaner) Name() browsers.Browser {
	return browsers.FirefoxBrowser
}

// Top level function to clear a Chromium user profile directory. Rather than
// saving important data to a Temp directory and then restoring (as previous version)
// now we simply go through the top level with a list of exceptions
// Example: clearProfile("Profile 1")
func (c *FirefoxCleaner) ClearProfile(doCache, doProfile bool) (error, int) {
	if c.scanOnly {
		return browsers.ErrInvalidOperation, 0
	}
	fmt.Printf("Clearing profile %q\n", c.ProfileName)

	c.cleanedSize = 0
	if len(c.ProfileName) == 0 {
		// we can only operate in AppData root only as no profile is given
		return cmn.ErrNoProfile, 40
	}

	// 1. Profile Cache
	if doCache {
		if err := c.clearCache(); err != nil {
			return err, 50
		}
	}

	// 2. Profile Data
	if doProfile {
		if err := c.eraseProfile(); err != nil {
			return err, 60
		}

		if err := c.clearExtensions(); err != nil {
			return err, 70
		}
	}

	c.logx.Printf("Profile %q cleared of private/junk data", c.ProfileName)
	return nil, 0
}

// This function should be implemented in all wiper browser plugins.
// It should print out the supposed location of the Data & Cache directories
// so that the user can verify prior to running the program for the 1st time.
func (c *FirefoxCleaner) Tell() bool {
	fmt.Printf("❋✦ %s Directories:\n", VARIANT)

	// the profile's existence has been verified at the constructor, and
	// the name had been normalized to lowercase by getProfiles()

	var profileSubDir string = "" // all profiles
	if !c.scanOnly {
		// user-profile specific
		profileSubDir = c.Profiles[strings.ToLower(c.ProfileName)].SubPath
	}

	if err, dataDir, cachesDir := GetFirefoxDirs(profileSubDir); err != nil {
		cmn.SpitOutError(1, err)
		return false
	} else {
		dataExists := cmn.IsDirectory(dataDir)
		cacheExists := cmn.IsDirectory(cachesDir)

		var sizeD, sizeC int64
		if dataExists {
			sizeD, _ = cmn.GetDirectorySize(dataDir)
		}
		if cacheExists {
			sizeC, _ = cmn.GetDirectorySize(cachesDir)
		}

		fmt.Printf("\tData : %5t %s %s\n", dataExists, dataDir, cmn.ReportByteCount(sizeD, c.sizeMode))
		fmt.Printf("\tCache: %5t %s %s\n", cacheExists, cachesDir, cmn.ReportByteCount(sizeC, c.sizeMode))
		/*		// list all registered Firefox user profiles
				fmt.Println("\tProfiles:")
				for _, pe := range c.Profiles {
					fmt.Printf("\t%s\n", pe)
				}*/
		return dataExists && cacheExists
	}
}

// Find all known user profiles. In Firefox ESR these are found in an INI file,
// therefore we do not need to scan a directory looking for profile directories.
func (c *FirefoxCleaner) FindProfileNames() ([]string, error) {
	names := make([]string, 0)

	err, mapping := getProfiles()
	if err != nil {
		return names, err
	}

	for k, v := range mapping {
		entry := k
		if v.IsDefault {
			entry = entry + " (default)"
		}
		names = append(names, entry)
	}

	return names, nil
}

// Browser data for ALL profiles. A user account has ONE browser AppDataRoot,
// but therein it may have more than one user Profile, each with its
// different settings, extensions, bookmarks, etc.
// Returns: true if GetDataDir() is the root of all user account profiles.
func (c *FirefoxCleaner) IdentifyAppDataRoot() bool {
	return IdentifyAppDataRoot()
}

// A user profile's cache directory.
// Returns: true if GetCacheDir()+profile is a valid browser Cache directory.
func (c *FirefoxCleaner) IdentifyProfileCache(profileDir string) bool {
	return IdentifyProfileCache(profileDir)
}

// A user profile specific data (extensions, Bookmarks, etc.).
// Applies to each of the profile directories under the directory identified
// by IdentifyAppDataRoot().
// Returns: true if directory contains browser user profile data & settings.
func (c *FirefoxCleaner) IdentifyProfileData(profileDir string) bool {
	return IdentifyProfileData(profileDir)
}

/* ----------------------------------------------------------------
 *					I n t e r n a l 	M e t h o d s
 *-----------------------------------------------------------------*/

// Clears the entire cache dir of a profile
// Example: clearCache("Profile 1")
// NOTE: Supports dry run.
func (c *FirefoxCleaner) clearCache() error {
	fmt.Println("\tClearing cache...")

	dry := cmn.NewDryRunner()
	if !c.doDryRun {
		dry.Disable()
	}

	cacheSize, err := cmn.GetDirectorySize(c.CacheRoot)
	if err != nil {
		c.logx.Printf("clearCache WARN %s", err)
	}

	if !IdentifyProfileCache(c.Profiles[c.ProfileName].SubPath) {
		c.logx.Printf("%s: %s", cmn.ErrNotBrowserCache, c.CacheRoot)
		return cmn.ErrNotBrowserCache
	}

	// 'Cache' 'Code Cache' and sometimes 'Storage'
	if err := dry.RemoveAll(c.CacheRoot); err != nil {
		werr := cmn.WrapError(err, 41, "Could not remove cache dir %q.\n\t%s", c.CacheRoot, cmn.ThisLocation(1))
		cmn.SpitOutError(1, werr)
		return werr
	}

	if RecreateCacheDir {
		if err := dry.MkDir(c.CacheRoot, PERMS); err != nil {
			c.logx.Printf("Could not recreate CACHE %s %s", c.CacheRoot, cmn.ThisLocation(1))
		}
	}

	c.cleanedSize += cacheSize
	fmt.Printf("\tDeleted %s bytes from cache\n", cmn.ReportByteCount(cacheSize, c.sizeMode))
	c.logx.Print("clearCache DONE")
	return nil
}

// erases a User Profile but keeps important profile data such as
// extensions and settings.
func (c *FirefoxCleaner) eraseProfile() error {
	fmt.Println("\tClearing profile")

	// (a )Identify it is a profile directory
	if !IdentifyProfileData(c.Profiles[c.ProfileName].SubPath) {
		return cmn.ErrNotBrowserProfile
	}

	// (b) we are going to clean the profile's top level
	fmt.Printf("%s DirCleanerRoot %s\n", cmn.ThisLocation(1), c.ProfileRoot)
	var filter cmn.IDirCleaner
	if !c.doDryRun {
		filter = cmn.NewDirCleaner(c.ProfileRoot, c.sizeMode, false)
	} else {
		filter = cmn.NewDirCleanerDryVFS(c.ProfileRoot, c.sizeMode, c.logx)
	}

	// (c) except these important profile items
	err := filter.CleanUp(FirefoxProfileExceptions)
	if err != nil {
		c.logx.Print(err)
		return cmn.WrapError(err, 51, "EraseProfile fault.")
	}

	c.cleanedSize += filter.CleanedSize()
	fmt.Printf("\t...Erased %s bytes\n", cmn.ReportByteCount(c.cleanedSize, c.sizeMode))
	return nil
}

// Apparently nothing to clear in Firefox Extensions
func (c *FirefoxCleaner) clearExtensions() error {
	return nil
}

// Removes all files matching a Pattern at Dir.
// Example: removeWithPattern("/home/lordofscripts/.cache", "*.log")
func (c *FirefoxCleaner) removeWithPatterns(dir string, patterns []string) error {
	dry := cmn.NewDryRunner()
	if !c.doDryRun {
		dry.Disable()
	}

	for _, pattern := range patterns {
		glob := dir + string(os.PathSeparator) + pattern
		files, err := filepath.Glob(glob)
		if err != nil {
			cmn.SpitOutError(1, err)
			return err
		}

		for _, fname := range files {
			var fileSize int64 = 0
			if finfo, err := os.Stat(fname); err != nil {
				return err
			} else {
				fileSize = finfo.Size()
			}
			// remove file or empty directory
			if err := dry.Remove(fname); err != nil {
				return err
			}
			c.cleanedSize += fileSize
		}
	}

	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Identify directory (from GetDataDir()) as a Chromium user account directory.
// Remember every user can have several profiles. This identifies just the
// root where the profiles are located
// At least 'System Profile', 'Avatars' & 'Safe Browsing' exist
func IdentifyAppDataRoot() bool {
	err, appdata := GetRootDataDir()
	if err != nil {
		fmt.Println(err) // TODO refactor app to pass error back to caller!
		return false
	}

	return cmn.IsDirectory(filepath.Join(appdata, "firefox-mpris")) &&
		cmn.IsDirectory(filepath.Join(appdata, "Crash Reports")) &&
		cmn.IsDirectory(filepath.Join(appdata, "Pending Pings")) &&
		cmn.IsFile(filepath.Join(appdata, "installs.ini")) == cmn.Yes &&
		cmn.IsFile(filepath.Join(appdata, "profiles.ini")) == cmn.Yes
}

// Identify GetCacheDir() as a proper Chromium Cache directory
// Both 'Cache' & 'Code Cache' dirs exist
// Linux: .cache/mozilla/firefox/
func IdentifyProfileCache(profileDir string) bool {
	err, userCacheDir := GetCacheDir(profileDir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return cmn.IsDirectory(filepath.Join(userCacheDir, "cache2")) &&
		cmn.IsDirectory(filepath.Join(userCacheDir, "startupCache"))
}

// At least Bookmarks exist && Cookies
func IdentifyProfileData(profileDir string) bool {
	err, userDir := GetDataDir(profileDir)
	if err != nil {
		fmt.Println(err) // TODO refactor app to pass error back to caller!
		return false
	}

	return cmn.IsDirectory(filepath.Join(userDir, "bookmarkbackups")) &&
		cmn.IsDirectory(filepath.Join(userDir, "extensions")) &&
		cmn.IsFile(filepath.Join(userDir, "places.sqlite")) == cmn.Yes &&
		cmn.IsFile(filepath.Join(userDir, "cookies.sqlite")) == cmn.Yes
}

// Gets the list of Firefox user profiles and their mapping to an actual
// directory. Unlike Chromium, Firefox uses a UNIQUE_ID.ProfileName format
// for their user-profile directories. The maping between profile name and
// that directory is on ¿FireFoxAppDir?/profiles.ini
// NOTE: The profile Name we normalize to lowercase.
func getProfiles() (error, map[string]firefoxProfile) {
	const (
		PROFILES_INI = "profiles.ini"
	)
	profiles := make(map[string]firefoxProfile, 0)

	// 1. Find loction
	err, pathStr := GetRootDataDir()
	if err != nil {
		return err, nil
	}

	// 1.1 Location of INI
	iniFilename := filepath.Join(pathStr, PROFILES_INI)
	if cmn.IsFile(iniFilename) != cmn.Yes {
		return fmt.Errorf("Couldn't find FIREFOX %q", iniFilename), nil
	}

	// 1.2 open INI with section & key names normalized to lowercase
	pcfg, err := ini.InsensitiveLoad(iniFilename)
	if err != nil {
		fmt.Printf("FIREFOX Fail to read file: %v", err)
		return err, nil
	}

	// 2. Iterate over INI Sections titled 'ProfileN' where N is a number
	for _, sectionName := range pcfg.SectionStrings() {
		if strings.HasPrefix(sectionName, "profile") {
			const (
				KEY_NAME    = "name"
				KEY_PATH    = "path"
				KEY_DEFAULT = "default"
			)
			section := pcfg.Section(sectionName)

			// 2.1 Read Name, Path (what we want) and Default
			name := section.Key(KEY_NAME).String()
			name = strings.ToLower(name)
			entry := firefoxProfile{
				name,
				section.Key(KEY_PATH).String(),
				section.Key(KEY_DEFAULT).MustBool(false),
			}
			profiles[name] = entry
		}
	}

	return nil, profiles
}
