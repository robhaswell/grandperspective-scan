# gpscan
A utility to create GrandPerspective scan files on other operating systems.

## What does it do?

[GrandPerspective](http://grandperspectiv.sourceforge.net/) is a small utility application for Mac that graphically shows the disk usage within a file system.

![Screenshot of GrandPerspective](resources/screenshot.png?raw=true "GrandPerspective")

`gpscan` is a cross-platform utility which can be run on almost any operating which will create scan data files which can be loaded by GrandPerspective to visualize disk usage on non-OSX machines.

![Screenshot of GrandPerspective's File menu demonstrating the Load Scan Data option](resources/load.png?raw=true "Load Save Data")

The Go source code can be cross-compiled for a large selection of platforms.

## Why was it created?

I wanted to use GrandPerspective (or any SequoiaView-clone) to visualize disk usage on a Linux server that I administrate.
Unfortunately I only have access to Windows and OSX GUI environments.

## How do I use it?

`gpscan` takes two arguments: The directory to scan, and a filename to save the result to.
It is recommended that you use a `.gpscan` file extension so that the file is automatically opened with GrandPerspective.
You can pass `-` as an output file in order to output to standard output.

```bash
$ gpscan path/to/dir scan.gpscan
```

## How do I install it?

Download a relevant release, unzip and copy to a location in your `$PATH`, e.g. `/usr/local/bin/gpscan`.

# Project status

This project is developed only to the point where I could complete my original task, however I will gladly review pull-requests and take bug reports in the Issues section.
