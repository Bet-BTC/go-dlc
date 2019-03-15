package dlc

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/p2pderivatives/dlc/internal/oracle"
	"github.com/p2pderivatives/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestContractExecutionTx(t *testing.T) {
	assert := assert.New(t)

	// A deal that has both amounts are > 0
	var damt1, damt2 btcutil.Amount = 1, 1
	b, _, dID, deal, err := setupContractorsUntilPubkeyExchange(damt1, damt2)
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	// fail without oracle's message commitment
	_, err = b.Contract.ContractExecutionTx(b.party, deal, dID)
	assert.NotNil(err)

	// set message commitment
	_, C := test.RandKeys()
	b.Contract.Oracle.Commitments[dID] = C

	// txout should have 2 entries
	tx, err := b.Contract.ContractExecutionTx(b.party, deal, dID)
	assert.Nil(err)
	assert.Len(tx.TxOut, 2)
	assert.Equal(int64(damt1), tx.TxOut[0].Value)
	assert.Equal(int64(damt2), tx.TxOut[1].Value)
}

// An edge case that a executing party tx takes all funds
func TestContractExecutionTxTakeAll(t *testing.T) {
	assert := assert.New(t)

	var damt1, damt2 btcutil.Amount = 1, 0
	b, _, dID, deal, err := setupContractorsUntilPubkeyExchange(damt1, damt2)
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	_, C := test.RandKeys()
	b.Contract.Oracle.Commitments[dID] = C

	tx, err := b.Contract.ContractExecutionTx(b.party, deal, dID)

	assert.NoError(err)
	assert.Len(tx.TxOut, 1)
	assert.Equal(int64(damt1), tx.TxOut[0].Value)
}

// An edge case that a executing party tx takes nothing
func TestContractExecutionTxTakeNothing(t *testing.T) {
	assert := assert.New(t)

	var damt1, damt2 btcutil.Amount = 0, 1
	b, _, dID, deal, err := setupContractorsUntilPubkeyExchange(damt1, damt2)
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	_, C := test.RandKeys()
	b.Contract.Oracle.Commitments[dID] = C

	tx, err := b.Contract.ContractExecutionTx(b.party, deal, dID)

	// asserions
	assert.Nil(tx)
	assert.NotNil(err)
	assert.IsType(&CETTakeNothingError{}, err)
}

func TestSignedContractExecutionTx(t *testing.T) {
	assert := assert.New(t)
	var err error

	// setup
	b1, b2, dID, deal, err := setupContractorsUntilPubkeyExchange(1, 1)
	if !assert.NoError(err) {
		assert.FailNow(err.Error())
	}

	privkey, C := test.RandKeys()
	b1.Contract.Oracle.Commitments[dID] = C
	b2.Contract.Oracle.Commitments[dID] = C
	osigs := [][]byte{privkey.D.Bytes()}
	oFixedMsg := &oracle.SignedMsg{Msgs: deal.Msgs, Sigs: osigs}

	err = b1.FixDeal(oFixedMsg, []int{0})
	assert.NoError(err)
	err = b2.FixDeal(oFixedMsg, []int{0})
	assert.NoError(err)

	// fail without the counterparty's signatures
	_, err = b1.SignedContractExecutionTx()
	assert.NoError(err)
	_, err = b2.SignedContractExecutionTx()
	assert.NoError(err)

	// exchange signs
	sig1, err := b1.SignContractExecutionTx(deal, dID)
	assert.NoError(err)
	sig2, err := b2.SignContractExecutionTx(deal, dID)
	assert.Nil(err)

	err = b1.AcceptCETxSignatures([][]byte{sig2})
	assert.Nil(err)
	err = b2.AcceptCETxSignatures([][]byte{sig1})
	assert.Nil(err)

	// no errors with the counterparty's sign
	tx1, err := b1.SignedContractExecutionTx()
	assert.NoError(err)
	tx2, err := b2.SignedContractExecutionTx()
	assert.NoError(err)

	// each party has a tx that has the same txin but has different txouts
	assert.Len(tx1.TxOut, 2)
	assert.Len(tx2.TxOut, 2)
	assert.Equal(
		tx1.TxIn[fundTxInAt].PreviousOutPoint,
		tx2.TxIn[fundTxInAt].PreviousOutPoint)
	assert.Equal(tx1.TxOut[0].Value, tx2.TxOut[1].Value)
	assert.Equal(tx1.TxOut[1].Value, tx2.TxOut[0].Value)

	// Both parties are able to send the CET
	err = runFundScript(b1, tx1)
	assert.Nil(err)
	err = runFundScript(b2, tx2)
	assert.Nil(err)
}

func setupContractorsUntilPubkeyExchange(
	damt1, damt2 btcutil.Amount,
) (b1, b2 *Builder, dID int, deal *Deal, err error) {
	msgs := [][]byte{{1}}
	deal = NewDeal(damt1, damt2, [][]byte{{1}})

	setupConds := func() *Conditions {
		conds := newTestConditions()
		conds.Deals = []*Deal{deal}
		return conds
	}

	// init first party
	b1 = setupBuilder(FirstParty, setupTestWallet, setupConds)
	if err = stepPrepare(b1); err != nil {
		return
	}

	// init second party
	b2 = setupBuilder(SecondParty, setupTestWallet, setupConds)
	if err = stepPrepare(b2); err != nil {
		return
	}

	// exchange pubkeys and utxos
	if err = stepSendRequirments(b1, b2); err != nil {
		return
	}
	if err = stepSendRequirments(b2, b1); err != nil {
		return
	}

	dID, deal, _ = b1.Contract.DealByMsgs(msgs)

	return b1, b2, dID, deal, nil
}

func runFundScript(b *Builder, tx *wire.MsgTx) error {
	fundtx, _ := b.Contract.FundTx()
	fout := fundtx.TxOut[fundTxOutAt]
	return test.ExecuteScript(fout.PkScript, tx, fout.Value)
}
