package alignmentprofile

import (
	a "github.com/hivdb/nucamino/types/amino"
	"reflect"
	"testing"
)

var exampleProfile = AlignmentProfile{
	StopCodonPenalty:         1,
	GapOpeningPenalty:        2,
	GapExtensionPenalty:      3,
	IndelCodonOpeningBonus:   4,
	IndelCodonExtensionBonus: 5,
	ReferenceSequences: ReferenceSeqs{
		"A": a.ReadString("TTALIEPPVYPIVEHSDEKTAHEEH"),
		"B": a.ReadString("CSNELVISHEADPVWRSAVLRGAP"),
	},
	GeneIndelScores: GenePositionalIndelScores{
		Gene("A"): PositionalIndelScores{
			3:  [2]int{4, 5},
			6:  [2]int{7, 8},
			-6: [2]int{7, 8},
			-9: [2]int{10, 11},
		},
		Gene("B"): PositionalIndelScores{
			2: [2]int{1, 2},
		},
	},
}

func TestGenes(t *testing.T) {
	genes := exampleProfile.Genes()
	expected := [2]Gene{"A", "B"}

	if len(genes) != len(expected) {
		t.Errorf("Expected %v to be %v", genes, expected)
	}

	for _, exp := range expected {
		found := false
		for _, g := range genes {
			if g == exp {
				found = true
			}
		}
		if !found {
			t.Errorf("Missing key: %v", exp)
		}
	}
}

func TestRawIndelScoresFromProfile(t *testing.T) {
	expected := map[string][]rawIndelScore{
		"A": []rawIndelScore{
			{Kind: "ins", Position: 3, Open: 4, Extend: 5},
			{Kind: "ins", Position: 6, Open: 7, Extend: 8},
			{Kind: "del", Position: 6, Open: 7, Extend: 8},
			{Kind: "del", Position: 9, Open: 10, Extend: 11},
		},
		"B": []rawIndelScore{
			{Kind: "ins", Position: 2, Open: 1, Extend: 2},
		},
	}
	constructed := exampleProfile.rawIndelScores()
	if !reflect.DeepEqual(constructed, expected) {
		t.Errorf("%v != %v", constructed, expected)
	}
}

func TestRoundTripToRawProfile(t *testing.T) {
	raw := exampleProfile.asRaw()
	constructed, err := raw.asProfile()
	if err != nil {
		t.Errorf("Unexpected error while converting profile to raw profile: %v", err)
		t.FailNow()
	}
	if !reflect.DeepEqual(*constructed, exampleProfile) {
		t.Errorf("%v != %v", *constructed, exampleProfile)
	}
}

func TestRetrievingPositionalIndelScores(t *testing.T) {
	_, found := exampleProfile.PositionalIndelScoresFor("A")
	if !found {
		t.Errorf("Failed to retrieve positional indel scores in example profile")
	}
	_, found = exampleProfile.PositionalIndelScoresFor("Z")
	if found {
		t.Errorf("Found positional indel scores for non-existent gene in example profile")
	}
}
