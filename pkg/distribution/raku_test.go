package distribution

import (
	"context"
	"testing"
)

func TestRaku(t *testing.T) {
	var d *Distribution
	var err error
	var name string

	fetcher := NewRakuFetcher()
	d, _ = New(1, "CPAN Upload: E/EL/ELIZABETH/Perl6/Hash-with-0.0.1.tar.gz")
	name, err = fetcher.FetchName(context.Background(), d.MetaURL())
	if err != nil {
		t.Fatal(err)
	}
	if name != "Hash-with" {
		t.Fatal("oops")
	}

	d, _ = New(1, "CPAN Upload: J/JN/JNTHN/Perl6/cro-zeromq-0.7.6.tar.gz")
	name, err = fetcher.FetchName(context.Background(), d.MetaURL())
	if err != nil {
		t.Fatal(err)
	}
	if name != "Cro::ZeroMQ" {
		t.Fatal("oops")
	}
}
