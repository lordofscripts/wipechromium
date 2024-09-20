/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Dry, Wet & Damp Runner for sensitive filesystem operations
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/lordofscripts/vfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// NO-OP a pure dry run where nothing happens
	DryRunTargetNOP DryRunTarget = iota
	// Running on real OS
	DryRunTargetOS
	// Running on Virtual Filesystem
	DryRunTargetVFS
)

var (
	ErrDryRunInvalidOperation = errors.New("Invalid Dry Run operation!")
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

// DryRunner operation mode: NOP/DRY, OS/WET or VFS/DAMP
type DryRunTarget int

// Quite simple way to enable/disable some dangerous file operations.
type DryRun struct {
	mu      sync.Mutex
	mode    DryRunTarget
	vfs     vfs.Filesystem
	actions *FileActions
	queries *FileQueries
}

// File/Directory object actions on a filesystem (real or not)
type FileActions struct {
	ActionRemoveAll func(path string) error
	ActionRemove    func(path string) error
	ActionMkDirAll  func(name string, perm os.FileMode) error
	ActionMkDir     func(name string, perm os.FileMode) error
	ActionRename    func(oldpath, newpath string) error
}

// File object queries on a filesystem (real or not)
type FileQueries struct {
	QueryIsFile func(path string) TriState
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Instantiates a new DryRun in which the basic OS file operations are
// replaced by NoOps which simply print what would have been done. So, instead
// of using os.RemoveAll() use dr.RemoveAll() after creating dr := NewDryRunner()
func NewDryRunner() *DryRun {
	dr := &DryRun{mode: DryRunTargetNOP, vfs: nil, actions: nil, queries: nil}
	dr.Enable()
	return dr
}

/* ----------------------------------------------------------------
 *				P u b l i c    M e t h o d s
 *-----------------------------------------------------------------*/

// @implements Stringer interface
func (t DryRunTarget) String() string {
	var str string
	switch t {
	case DryRunTargetNOP:
		str = "NOP"
		break
	case DryRunTargetOS:
		str = "OS"
		break
	case DryRunTargetVFS:
		str = "VFS"
		break
	default:
		str = ""
	}
	return str
}

// @implements Stringer interface
func (d *DryRun) String() string {
	return d.mode.String()
}

// Get current operation mode DRY/NOP, WET/OS or DAMP/VFS
func (d *DryRun) GetMode() DryRunTarget {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.mode
}

// True if running the actions will not result in filesystem changes (DRY)
func (d *DryRun) IsSafeRun() bool {
	return d.mode == DryRunTargetNOP
}

// Enable WET run, i.e. on the REAL filesystem!
func (d *DryRun) Disable() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.vfs = nil
	d.mode = DryRunTargetOS

	d.actions, d.queries = d.getOsMapping()
}

// Enable DAMP run (on the selected virtual filesystem)
func (d *DryRun) EnableOn(afs vfs.Filesystem) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.vfs = afs
	d.mode = DryRunTargetVFS
	d.actions, d.queries = d.getVfsMapping(afs)
	//d.PrintAddress("ActionRemove", ActionRemove)
	//d.PrintAddress("afs.Remove", afs.Remove)
}

// Enable true DRY run (No Operations)
func (d *DryRun) Enable() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.vfs = nil
	d.mode = DryRunTargetNOP

	d.actions, d.queries = d.getNopMapping()
}

// Delete a directory and all its children action
func (d *DryRun) RemoveAll(path string) error {
	return d.actions.ActionRemoveAll(path)
}

// File/Dir delete action
func (d *DryRun) Remove(path string) error {
	return d.actions.ActionRemove(path)
}

// Make directories and all its parents action
func (d *DryRun) MkDirAll(name string, perm os.FileMode) error {
	return d.actions.ActionMkDirAll(name, perm)
}

// Make a parent directory action. Does not ensure parents are made!
func (d *DryRun) MkDir(name string, perm os.FileMode) error {
	return d.actions.ActionMkDir(name, perm)
}

// File/Directory rename action
func (d *DryRun) Rename(oldpath, newpath string) error {
	return d.actions.ActionRename(oldpath, newpath)
}

// helper to alias the VFS method to the Action* signature
func (d *DryRun) RemoveAllVFS(path string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.vfs == nil {
		return ErrDryRunInvalidOperation
	}

	vfs.RemoveAll(d.vfs, path)
	return nil
}

