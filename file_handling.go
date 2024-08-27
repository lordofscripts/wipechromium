/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Get HOME directory and return the path of fileOrDir right underneath.
// If HOME is /home/toering then AtHome(".config") returns /home/toering/.config
func AtHome(fileOrDir string) string {
	atHome := ""
	if home, err := os.UserHomeDir(); err == nil {
		atHome = filepath.Join(home, fileOrDir)
	} else {
		log.Fatal("Could not get HOME directory")
	}
	return atHome
}

// Check whether path exists and it is a directory
func IsDirectory(path string) bool {
	if finfo, err := os.Stat(path); err == nil {
		return finfo.IsDir()
	}
	// not exist or not a directory
	return false
}

func IsFile(path string) TriState {
	if finfo, err := os.Stat(path); err == nil {
		if !finfo.IsDir() {
			return Yes
		} else {
			return No
		}
	}
	// not exist or not a directory
	return Undecided
}

func MoveDir(src, dest string) error {
	// ensure both are directories
	if !IsDirectory(src) {
		return fmt.Errorf("Not a directory: %s", src)
	}
	if !IsDirectory(dest) {
		return fmt.Errorf("Not a directory: %s", dest)
	}

	// Rename will do Move IF both are directories
	if err := os.Rename(src, dest); err != nil {
		return err
	}

	log.Printf("Moved directory to %s", dest)
	return nil
}

// Move and optionally rename a file.
// Example: moveFile("code/paragraph.go", "code/text/text_handling.go")
func MoveFile(src, dest string, notify bool) error {
	// In GO rename allows to rename and move a file in one step
	if err := os.Rename(src, dest); err != nil {
		log.Print("REN-F", src)
		log.Print("REN-T", dest)
		return err
	}

	if notify {
		log.Printf("Moved %s", src)
	}
	return nil
}

// Same as moveFile except the destination is a directory rather than a
// fully-qualified filename
// Example: moveFileTo("/home/pi/delete.me", "/tmp")
func MoveFileTo(src, dest string, notify bool) error {
	if !IsDirectory(dest) {
		return fmt.Errorf("Not a directory: %s", dest)
	}
	//filename := filepath.Base(src)
	//destFull := filepath.Join(dest, filename)
	//return moveFile(src, destFull, true)
	return MoveFile(src, dest, true)
}

// Example: changePath("/home/pi/test.sh", "/tmp/anydir")
func ChangePath(src, dest string) string {
	base := filepath.Base(src)
	result := filepath.Join(dest, base)
	return result
}

// Removes all files matching a Pattern at Dir.
// Example: removeWithPattern("/home/lordofscripts/.cache", "*.log")
func RemoveWithPattern(dir, pattern string) error {
	glob := dir + string(os.PathSeparator) + pattern
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	for _, fname := range files {
		if err := os.Remove(fname); err != nil {
			return err
		}
	}

	log.Printf("Deleted %s on %s", pattern, dir)
	return nil
}

func MoveWithPattern(dir, pattern string) error {
	glob := dir + string(os.PathSeparator) + pattern
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	for _, fname := range files {
		if err := os.Remove(fname); err != nil {
			return err
		}
	}

	log.Printf("Deleted %s on %s", pattern, dir)
	return nil
}
