package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/dgarage/dlc/internal/script"
)

func (w *wallet) NewPubkey() (pub *btcec.PublicKey, err error) {
	// TODO: generate new pubkey and address using newAddress
	_, err = w.newAddress(waddrmgr.KeyScopeBIP0084, []byte{}, uint32(1), uint32(1))
	if err != nil {
		return
	}
	return
}

func (w *wallet) NewWitnessPubkeyScript() (pkScript []byte, err error) {
	var pub *btcec.PublicKey
	pub, err = w.NewPubkey()
	if err != nil {
		return
	}
	return script.P2WPKHpkScript(pub)
}

// NewAddress returns a new ManagedAddress for a given scope and account number.
// NOTE: this function callsNextExternalAddresses to generate a ManagadAdddress.
func (w *wallet) newAddress(scope waddrmgr.KeyScope, privPass []byte,
	account uint32, numAddresses uint32) ([]waddrmgr.ManagedAddress, error) {
	// unlock Manager
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		e := w.manager.Unlock(ns, privPass)
		return e
	})
	if err != nil {
		return nil, err
	}

	// get ScopedKeyManager
	scopedMgr, err := w.manager.FetchScopedKeyManager(scope)
	if err != nil {
		return nil, err
	}

	var addrs []waddrmgr.ManagedAddress
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var e error
		addrs, e = scopedMgr.NextExternalAddresses(ns, account, numAddresses)
		return e
	})
	if err != nil {
		return nil, err
	}

	return addrs, nil
}
