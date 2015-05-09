package imap

import (
	"testing"
)

func testRange(t *testing.T, rangeStr string, expectedMin SequenceNumber, expectedMax SequenceNumber, expectedErr error) {
	rng, err := interpretMessageRange(rangeStr)
	if rng.min != expectedMin {
		t.Errorf("Range '%s': min '%s' did not match expected '%s'", rangeStr, rng.min, expectedMin)
	}
	if rng.max != expectedMax {
		t.Errorf("Range '%s': max '%s' did not match expected '%s'", rangeStr, rng.max, expectedMax)
	}
	assertErr(t, expectedErr, err)
}

func assertErr(t *testing.T, expectedErr error, actualErr error) {
	if expectedErr == nil && actualErr != nil {
		t.Errorf("Expected nil error but got %s", actualErr.Error())
	} else if expectedErr != nil && actualErr == nil {
		t.Errorf("Expected error %s but got nil error", expectedErr.Error())
	} else if expectedErr != actualErr {
		t.Errorf("Expected error %s but got error %s", expectedErr.Error(), actualErr.Error())
	}
}

func testSet(t *testing.T, setStr string, expectedSet []SequenceRange, expectedErr error) {
	set, err := interpretSequenceSet(setStr)
	assertErr(t, expectedErr, err)

	if len(set) != len(expectedSet) {
		t.Errorf("Sequence set %s\n"+
			"\tLength %d does not match expected %d", setStr, len(set), len(expectedSet))
	}
}

func TestFindMessageRange(t *testing.T) {
	testRange(t, "15:95", "15", "95", nil)
	testRange(t, "53:*", "53", "*", nil)
	testRange(t, "35", "35", "", nil)
	testRange(t, "*", "*", "", nil)
	testRange(t, "5*", "", "", errInvalidRangeString("5*"))
	testRange(t, "*5*", "", "", errInvalidRangeString("*5*"))
	testRange(t, "hello", "", "", errInvalidRangeString("hello"))
}

func TestSequenceSet(t *testing.T) {
	testSet(t, "118:*", []SequenceRange{
		SequenceRange{min: "118", max: "*"},
	}, nil)
	testSet(t, "1,3,4:14", []SequenceRange{
		SequenceRange{min: "1", max: ""},
		SequenceRange{min: "3", max: ""},
		SequenceRange{min: "4", max: "14"},
	}, nil)
	testSet(t, "1,3,8:14,18:*", []SequenceRange{
		SequenceRange{min: "1", max: ""},
		SequenceRange{min: "3", max: ""},
		SequenceRange{min: "8", max: "14"},
		SequenceRange{min: "18", max: "8"},
	}, nil)
	testSet(t, "1,3,:8:14,18:*", nil, errInvalidSequenceSetString("1,3,:8:14,18:*"))
}

func TestSequenceNumber(t *testing.T) {
	a := SequenceNumber("*")
	if a.Last() == false {
		t.Errorf("Last() function for sequence number of * should return true")
	}
	if a.Nil() == true {
		t.Errorf("Nil() function for non-blank sequence number should return false")
	}
	if _, err := a.Value(); err == nil {
		t.Errorf("Value() function for sequence number of * should return an error")
	}

	b := SequenceNumber("56")
	if b.Last() == true {
		t.Errorf("Last() function for sequence number of 56 should return false")
	}
	if v, _ := b.Value(); v != 56 {
		t.Errorf("Value() function for sequence number of 56 should return integer 56")
	}
	if _, err := b.Value(); err != nil {
		t.Errorf("Value() function for sequence number of 56 should not return an error")
	}

	c := SequenceNumber("")
	if c.Last() == true {
		t.Errorf("Last() function for blank sequence number should return false")
	}
	if c.Nil() == false {
		t.Errorf("Nil() function for blank sequence number should return true")
	}
	if _, err := c.Value(); err == nil {
		t.Errorf("Value() function for blank sequence number should return an error")
	}
}
