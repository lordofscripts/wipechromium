/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"os"
	"fmt"
	"slices"	// GO v1.18
	"path/filepath"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DirCleaner struct {
	Root	string
	cleanedSize int64
	removedQty	int
	skippedQty	int
	logx		ILogger
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDirCleaner(root string, logger ...ILogger) *DirCleaner {
	const cName = "DirCleaner"
	var logCtx ILogger
	if len(logger) == 0 {
		logCtx = NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}
	return &DirCleaner{root, 0, 0, 0, logCtx}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (d *DirCleaner) String() string {
	return fmt.Sprintf("DirCleaner %q del:%d skip:%d size:%d", d.Root, d.removedQty, d.skippedQty, d.cleanedSize)
}

func (d *DirCleaner) CleanUp(exceptions []string) error {
	d.cleanedSize = 0
	d.removedQty = 0
	d.skippedQty = 0
	entries, err := os.ReadDir(d.Root)
    if err != nil {
        d.logx.Print(err)
        return err
    }

	var executor func(string) error = nil

	execRemoveRecursive := func(path string) error {
		return os.RemoveAll(path)
	}

	execRemoveSingle := func(path string) error {
		return os.Remove(path)
	}

    for _, item := range entries {
   		if !slices.Contains(exceptions, item.Name()) {
   			// get file/dir size
   			if finfo, err := item.Info(); err == nil {
   				// recalculate saved space
				d.cleanedSize += finfo.Size()
				if finfo.IsDir() {
					executor = execRemoveRecursive
				} else {
					executor = execRemoveSingle
				}
   			} else {
				d.logx.Print("DirCleaner WARN:", err)
   			}

			// delete
			fullPath := filepath.Join(d.Root, item.Name())
			if err := executor(fullPath); err != nil {
				return err
			}
			d.removedQty += 1
   		} else {
   			d.skippedQty += 1
			d.logx.Printf("DirCleaner skipping %s", item.Name())
   		}
    }
    return nil
}

func (d *DirCleaner) CleanedSize() int64 {
	return d.cleanedSize
}


/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/


