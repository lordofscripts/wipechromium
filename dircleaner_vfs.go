/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Recursive Directory cleanup using Virtual File System
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"os"
	"fmt"
	"slices"	// GO v1.18
	"path/filepath"

	"github.com/blang/vfs"
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
	Root	string
	cleanedSize int64
	removedQty	int
	skippedQty	int
	logx		ILogger
	vfs			vfs.Filesystem
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDirCleanerVFS(fs vfs.Filesystem, root string, logger ...ILogger) *DirCleanerVFS {
	const cName = "DirCleaner"
	var logCtx ILogger
	if len(logger) == 0 {
		logCtx = NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}
	return &DirCleanerVFS{root, 0, 0, 0, logCtx, fs}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (d *DirCleanerVFS) String() string {
	return fmt.Sprintf("DirCleaner %q del:%d skip:%d size:%d", d.Root, d.removedQty, d.skippedQty, d.cleanedSize)
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

/*   			if finfo, err := item.Info(); err == nil {
   				// recalculate saved space
				d.cleanedSize += finfo.Size()
				if finfo.IsDir() {
					executor = execRemoveRecursive
				} else {
					executor = execRemoveSingle
				}
   			} else {
				d.logx.Print("DirCleaner WARN:", err)
   			}*/

			// delete
			fullPath := filepath.Join(d.Root, item.Name())
			if err := executor(fullPath); err != nil {
				return err
			} else {
				if item.IsDir() {
					// LIMITATION OF VFS SO FAR
					//d.cleanedSize += GetDirectorySizeVFS(d.vfs, fullPath)
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

func GetDirectorySize(folder string) int64 {
	// Step 1: remember subdirectories we must recurse
	var folders []string
    // Step 2: read directory and handle errors.
    dirRead, err := os.Open(folder)
    if err != nil {
        panic(err)
    }
    dirFiles, err := dirRead.Readdir(0)
    if err != nil {
        panic(err)
    }

    // Step 3: sum up Size of all files in the directory.
    sum := int64(0)
    for _, fileHere := range dirFiles {
		if fileHere.IsDir() {
			// Size() returns the size of the directory entry not the actual
			// sum of file sizes in that directory. Therefore, we skip it for later.
			folders = append(folders, filepath.Join(folder, fileHere.Name()))
		} else {
			csize := fileHere.Size()
		    sum += csize
			//fmt.Printf("%8d %5t %s\n", csize, fileHere.IsDir(), fileHere.Name())
		}
    }

    // Step 4: close directory and return the sum.
    dirRead.Close()

	// Step: 5: Iterate recursively through the directories we encountered
	for _, dn := range folders {
		sum += GetDirectorySize(dn)
	}
    return sum
}


//  UNFORTUNATELY vfs does not support neither Readdir() nor ReadDir()
func GetDirectorySizeVFS(fs vfs.Filesystem, folder string) int64 {
	// Step 1: remember subdirectories we must recurse
	var folders []string
    // Step 2: read directory and handle errors.
/*    dirRead, err := vfs.Open(fs, folder)
    if err != nil {
        panic(err)
    }*/
    //dirFiles, err := dirRead.ReadDir(0)
    dirFiles,err := fs.ReadDir(folder)
    if err != nil {
        panic(err)
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
		sum += GetDirectorySize(dn)
	}
    return sum
}
