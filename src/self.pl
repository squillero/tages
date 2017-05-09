#!/usr/bin/perl -w

# Permission to use, copy, modify, and/or distribute this software
# for any purpose with or without fee is hereby granted.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

use IPC::Open2;

$| = 1;

@DIRS = <Tages_*>;

for($t1=0; $t1<=$#DIRS; ++$t1) {
    $o1 = $DIRS[$t1];
    for($t2=$t1+1; $t2<=$#DIRS; ++$t2) {
	$o2 = $DIRS[$t2];
	print "Match: $o1 vs. $o2\n";
	$pid1 = open2(*Tages1, *Opponent1, "cd $o1; ./Tages -s=0 TCIAIG_TestStrategy 2>> full.log");
	$pid2 = open2(*Tages2, *Opponent2, "cd $o2; ./Tages -s=0 TCIAIG_TestStrategy 2>> full.log");
	$over = 0;
	while(!$over) {
	    $m1 = <Tages1>;
	    $m2 = <Tages2>;
	    if(defined($m1) and defined($m2)) {
		print ".";
		print Opponent1 $m2;
		print Opponent2 $m1;
	    } else {
		$over = 1;
	    }
	}
	print "\n";
	close(Tages1);
	close(Tages2);
	close(Opponent1);
	close(Opponent2);
	waitpid($pid1, 0);
	waitpid($pid2, 0);
    }
}