// helper to alias the VFS method to the Action* signature
func (d *DryRun) MkDirAllVFS(name string, perm os.FileMode) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.vfs == nil {
		return ErrDryRunInvalidOperation
	}

	vfs.MkdirAll(d.vfs, name, perm)
	return nil
}

// Whether path is a FILE (Yes) or not (No) or does not exist (Undecided)
func (d *DryRun) IsFile(path string) TriState {
	d.mu.Lock()
	defer d.mu.Unlock()

	// for OS & NOP
	var OsStat func(string) (os.FileInfo, error) = os.Stat
	// for VFS
	if d.vfs != nil {
		OsStat = d.vfs.Stat
	}

	if finfo, err := OsStat(path); err == nil {
		if !finfo.IsDir() {
			return Yes
		} else {
			return No
		}
	}
	// not exist or not a directory
	return Undecided
}

// Handy for testing
func (d *DryRun) Dump() {
	fmt.Printf("((DryRun Dump %q))\n", d.mode)
	fmt.Println("  ☢ Actions") // \u2622 (radioactive)
	d.PrintAddress("RemoveAll", d.actions.ActionRemoveAll)
	d.PrintAddress("Remove", d.actions.ActionRemove)
	d.PrintAddress("MkdirAll", d.actions.ActionMkDirAll)
	d.PrintAddress("Mkdir", d.actions.ActionMkDir)
	d.PrintAddress("Rename", d.actions.ActionRename)

	fmt.Println("  ☢ Queries") // \u2622 (radioactive)
	d.PrintAddress("IsFile", d.queries.QueryIsFile)
}

// recursively show a Virtual File System starting at root.
// Returns: number of directories & files encountered
func (d *DryRun) DumpFS(root string, isRecursing ...bool) (dirCnt, fileCnt int64) {
	if len(isRecursing) == 0 {
		fmt.Printf("((DryRun Dump %q))\n", d.mode)
	}
	fileCnt = int64(0)
	dirCnt = int64(0)
	otherDirs := make([]string, 0)
	// NOTES:
	// 1. os.ReadDir(path) ([]os.DirEntry, error) requires DirEntry.Info()
	// 2. vfs.ReadDir(path) ([]os.FileInfo, error)
	var err error
	var entries []os.FileInfo
	if d.vfs != nil {
		// DryRunTargetVFS
		entries, err = d.vfs.ReadDir(root)
	} else {
		dentries, err1 := os.ReadDir(root)
		if err1 == nil {
			// iterate over []DirEntry to convert to []FileInfo
			for _, de := range dentries {
				if finfo, err2 := de.Info(); err2 != nil {
					err = err2
					break
				} else {
					entries = append(entries, finfo)
				}
			}
		} else {
			err = err1
		}
	}

	if err == nil {
		for _, fso := range entries {
			current := filepath.Join(root, fso.Name())
			if fso.IsDir() {
				otherDirs = append(otherDirs, current)
				dirCnt++
				fmt.Printf("\t\tⅮ %s\n", current)
			} else {
				fileCnt++
				fmt.Printf("\t\tℲ %s\n", current)
			}
		}

		// recurse directories
		for _, dirName := range otherDirs {
			subDirCnt, subFileCnt := d.DumpFS(dirName, true)
			dirCnt += subDirCnt
			fileCnt += subFileCnt
		}
	} else {
		fmt.Println("Error dumping Filesystem", err.Error())
	}

	return dirCnt, fileCnt
}

/* ----------------------------------------------------------------
 *				P r i v a t e    M e t h o d s
 *-----------------------------------------------------------------*/

// Action/Query mappings for OS
func (d *DryRun) getOsMapping() (*FileActions, *FileQueries) {
	var myOsActions *FileActions = &FileActions{
		ActionRemoveAll: os.RemoveAll,
		ActionRemove:    os.Remove,
		ActionMkDirAll:  os.MkdirAll,
		ActionMkDir:     os.Mkdir,
		ActionRename:    os.Rename,
	}
	var myOsQueries *FileQueries = &FileQueries{
		QueryIsFile: IsFile,
	}

	return myOsActions, myOsQueries
}

