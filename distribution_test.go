package main

import "testing"

var body1 = `
The uploaded file

    App-Mi6-0.1.6.tar.gz

has entered CPAN as

  file: $CPAN/authors/id/S/SK/SKAJI/Perl6/App-Mi6-0.1.6.tar.gz
  size: 17560 bytes
   md5: 9a8bc3617c993a9fe735a257d37fe471
  sha1: 52e3a727826a35e5df4a00d4755e31ed700eed7a

CPAN Testers will start reporting results in an hour or so:

  http://matrix.cpantesters.org/?dist=App-Mi6

Request entered by: SKAJI (Shoichi Kaji)
Request entered on: Sun, 25 Mar 2018 04:39:30 GMT
Request completed:  Sun, 25 Mar 2018 04:40:35 GMT

Thanks,
-- 
paused, v1049
`

var body2 = `
The uploaded file

    App-cpm-0.963.tar.gz

has entered CPAN as

  file: $CPAN/authors/id/S/SK/SKAJI/App-cpm-0.963.tar.gz
  size: 39747 bytes
   md5: f7643128ca9fd85ab172c8a4baaa5f92
  sha1: 06f20449745b20955382ebbe3cd1bf4b5e2ec5d0

CPAN Testers will start reporting results in an hour or so:

  http://matrix.cpantesters.org/?dist=App-cpm

Request entered by: SKAJI (Shoichi Kaji)
Request entered on: Sun, 25 Mar 2018 08:23:16 GMT
Request completed:  Sun, 25 Mar 2018 08:24:08 GMT

Thanks,
-- 
paused, v1049
`

func TestDistribution(t *testing.T) {
	var d *Distribution
	var err error

	d, err = NewDistribution("foo bar baz")
	if err == nil {
		t.Fatal("oops")
	}

	d, err = NewDistribution(body1)
	if err != nil {
		t.Fatal(err)
	}
	if !d.IsPerl6 {
		t.Fatal("oops")
	}

	d, err = NewDistribution(body2)
	if err != nil {
		t.Fatal(err)
	}
	if d.IsPerl6 {
		t.Fatal("oops")
	}

	d, err = NewDistribution("S/SK/SKAJI/App-cpm-0.963-TRIAL.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if d.IsPerl6 {
		t.Fatal("oops")
	}
	if d.Distname != "App-cpm" {
		t.Fatal("oops")
	}

	d, err = NewDistribution("https://cpan.metacpan.org/authors/id/S/SH/SHAY/perl-5.24.4-RC1.tar.bz2")
	if err != nil {
		t.Fatal(err)
	}
	if d.IsPerl6 {
		t.Fatal("oops")
	}
	if d.Distname != "perl" {
		t.Fatal("oops")
	}
}
