package alphavantage

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

// These tests are offline: they exercise the URI driver's pure string functions
// and the host wiring (classify, locate, resolve), which need no network.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "alphavantage" {
		t.Errorf("Scheme = %q, want alphavantage", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "alphavantage" {
		t.Errorf("Identity.Binary = %q, want alphavantage", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		in, typ, id string
	}{
		{"IBM", "symbol", "IBM"},
		{"AAPL", "symbol", "AAPL"},
		{"GOOGL", "symbol", "GOOGL"},
		{"tesco", "query", "tesco"},
		{"technology", "query", "technology"},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if err != nil || typ != tc.typ || id != tc.id {
			t.Errorf("Classify(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.in, typ, id, err, tc.typ, tc.id)
		}
	}
}

func TestClassifyEmpty(t *testing.T) {
	_, _, err := Domain{}.Classify("")
	if err == nil {
		t.Error("Classify(\"\") = nil error, want error")
	}
}

func TestLocate(t *testing.T) {
	cases := []struct {
		typ, id, want string
	}{
		{"symbol", "IBM", "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=IBM&apikey=demo"},
		{"query", "tesco", "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=demo"},
	}
	for _, tc := range cases {
		got, err := Domain{}.Locate(tc.typ, tc.id)
		if err != nil || got != tc.want {
			t.Errorf("Locate(%q, %q) = (%q, %v), want (%q, nil)", tc.typ, tc.id, got, err, tc.want)
		}
	}
}

func TestLocateUnknownType(t *testing.T) {
	_, err := Domain{}.Locate("unknown", "IBM")
	if err == nil {
		t.Error("Locate(unknown) = nil error, want error")
	}
}

func TestHostWiring(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}

	got, err := h.ResolveOn("alphavantage", "IBM")
	if err != nil || got.String() != "alphavantage://symbol/IBM" {
		t.Errorf("ResolveOn = (%q, %v), want alphavantage://symbol/IBM", got.String(), err)
	}
}
