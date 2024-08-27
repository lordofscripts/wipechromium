# Lord of Scripts&trade; WipeChromium

[![Years](https://badges.pufler.dev/years/lordofscripts)](https://badges.pufler.dev)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/wipechromium?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/wipechromium)
[![GitHub](https://img.shields.io/github/license/lordofscripts/wipechromium)](https://github.com/lordofscripts/wipechromium/blob/master/LICENSE)
![Tests](https://github.com/lordofscripts/wipechromium/actions/workflows/crossbuild.yml/badge.svg)
[![Coverage](https://coveralls.io/repos/github/lordofscripts/wipechromium/badge.svg?branch=main)](https://coveralls.io/github/lordofscripts/wipechromium?branch=main)
[![Visits](https://badges.pufler.dev/visits/lordofscripts/wipechromium)](https://badges.pufler.dev)
[![Created](https://badges.pufler.dev/created/lordofscripts/wipechromium)](https://badges.pufler.dev)
[![Updated](https://badges.pufler.dev/updated/lordofscripts/wipechromium)](https://badges.pufler.dev)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/wipechromium)

*WipeChromium* is a small utility written in Go (v1.22) whose purpose is doing
something I need to do often: wipe out my Chromium browser data while keeping
only the important stuff such as **Bookmarks, Settings, Extensions & Web Applications**

Usually you have to use menu options that move from time to time or are too
cumbersome to find. Additionally, people tend not to know or remember the
command-line options of each browser. With this application there are a few
simple options that allow you to regain greater control of your privacy and of
your disk space.

*Note:* Like all software, you agree to use it at your OWN risk. **I am NOT
responsible for any data loss**.

Use Cases:

* You want to *free up space* in your disk
* You worry about *scripts* that may have remained after your browsing
* You are tired of big companies profitting from you with *tracking cookies*
* You want to *start afresh* without reinstalling everything.
* You want to protect your *privacy*
* You are worried about security and want to make sure any usernames/passwords
  are removed from your profile.
* You want erase all the history and whatever else browser-makers track!

For example, on a test run my user profile directory diminished from 126 MB to
only 234 KB in size while keeping all important data.

## Requirements

Go Version: >= v1.18

## Installation

> `go install github.com/lordofscripts/wipechromium`

## Usage

Before continuing keep in mind that regardless of your OS, you, as a user, will
run within the context of your *user account*, for example `lordofscripts` which
is on `/home/lordofscripts` on Unix/Linux, or `C:\Users\LordOfScripts\` on
Microsoft Windows, or `/Users/LordOfScripts` on MacOS.

When you use most modern browsers,such as Chromium-based browsers, within your
user account (i.e. `lordofscripts`), you can use the browser with different
*user profiles* such as `Profile 1`, `Profile Work`, `Profile Gaming`, although
most users just use their `Default` user profile. This software operats on
these **User Profiles**.

### Features
* At present it supports *Chromium* but it is designed to support extra browsers.
* It can `-scan` your system for browser data & cache directories.
* You can wipe out your entire cache,
* You can wipe out most of your user profile data,
* You can wipe the cache & data in one go,
* It keeps your precious data: Settings, Web applications, File systems, Bookmarks & Extensions.

#### Known Limitations

Although care has been taken to support file operations for Linux/Unix, MacOS
and Windows, I have only tested it with Linux/Unix.


### Useful combinations

If you think there is a malfunction of some kind, you can add the `-log` flag
to any of the command variations mentioned below. Else it is better to keep
logging off (default). But In general:

> `wipechromium [-browser Chromium][-log] {-scan|-cache|-profile} {-name NAME}`

Notice that the `-browser` option takes as parameter the browser type. Currently
it defaults to Chromium because it's the only one supported at the moment. As
such, you can omit this option in your command-line for now.

#### Get Help

Gives you a comprehensive guide of all command-line options:

> `wipechromium -help`

#### Scan for browsers

Scans your home directory for data and/or cache of supported browsers:

> `wipechromium -scan`

I actually **recommend** that you run with this option the first time before
trying anything else. If it says it detected Chromium (or other browser's) data,
it will tell you. If it doesn't detect it there is no purpose in running the
other commands.

#### Clear Profile's Cache

Let's say your gaming profile is `Dart Vader` and that it has grown big. Or you
just want to feel more secure:

> `wipechromium -browser Chromium -name "Dart Vader" -cache`

This command will wipe out the entire Cache for the named user profile.

#### Clear Profile's User Data

Let's say your daily profile is `Profile X` and every now and then (you should!)
you like to clean it up in case there are possible malware scripts, or those
pesky tracking cookies from greedy companies that want to abuse your privacy.

But wait a minute! you have bookmarks of useful websites right? maybe you
installed some browser extensions, or perhaps you installed one of those
**Progressive Web Applications** and you certainly don't want to reinstall those
things. And especially PGA's have their own File System in case the Cloud is not
available. You don't want to lose those! Well, I don't, and this script takes
care of keeping that data intact while wiping out the rest in that profile.
After all, next time you fire your browser under that profile, it will recreate
them and you start afresh.)

> `wipechromium -browser Chromium -name "Profile X" -profile`

This command will wipe out the entire Cache for the named user profile.

#### Clear Both User Profile Data & Cache

This option is equivalent to using both `-cache` and `-profile` together. Like
this:

> `wipechromium -browser Chromium -name "Profile X" -cache -profile`

However, since usually you would want to do both, the internal logic enables
*both* if you don't set any; therefore, it is equivalent to this:

> `wipechromium -browser Chromium -name "Profile X"`

which will clean up both the profile data and the profile cache in one run.

-----
> All Rights Reserved [LordOfScripts&trade;](https://allmylinks.com/lordofscripts)

