#!/bin/sh

cd $GOPATH/src/github.com/blicero/dwd/

rm -vf bak.dwd dwd dbg.build.log && du -sh . && git fsck --full && git reflog expire --expire=now && git gc --aggressive --prune=now && du -sh .

