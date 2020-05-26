package datastore

import (
	"os"

	"github.com/input-output-hk/jorvit/internal/loader"
)

type Proposals struct {
	List *[]*loader.ProposalData `json:"proposals"`
}

func (b *Proposals) Initialize(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	b.List, err = loader.LoadData(file)
	if err != nil {
		return err
	}
	for _, v := range *b.List {
		v.ChainProposal.VoteOptions.Blank = 0
		v.ChainProposal.VoteOptions.Yes = 1
		v.ChainProposal.VoteOptions.No = 2
	}
	return nil
}

func (b *Proposals) All() *[]*loader.ProposalData {
	return b.List
}

func (b *Proposals) SearchID(internalID string) *loader.ProposalData {
	ret := FilterSingle(b.List, func(v *loader.ProposalData) bool {
		return v.InternalID == internalID
	})
	return ret
}

func (b *Proposals) Total() int {
	return len(*b.List)
}

func Filter(vs *[]*loader.ProposalData, f func(*loader.ProposalData) bool) *[]*loader.ProposalData {
	ret := make([]*loader.ProposalData, 0)
	for _, v := range *vs {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return &ret
}

func FilterSingle(vs *[]*loader.ProposalData, f func(*loader.ProposalData) bool) *loader.ProposalData {
	for _, v := range *vs {
		if f(v) {
			return v
		}
	}
	return nil
}
