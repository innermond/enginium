package action_test

import (
	"strconv"
	"testing"

	"github.com/innermond/printoo/printoo"
	"github.com/innermond/printoo/printoo/action"
	"github.com/innermond/printoo/printoo/validation"
)

var tests = map[string][]map[string][]string{
	"simple row": {
		{"longname": {"test-1"}},
		{"longname": {"test0"}, "phone": {"0723158754"}},
		{"longname": {"test1"}, "phone": {"0723158753"}, "email": {"bg1@bg.ro"}},
		{"longname": {"test2"}, "phone": {"0723158752"}, "email": {"bg2@bg.ro"}, "is_male": {"1"}},
		{"longname": {"test3"}, "phone": {"0723158751"}, "email": {"bg3@bg.ro"}, "is_male": {"0"}, "is_client": {"1"}},
		{"longname": {"test4"}, "phone": {"0723158750"}, "email": {"bg4@bg.ro"}, "is_male": {"1"}, "is_client": {"1"}, "is_contractor": {"1"}},
	},
}

func TestAction_AddPerson(t *testing.T) {
	do := openAction(t)
	for n, tcs := range tests {
		t.Run(n, func(t *testing.T) {
			for _, d := range tcs {
				p, ok, verr := printoo.NewPerson(d)
				if !ok {
					t.Fatalf("%v", verr)
				}
				p, err := do.AddPerson(p)
				if err != nil {
					t.Fatalf("insert: %v", err)
				}
				if p.Id == 0 {
					t.Errorf("lid expected %d got %d", 0, p.Id)
				}
				t.Logf("lid %d", p.Id)
				// delete inserted
				t.Run(p.Longname, func(t *testing.T) {
					err = do.DeletePerson(p)
					if err != nil {
						t.Errorf("delete: %v", err)
					}
					t.Logf("did %d", p.Id)
				})
			}
		})
	}
}

func TestAction_EditPerson(t *testing.T) {
	do := openAction(t)
	for n, tcs := range tests {
		t.Run(n, func(t *testing.T) {
			for _, d := range tcs {
				p, ok, verr := printoo.NewPerson(d)
				if !ok {
					t.Fatalf("%v", verr)
				}
				p, err := do.AddPerson(p)
				if err != nil {
					t.Fatalf("insert: %v", err)
				}
				if p.Id == 0 {
					t.Errorf("lid expected %d got %d", 0, p.Id)
				}
				t.Logf("iid %d", p.Id)
				// edit + delete inserted
				t.Run(p.Longname, func(t *testing.T) {
					p, err = do.GetPerson(p.Id)
					if err != nil {
						t.Fatalf("Get %v", err)
					}
					p.Longname += "--edit"
					err = do.EditPerson(p)
					if err != nil {
						t.Errorf("edit: %v", err)
					}
					t.Logf("eid %d", p.Id)
					err = do.DeletePerson(p)
					if err != nil {
						t.Errorf("delete: %v", err)
					}
					t.Logf("did %d", p.Id)
				})
			}
		})
	}
}

func TestAction_GetPerson(t *testing.T) {
	do := openAction(t)
	ids := []string{"1", "2", "3"}
	for _, sid := range ids {
		t.Run(sid, func(t *testing.T) {
			id, err := strconv.Atoi(sid)
			if err != nil {
				t.Fatalf("id %q %v", id, err)
			}
			p, err := do.GetPerson(id)
			if err != nil {
				t.Errorf("Get %v", err)
			}

			if p.IsMale != nil {
				t.Log(*p.IsMale)
			} else {
				t.Log(p.IsMale)
			}
		})
	}
}

func TestAction_AddExtraPhones(t *testing.T) {
	do := openAction(t)
	tests := map[string][]string{
		"1": {"0000000000", "6886856821", "0723158571"},
		"2": {"0000000000", "6886856821", "0723158571"},
	}
	for sid, pps := range tests {
		t.Run(sid, func(t *testing.T) {
			for _, pp := range pps {
				ok := validation.IsPhone(pp, 10)
				if !ok {
					t.Fatalf("%v", pp)
				}
				pid, err := strconv.Atoi(sid)
				if err != nil {
					t.Fatalf("toid: %v", err)
				}
				err = do.AddExtraPhone(pid, pp)
				if err != nil {
					t.Fatalf("add extra: %v", err)
				}
				t.Logf("pp %s", pp)
				// delete inserted
				t.Run(pp, func(t *testing.T) {
					err = do.DeleteExtraPhone(pid, pp)
					if err != nil {
						t.Errorf("delete: %v", err)
					}
					t.Logf("did %s", pp)
				})
			}
		})
	}
}

func openAction(t *testing.T) action.PersonManager {
	db, err := printoo.Open("root:M0b1d1c3@tcp(:3306)/printoo")
	if err != nil {
		t.Fatalf("%v", err)
	}
	do := action.NewHave(db)
	if do == nil {
		t.Fatal("unexpected nil")
	}
	return do
}
