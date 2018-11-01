package script

import (
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dgarage/dlc/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestP2WPKHpkScript(t *testing.T) {
	assert := assert.New(t)

	priv, pub := test.RandKeys()
	amt := int64(10000)

	// create P2WPKHpkScript
	pkScript, err := P2WPKHpkScript(pub)
	assert.Nil(err)

	// prepare source tx
	sourceTx := test.SourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(amt, pkScript))

	// create redeem tx
	redeemTx := createRedeemTx(sourceTx)

	// witness signature
	sign, err := WitnessSignature(redeemTx, 0, amt, pkScript, priv)
	assert.Nil(err)

	// redeem script
	wt := wire.TxWitness{sign, pub.SerializeCompressed()}
	redeemTx.TxIn[0].Witness = wt

	// execute script
	err = executeScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}

func TestMultiSigScript2of2(t *testing.T) {
	assert := assert.New(t)

	priv1, pub1 := test.RandKeys()
	priv2, pub2 := test.RandKeys()
	amt := int64(10000)

	script, err := MultiSigScript2of2(pub1, pub2)
	assert.Nil(err)
	pkScript, err := P2WSHpkScript(script)
	assert.Nil(err)

	// prepare source tx
	sourceTx := test.SourceTx()
	sourceTx.AddTxOut(wire.NewTxOut(amt, pkScript))

	// create redeem tx
	redeemTx := createRedeemTx(sourceTx)

	// witness signatures
	sign1, err := WitnessSignature(redeemTx, 0, amt, script, priv1)
	assert.Nil(err)
	sign2, err := WitnessSignature(redeemTx, 0, amt, script, priv2)
	assert.Nil(err)

	// redeem script
	wt := wire.TxWitness{[]byte{}, sign1, sign2, script}
	redeemTx.TxIn[0].Witness = wt

	// execute script
	err = executeScript(pkScript, redeemTx, amt)
	assert.Nil(err)
}

func createRedeemTx(sourceTx *wire.MsgTx) *wire.MsgTx {
	txHash := sourceTx.TxHash()
	outPt := wire.NewOutPoint(&txHash, 0)
	tx := wire.NewMsgTx(test.TxVersion)
	tx.AddTxIn(wire.NewTxIn(outPt, nil, nil))
	return tx
}

func executeScript(pkScript []byte, tx *wire.MsgTx, amt int64) error {
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(pkScript, tx, 0, flags, nil, nil, amt)
	if err != nil {
		return err
	}
	return vm.Execute()
}