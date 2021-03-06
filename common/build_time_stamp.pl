#!/usr/bin/env perl
# -*- mode: cperl; coding: utf-8; -*-
# /home/krylon/go/src/pepper/common/build_time_stamp.pl
# created at 22. 07. 2019 by Benjamin Walkenhorst
# (c) 2019 Benjamin Walkenhorst <krylon@gmx.net>
# Time-stamp: <2020-11-23 20:27:22 krylon>
#  Redistribution and use in source and binary forms, with or without
#  modification, are permitted provided that the following conditions
#  are met:
#  1. Redistributions of source code must retain the copyright
#     notice, this list of conditions and the following disclaimer.
#  2. Redistributions in binary form must reproduce the above copyright
#     notice, this list of conditions and the following disclaimer in the
#     documentation and/or other materials provided with the distribution.
#
#  THIS SOFTWARE IS PROVIDED BY BENJAMIN WALKENHORST ``AS IS'' AND
#  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
#  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
#  ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
#  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
#  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
#  OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
#  HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
#  LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
#  OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
#  SUCH DAMAGE.

use strict;
use warnings;
use diagnostics;
use utf8;
use feature qw(say);

use Carp;
use English '-no_match_vars';
use Time::Piece;

my $now = localtime;
my $year = $now->year;
my $mon = $now->mon;
my $day = $now->mday;
my $hour = $now->hour;
my $min = $now->min;
my $sec = $now->sec;

my $outpath = "build_stamp_gen.go";

open(my $fh, ">", $outpath)
  or croak "Error opening $outpath: $OS_ERROR";

print {$fh} <<"EOF";
// Code generated by build_time_stamp.pl. DO NOT EDIT.
package common

import "time"

// BuildStamp is the time and date when the program was last built.
var BuildStamp = time.Date($year, $mon, $day, $hour, $min, $sec, 0, time.Local)

EOF

close $fh;

# Local Variables: #
# compile-command: "perl -c /home/krylon/go/src/pepper/common/build_time_stamp.pl" #
# End: #
