// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import btcec "github.com/btcsuite/btcd/btcec"
import btcjson "github.com/btcsuite/btcd/btcjson"
import btcutil "github.com/btcsuite/btcutil"
import mock "github.com/stretchr/testify/mock"
import wallet "github.com/dgarage/dlc/internal/wallet"
import wire "github.com/btcsuite/btcd/wire"

// Wallet is an autogenerated mock type for the Wallet type
type Wallet struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Wallet) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListUnspent provides a mock function with given fields:
func (_m *Wallet) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	ret := _m.Called()

	var r0 []btcjson.ListUnspentResult
	if rf, ok := ret.Get(0).(func() []btcjson.ListUnspentResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPubkey provides a mock function with given fields:
func (_m *Wallet) NewPubkey() (*btcec.PublicKey, error) {
	ret := _m.Called()

	var r0 *btcec.PublicKey
	if rf, ok := ret.Get(0).(func() *btcec.PublicKey); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcec.PublicKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectUnspent provides a mock function with given fields: amt, feePerTxIn, feePerTxOut
func (_m *Wallet) SelectUnspent(amt btcutil.Amount, feePerTxIn btcutil.Amount, feePerTxOut btcutil.Amount) ([]btcjson.ListUnspentResult, btcutil.Amount, error) {
	ret := _m.Called(amt, feePerTxIn, feePerTxOut)

	var r0 []btcjson.ListUnspentResult
	if rf, ok := ret.Get(0).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) []btcjson.ListUnspentResult); ok {
		r0 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	var r1 btcutil.Amount
	if rf, ok := ret.Get(1).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) btcutil.Amount); ok {
		r1 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		r1 = ret.Get(1).(btcutil.Amount)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(btcutil.Amount, btcutil.Amount, btcutil.Amount) error); ok {
		r2 = rf(amt, feePerTxIn, feePerTxOut)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Unlock provides a mock function with given fields: privPass
func (_m *Wallet) Unlock(privPass []byte) error {
	ret := _m.Called(privPass)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(privPass)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WitnessSignature provides a mock function with given fields: tx, idx, amt, sc, pub
func (_m *Wallet) WitnessSignature(tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey) ([]byte, error) {
	ret := _m.Called(tx, idx, amt, sc, pub)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey) []byte); ok {
		r0 = rf(tx, idx, amt, sc, pub)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey) error); ok {
		r1 = rf(tx, idx, amt, sc, pub)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WitnessSignatureWithCallback provides a mock function with given fields: tx, idx, amt, sc, pub, privkeyConverter
func (_m *Wallet) WitnessSignatureWithCallback(tx *wire.MsgTx, idx int, amt btcutil.Amount, sc []byte, pub *btcec.PublicKey, privkeyConverter wallet.PrivateKeyConverter) ([]byte, error) {
	ret := _m.Called(tx, idx, amt, sc, pub, privkeyConverter)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey, wallet.PrivateKeyConverter) []byte); ok {
		r0 = rf(tx, idx, amt, sc, pub, privkeyConverter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*wire.MsgTx, int, btcutil.Amount, []byte, *btcec.PublicKey, wallet.PrivateKeyConverter) error); ok {
		r1 = rf(tx, idx, amt, sc, pub, privkeyConverter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
