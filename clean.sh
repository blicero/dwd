#!/bin/sh
# Time-stamp: <2021-07-30 18:06:08 krylon>

cd $GOPATH/src/github.com/blicero/dwd/

rm -vf bak.dwd dwd dbg.build.log && \
    du -sh . && \
    git fsck --full && \
    git reflog expire --expire=now && \
    git gc --aggressive --prune=now && \
    du -sh .