// Action/Query mappings for VFS
func (d *DryRun) getVfsMapping(afs vfs.Filesystem) (*FileActions, *FileQueries) {
	actions := &FileActions{
		ActionRemoveAll: d.RemoveAllVFS,
		ActionRemove:    afs.Remove,
		ActionMkDirAll:  d.MkDirAllVFS,
		ActionMkDir:     afs.Mkdir,
		ActionRename:    afs.Rename,
	}
	queries := &FileQueries{
		QueryIsFile: d.IsFile,
	}

	return actions, queries
}

// Action/Query mappings for NOP
func (d *DryRun) getNopMapping() (*FileActions, *FileQueries) {
	actions := &FileActions{
		ActionRemoveAll: d.dryRemoveAll,
		ActionRemove:    d.dryRemove,
		ActionMkDirAll:  d.dryMkDirAll,
		ActionMkDir:     d.dryMkDir,
		ActionRename:    d.dryRename,
	}
	queries := &FileQueries{
		QueryIsFile: d.IsFile,
	}

	return actions, queries
}

// NOP equivalent of os.RemoveAll()
func (d *DryRun) dryRemoveAll(path string) error {
	fmt.Printf("\t%c os.RemoveAll %s\n", CHR_HIGHVOLTAGE, FromHome(path))
	return nil
}

// NOP equivalent of os.Remove()
func (d *DryRun) dryRemove(path string) error {
	fmt.Printf("\t%c os.Remove %s\n", CHR_HIGHVOLTAGE, FromHome(path))
	return nil
}

// NOP equivalent of os.MkdirAll()
func (d *DryRun) dryMkDirAll(name string, perm os.FileMode) error {
	fmt.Printf("\t%c os.MkdirAll %s %O\n", CHR_HIGHVOLTAGE, FromHome(name), perm)
	return nil
}

// NOP equivalent of os.Mkdir()
func (d *DryRun) dryMkDir(name string, perm os.FileMode) error {
	fmt.Printf("\t%c os.Mkdir %s %O\n", CHR_HIGHVOLTAGE, FromHome(name), perm)
	return nil
}

// NOP equivalent of os.Rename()
func (d *DryRun) dryRename(oldpath, newpath string) error {
	fmt.Printf("\t%c os.Rename %s -> %s\n", CHR_HIGHVOLTAGE, FromHome(oldpath), FromHome(newpath))
	return nil
}

/* ----------------------------------------------------------------
 *				I n t e r n a l    M e t h o d s
 *-----------------------------------------------------------------*/

func (d *DryRun) PrintAddress(name string, x any) {
	fmt.Printf("\t&%-10s is 0x%0x\n", name, reflect.ValueOf(x).Pointer())
}

// Just in case you want to re-test your selected mappings are OK
func (d *DryRun) AssertMapping(mode DryRunTarget) error {
	var e error
	var actions *FileActions
	var queries *FileQueries

	if mode == DryRunTargetOS {
		actions, queries = d.getOsMapping()
	} else if mode == DryRunTargetNOP {
		actions, queries = d.getNopMapping()
	} else if mode == DryRunTargetVFS {
		actions, queries = d.getVfsMapping(d.vfs)
	} else {
		panic("unsupported DryRun Target")
	}

	if !d.assertA(d.actions.ActionRemoveAll, actions.ActionRemoveAll) {
		e = errors.Join(e, fmt.Errorf("ActionRemoveAll"))
	}
	if !d.assertA(d.actions.ActionRemove, actions.ActionRemove) {
		e = errors.Join(e, fmt.Errorf("ActionRemove"))
	}
	if !d.assertB(d.actions.ActionMkDirAll, actions.ActionMkDirAll) {
		e = errors.Join(e, fmt.Errorf("ActionMkdirAll"))
	}
	if !d.assertB(d.actions.ActionMkDir, actions.ActionMkDir) {
		e = errors.Join(e, fmt.Errorf("ActionMkdir"))
	}
	if !d.assertC(d.actions.ActionRename, actions.ActionRename) {
		e = errors.Join(e, fmt.Errorf("ActionRename"))
	}
	if !d.assertQ1(d.queries.QueryIsFile, queries.QueryIsFile) {
		e = errors.Join(e, fmt.Errorf("QueryIsFile"))
	}

	return e
}

// Works for Remove*
func (d *DryRun) assertA(func1, func2 func(string) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for Mkdir*
func (d *DryRun) assertB(func1, func2 func(string, os.FileMode) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for Rename
func (d *DryRun) assertC(func1, func2 func(string, string) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for QueryIsFile
func (d *DryRun) assertQ1(func1, func2 func(string) TriState) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
