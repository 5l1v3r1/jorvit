package datastore

import (
	"fmt"
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

	payloads := make(map[string]int)

	for _, v := range *b.List {
		if v.VoteAction == "" {
			v.VoteAction = "off_chain"
		} else if v.VoteAction != "off_chain" {
			return fmt.Errorf("%s - expected to be one of (%s) - but [%s] provided", "chain_vote_action", "off_chain", v.VoteAction)
		}

		if v.VoteType == "" {
			v.VoteType = "public"
		} else if v.VoteType != "public" && v.VoteType != "private" {
			return fmt.Errorf("%s - expected to be one of (%s, %s) - but [%s] provided", "chain_vote_type", "public", "private", v.VoteType)
		}
		v.ChainVotePlan.Payload = v.VoteType
		payloads[v.VoteType] += 1
		if len(payloads) > 1 {
			return fmt.Errorf("%s - expected to be only one of (%s, %s) for the whole list, but [%s] provided. Please split your list so it contains only one (%s)", "chain_vote_type", "public", "private", "multiple", "chain_vote_type")
		}
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

func (b *Proposals) Payloads() map[string]int {
	payloads := make(map[string]int)

	for _, v := range *b.List {
		if v.VoteType == "" {
			v.VoteType = "public"
		}
		payloads[v.VoteType] += 1
	}

	return payloads
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

type Funds struct {
	List *[]*loader.FundData `json:"funds"`
}

func (b *Funds) Initialize(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	b.List, err = loader.LoadFundData(file)
	if err != nil {
		return err
	}

	for _, v := range *b.List {
		v.Voteplans = make([]loader.ChainVotePlan, 0)
	}

	return nil
}

// tmp
func (b *Funds) First() *loader.FundData {
	if len(*b.List) == 0 {
		return nil
	}
	return (*b.List)[0]
}

func (b *Funds) Total() int {
	return len(*b.List)
}
