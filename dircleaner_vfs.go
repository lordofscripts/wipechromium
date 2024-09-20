/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Recursive Directory cleanup using Virtual File System
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
	"os"
	"path/filepath"
	"slices" // GO v1.18

	"github.com/lordofscripts/vfs"
	"github.com/lordofscripts/vfs/bucketfs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var _ IDirCleaner = (*DirCleanerVFS)(nil)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DirCleanerVFS struct {
	Root        string
	cleanedSize int64
	removedQty  int
	skippedQty  int
	sizeMode    SizeMode
	logx        ILogger
	vfs         vfs.Filesystem
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// NewDirCleanerVFS creates a new (recursive) directory cleaner instance with
// the selected Virtual File System instance.
func NewDirCleanerVFS(fs vfs.Filesystem, root string, sizing SizeMode, logger ...ILogger) *DirCleanerVFS {
	const cName = "DirCleanerVFS"
	var logCtx ILogger
	if len(logger) == 0 {
		logCtx = NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}
	return &DirCleanerVFS{root, 0, 0, 0, sizing, logCtx, fs}
}

// NewDirCleanerDryVFS creates a new (recursive) directory cleaner instance with
// the Bit-Bucket Virtual File System in a mode that allows us to do a Dry Run.
// Note: This is equivalent to the other ctor. but creates the bucketfs internally
// before calling the default ctor.
func NewDirCleanerDryVFS(root string, sizing SizeMode, logger ...ILogger) *DirCleanerVFS {
	fs := bucketfs.Create()
	return NewDirCleanerVFS(fs, root, sizing, logger...)
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Stringer interface
func (d *DirCleanerVFS) String() string {
	return fmt.Sprintf("DirCleaner %q del:%d skip:%d size:%s", d.Root,
		d.removedQty,
		d.skippedQty,
		ReportByteCount(d.cleanedSize, d.sizeMode))
}

func (d *DirCleanerVFS) CleanUp(exceptions []string) error {
	d.cleanedSize = 0
	d.removedQty = 0
	d.skippedQty = 0
	entries, err := d.vfs.ReadDir(d.Root)
	if err != nil {
		d.logx.Print(err)
		return err
	}

	var executor func(string) error = nil

	execRemoveRecursive := func(path string) error {
		return vfs.RemoveAll(d.vfs, path)
	}

	execRemoveSingle := func(path string) error {
		return d.vfs.Remove(path)
	}

	for _, item := range entries {
		if !slices.Contains(exceptions, item.Name()) {
			if item.IsDir() {
				executor = execRemoveRecursive
			} else {
				executor = execRemoveSingle
			}

			fullPath := filepath.Join(d.Root, item.Name())
			if err := executor(fullPath); err != nil {
				return err
			} else {
				if item.IsDir() {
					// LIMITATION OF VFS SO FAR
					folderSize, _ := GetDirectorySizeVFS(d.vfs, fullPath)
					d.cleanedSize += folderSize
				} else {
					// get file size
					d.cleanedSize += item.Size()
				}
			}

			d.removedQty += 1
		} else {
			d.skippedQty += 1
			d.logx.Printf("DirCleaner skipping %s", item.Name())
		}
	}
	return nil
}

func (d *DirCleanerVFS) CleanedSize() int64 {
	return d.cleanedSize
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetDirectorySizeVFS(fs vfs.Filesystem, folder string) (int64, error) {
	// Step 1: remember subdirectories we must recurse
	var folders []string
	// Step 2: read directory and handle errors.
	/*    dirRead, err := vfs.Open(fs, folder)
	      if err != nil {
	          panic(err)
	      }*/
	//dirFiles, err := dirRead.ReadDir(0)
	dirFiles, err := fs.ReadDir(folder)
	if err != nil {
		return 0, err
	}

	// Step 3: sum up Size of all files in the directory.
	sum := int64(0)
	for _, fileHere := range dirFiles {
		if fileHere.IsDir() {
			// Size() returns the size of the directory entry not the actual
			// sum of file sizes in that directory. Therefore, we skip it for later.
			folders = append(folders, filepath.Join(folder, fileHere.Name()))
		} else {
			var csize int64 = fileHere.Size()
			sum += csize
			//fmt.Printf("%8d %5t %s\n", csize, fileHere.IsDir(), fileHere.Name())
		}
	}

	// Step 4: close directory and return the sum.
	//dirRead.Close()

	// Step: 5: Iterate recursively through the directories we encountered
	for _, dn := range folders {
		folderSize, _ := GetDirectorySize(dn)
		sum += folderSize
	}
	return sum, nil
}

// Utility to replicate a real filesystem into the selected Virtual File System.
// Returns: dirCnt, fileCnt, error
func MimicFileSystem(root string, fs vfs.Filesystem) (int64, int64, error) {
	// ensure a directory is made
	cloneDirectoryOnly := func(rootSrc string, all bool) error {
		if fileStat, err := os.Stat(rootSrc); err != nil {
			return err
		} else {
			if all {
				if err := vfs.MkdirAll(fs, rootSrc, fileStat.Mode()); err != nil {
					return err
				}
			} else {
				if err := fs.Mkdir(rootSrc, fileStat.Mode()); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// ensure dummy file is created
	cloneFileOnly := func(fullpath, content string) error {
		if fileStat, err := os.Stat(fullpath); err != nil {
			return err
		} else {
			if err := vfs.WriteFile(fs, fullpath, []byte(content), fileStat.Mode()); err != nil {
				return err
			}
		}
		return nil
	}

	// 1. Create root
	if err := cloneDirectoryOnly(root, true); err != nil {
		return -1, -1, err
	}

	// 2. Walk that root
	var folders []string
	var dirCnt, fileCnt int64

	if dirFiles, err := os.ReadDir(root); err != nil {
		return -1, -1, err
	} else {
		// 2.2 through each file/dir at this level
		for _, fileHere := range dirFiles {
			if fileHere.IsDir() {
				// Size() returns the size of the directory entry not the actual
				// sum of file sizes in that directory. Therefore, we skip it for later.
				dirName := filepath.Join(root, fileHere.Name())
				folders = append(folders, dirName)
				if err := cloneDirectoryOnly(dirName, false); err != nil {
					return -1, -1, err
				}
				dirCnt++
			} else {
				fileName := filepath.Join(root, fileHere.Name())
				if err := cloneFileOnly(fileName, "..."); err != nil {
					return -1, -1, err
				}
				fileCnt++
			}
		}
	}

	// 3: Iterate recursively through the directories we encountered
	for _, dirName := range folders {
		if dC, fC, err := MimicFileSystem(dirName, fs); err != nil {
			return -1, -1, err
		} else {
			dirCnt += dC
			fileCnt += fC
		}
	}
	return dirCnt, fileCnt, nil
}
