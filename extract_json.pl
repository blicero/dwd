#!/usr/bin/perl
# -*- mode: cperl; coding: utf-8; -*-
# /home/krylon/go/src/github.com/blicero/dwd/extract_json.pl
# created at 26. 07. 2021 by Benjamin Walkenhorst
# (c) 2021 Benjamin Walkenhorst <krylon@gmx.net>
# Time-stamp: <2021-07-26 20:25:55 krylon>
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

# Extract and format the JSON data from the response from DWD's web service

sub slurp {
  my ($file) = @_;
  open my $fh, '<', $file or die "Cannot open $file: $OS_ERROR";
  local $/ = undef;
  my $cont = <$fh>;
  close $fh;
  return $cont;
}

sub xtract {
  my ($content) = @_;
  if ($content =~ /^warnWetter[.]loadWarnings\((.*)\);/) {
    return $1;
  } else {
    return;
  }
}

foreach my $file (@ARGV) {
  my $content = slurp($file);
  my $json = xtract($content);
  my $dumpfile = "$file.x.json";
  open(my $fh, '>', $dumpfile)
    or die "Cannot open $dumpfile: $OS_ERROR";

  print {$fh} $json;
  print {$fh} "\n";
  close $fh;

  system "json_pp < $dumpfile > $file";

  unlink $dumpfile;
}

# Local Variables: #
# compile-command: "perl -c /home/krylon/go/src/github.com/blicero/dwd/extract_json.pl" #
# End: #
