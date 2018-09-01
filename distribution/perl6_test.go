package distribution

import (
	"context"
	"testing"
)

func TestPerl6(t *testing.T) {
	var d *Distribution
	var err error
	var name string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetcher := NewPerl6Fetcher()
	d, _ = New("CPAN Upload: E/EL/ELIZABETH/Perl6/Hash-with-0.0.1.tar.gz")
	name, err = fetcher.FetchName(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Hash-with" {
		t.Fatal("oops")
	}

	d, _ = New("CPAN Upload: J/JN/JNTHN/Perl6/cro-zeromq-0.7.6.tar.gz")
	name, err = fetcher.FetchName(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Cro::ZeroMQ" {
		t.Fatal("oops")
	}
}
