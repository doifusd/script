#!/bin/sh

rm "/Users/$(whoami)/Library/Application Support/Beyond Compare/registry.dat"

rm /Applications/Beyond\ Compare.app/Contents/MacOS/BCompare/BCompare
rm /Applications/Beyond\ Compare.app/Contents/MacOS/BCompare/BCompare

cp /Applications/Beyond\ Compare.app/Contents/MacOS/BCompare.real /Applications/Beyond\ Compare.app/Contents/MacOS/BCompare
