# WipeChromium Developer Notes

We all know the human mind, right? During development we have fresh information,
but after a few months of not maintaining the software, most programmers forget
to a certain degree those fine details. If forgotten, we have to relearn or
perhaps we make wrong assumptions that introduce bugs. I have confirmed this
throughout **many** years as a developer. Therefore, these *Developer Notes*.
These shall also serve for contributors.

## General Notes

* The entire application is written in GO v1.22.1 but AFAIK the lowest possible
  version is GO v1.18.
* The `Makefile` is a mere convenience.

We should be aware that on every computer we have USER accounts that require
login credentials. Then there are the Internet Browsers installed on the system
which *each* USER account can use.

Modern browsers support the so called PROFILES which are sometimes used when
multiple people share a single computer USER account, but in many cases it is
the same USER having several PROFILES with different settings like Work, Personal
and Guest (for illustrative purposes!).

### Application Program

**Source**: `cmd/wiper/*go`

Currently this package has only one deliverable: the wiper executable.

---
## Browsers

It started with only Chromium support but with other browsers in mind. For
that reason, I used an interface defined in `browsers` and underneath that
there is one directory for each supported browser: `chromium` & `firefox`.

Every supported browser:

* Should implement `browsers.IBrowser`,
* Must have an enumerated value of type `browsers.Browser`, and
* An entry in `browsers.SupportedBrowsers`
* A sub-package `browsers/NAME` which implements the browser cleaner.

Each browser sub-package ideally has:

* The main cleaner implementation of `browsers.IBrowser` in a filename named
  after the browser, i.e. `chromium.go`
* One OS-specific file which should ensure it works on Linux/Unix, MacOS and
  Windows, i.e. `chromium_unix.go`, `chromium_darwin.go` & `chromium_windows.go`.

Typically we see that under a root data directory, there are sub-directories
for each **browser** *user profile*. The cleaner should be aware of profile-specific
data that should be preserved.

Keep in mind that the paths used by the cleaner of each browser vary per OS.

* `GetRootDataDir()` gets the root data directory for all the browser settings
  for that computer USER account. Underneath this directory we typically have
  one or more user PROFILES. Also see `IdentifyAppDataRoot()`.
* `GetCacheDir()` gets the *Cache* directory of a PROFILE. Usually we can
  wipe it out in its entirety. In some OS (Unix/Linux) the cache is under a
  different root directory. In others, it *may* be under the PROFILE directory.
  See also `IdentifyProfileCache()`.
* `GetDataDir()` gets the data, settings, bookmarks, etc. directory of a given
  PROFILE. It is usually under `RootDataDir`. See also `IdentifyProfileData()`.
* `GetBROWSERDirs()` is a convenience function that gets both the Data & Cache
  directories (or an error). For example `GetChromiumDirs()`.

There are also a few functions (and methods) paired to the above functions
whose only purpose is to take categorically **identify** that the path is what
it is: a user cache, a user profile, or a browser root directory. Their purpose
is to look for certain files and subdirectories in those Cache/Data/Root
directories that are always present for a given browser. This is to safeguard,
to some extent, that we do not arbitrarily wipe out some other data!

> Dry, Wet & Damp Runs

For those developing support for extra browsers, it is likely that some bugs
might creep into your `IBrowser` cleaner implementation. Since we deal with
cleaning & wiping, a seemingly inocuous bug may result in files or directories
getting deleted when it was not meant to.

For that reason, it is important that -at least initially- you do dry runs
at first so that it only tells you what it is about to do without actually
doing it!

For that reason, instead of using `os.*` calls directly, you should consider
emulating them. Once your cleaner is ready for production, you may choose to
keep it with the filesystem proxy `DryRun` or revert to pure `os.*` calls.
*For more information see the §"Dry Runner" section.

### Chromium

**Source**: `browsers/chromium/*`

The first and principal browser implemented up to `v0.3.1`. It is based on
the `Chromium` browser that is installed on Debian and many other Linux
distributions. Many other browsers are derived from this engine.


### FireFox

**Source**: `browsers/firefox/*`

The 2nd browser I implemented "just because". Since I am using a Raspberry Pi
running Debian, the distributed variant is *Firefox ESR*.

Firefox supports user Profiles which are created via the `about:profiles`.
It stores the list of profiles in a file named `profiles.ini` under the
`RootDataDir`. However, unlike Chromium, the name of the profile is different
than the actual directory name of that profile, and that mapping is declared
in that `profiles.ini` file.

Notice that unlike the Chromium implementation, the `GetFirefoxDirs()` function
receives a *Profile (sub-)Directory* which has already been translated in the
constructor from the provided `ProfileName`

## Internals

### Virtual File System

This application has a small dependency of a Virtual File System module which
supports Dummy, OS and Memory Filesystems. I carefully chose that one because
it does not rely on any other modules, and does what it does without extra
baggage.

For detailed information see [this](https://github.com/blang/vfs)]. In particular,
*WipeChromium* uses the Memory File System for testing `DirCleanerVFS` and
`DryRun` objects.

#### Dry Runner

The *Dry Runner* is a relatively simple filesystem proxy for safe development
of new browser cleaners. But also helpful in running `wipechromium` with a
`-dry` flag so that it doesn't actually removes anything but tells the user
what it **would** remove otherwise in a normal run.

The `DryRun` object has three modes of operation:

* `DryRunTargetNOP`: (Dry run) does not do any `os.Remove*()`, `os.Makedir*()` or
  `os.Rename()` operations on ANY filesystem. It simply replaces those sensitive
  calls with a console message indicating which operation was invoked and on
  which file or directory.
* `DryRunTargetOS`: (Wet run) actually operate on the REAL filesystem. It is
  like using `os.*` directly but with the flexibility of switching easily
  without re-coding.
* `DryRunTargetVFS`: (Damp run) does NOT operate on the real filesystem, but
  instead it works on the supplied memory filesystem.

