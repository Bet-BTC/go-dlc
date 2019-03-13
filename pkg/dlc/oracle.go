package dlc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/p2pderivatives/dlc/pkg/oracle"
	"github.com/p2pderivatives/dlc/pkg/schnorr"
)

// Oracle contains pubkeys and commitments and signature received from oracle
type Oracle struct {
	PubkeySet   *oracle.PubkeySet  // Oracle's pubkey set
	Commitments []*btcec.PublicKey // Commitments for deals
	Sig         []byte             // Signature for a fixed deal
	SignedMsgs  [][]byte           // Messages signed by Oracle
}

func NewOracle(n int) *Oracle {
	return &Oracle{
		Commitments: make([]*btcec.PublicKey, n)}
}

// PrepareOracleCommitments prepares oracle's commitments for all deals
func (d *DLC) PrepareOracleCommitments(
	V *btcec.PublicKey, Rs []*btcec.PublicKey) {
	for i, deal := range d.Conds.Deals {
		C := schnorr.CommitMulti(V, Rs, deal.Msgs)
		d.Oracle.Commitments[i] = C
	}
}

// SetOraclePubkeySet sets oracle's pubkey set
func (b *Builder) SetOraclePubkeySet(pubset *oracle.PubkeySet) {
	b.dlc.PrepareOracleCommitments(
		pubset.Pubkey, pubset.CommittedRpoints)
	b.dlc.Oracle.PubkeySet = pubset
}

// FixDeal fixes a deal by setting the signature provided by oracle
func (d *DLC) FixDeal(msgs [][]byte, sigs [][]byte) error {
	dID, _, err := d.DealByMsgs(msgs)
	if err != nil {
		return err
	}

	C := d.Oracle.Commitments[dID]
	s := schnorr.SumSigs(sigs)

	ok := schnorr.Verify(C, s)
	if !ok {
		return errors.New("invalid oracle signature")
	}

	// set fixed messages and signature for it
	d.Oracle.SignedMsgs = msgs
	d.Oracle.Sig = s

	return nil
}

// FixDeal fixes a deal by a oracle's signature set by picking up required messages and sigs
func (b *Builder) FixDeal(fm *oracle.SignedMsg, idxs []int) error {
	msgs := [][]byte{}
	sigs := [][]byte{}
	for _, idx := range idxs {
		msgs = append(msgs, fm.Msgs[idx])
		sigs = append(sigs, fm.Sigs[idx])
	}
	return b.dlc.FixDeal(msgs, sigs)
}

// FixedDeal returns a fixed deal
func (d *DLC) FixedDeal() (idx int, deal *Deal, err error) {
	if !d.HasDealFixed() {
		err = newNoFixedDealError()
		return
	}
	return d.DealByMsgs(d.Oracle.SignedMsgs)
}

// HasDealFixed checks if a deal has been fixed
func (d *DLC) HasDealFixed() bool {
	return d.Oracle.SignedMsgs != nil && d.Oracle.Sig != nil
}
