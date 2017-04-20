package printoo_test

import (
	"strings"
	"testing"

	"github.com/innermond/printoo/printoo"
)

func TestLongname(t *testing.T) {
	ds := map[string][]map[string][]string{
		"empty cases": {
			{"longname": {""}},
			{"longname": {"     "}},
		},
		"length lower": {
			{"longname": {"<3x"}},
		},
		"length greather": {
			{"longname": {strings.Repeat("x", 51)}},
		},
		"only printable": {
			{"longname": {"x\nx"}},
			{"longname": {"x\tx"}},
			{"longname": {"contain \nx\007x"}},
		},
		"required": {
			{"non-longname": {"test"}},
			{"id": {"1"}},
		},
	}
	t.Log("Invalid data should not be validated")
	for tc, td := range ds {
		t.Run(tc, func(t *testing.T) {
			for _, d := range td {
				_, ok, verr := printoo.NewPerson(d)
				if ok {
					t.Errorf("%s: %s", tc, "WRONG data accepted")
				}
				t.Log(verr)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	ds := map[string][]map[string][]string{
		"empty cases": {
			{"phone": {""}},
		},
		"length lower": {
			{"phone": {strings.Repeat("0", 9)}},
		},
		"length greather": {
			{"phone": {strings.Repeat("0", 16)}},
		},
		"not a phone": {
			{"phone": {"abc"}},
		},
	}
	t.Log("Invalid data should not be validated")
	for tc, td := range ds {
		t.Run(tc, func(t *testing.T) {
			for _, d := range td {
				_, ok, verr := printoo.NewPerson(d)
				if ok {
					t.Errorf("%s: %s", tc, "WRONG data accepted")
				}
				t.Log(verr)
			}
		})
	}
}
