package distribution

import (
	"testing"
)

func TestDistribution(t *testing.T) {
	var d *Distribution
	var err error

	d, err = New("id", "foo bar baz")
	if _, ok := err.(*Error); !ok {
		t.Fatal("oops")
	}
	if err == nil {
		t.Fatal("oops")
	}

	d, err = New("id", "CPAN Upload: S/SK/SKAJI/Perl6/App-Mi6-0.1.6.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsRaku {
		t.Fatal("oops")
	}
	if d.MainModule != "App::Mi6" {
		t.Fatal("oops")
	}

	d, err = New("id", "CPAN Upload:S/SK/SKAJI/App-cpm-0.963.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if d.IsRaku {
		t.Fatal("oops")
	}

	d, err = New("id", "S/SK/SKAJI/App-cpm-0.963-TRIAL.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if d.IsRaku {
		t.Fatal("oops")
	}
	if d.Distvname != "App-cpm-0.963-TRIAL" {
		t.Fatal("oops")
	}

	d, err = New("id", "https://cpan.metacpan.org/authors/id/S/SH/SHAY/perl-5.24.4-RC1.tar.bz2")
	if err != nil {
		t.Fatal(err)
	}
	if d.IsRaku {
		t.Fatal("oops")
	}
	if d.Distvname != "perl-5.24.4-RC1" {
		t.Fatal("oops")
	}

}
