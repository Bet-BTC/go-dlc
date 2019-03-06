package dlcmgr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb" // blank import for bolt db driver
	"github.com/dgarage/dlc/pkg/dlc"
	"github.com/stretchr/testify/assert"
)

var (
	testDBPrefix = "dlcmgr_"
	testDBName   = "dlcmgr.db"
)

func TestCreateAndOpen(t *testing.T) {
	assert := assert.New(t)

	db, closeFunc := newWalletDB()
	defer closeFunc()

	manager, err := Create(db)
	assert.NoError(err)
	assert.NotNil(manager)

	manager, err = Open(db)
	assert.NoError(err)
	assert.NotNil(manager)
}

func TestStoreContract(t *testing.T) {
	assert := assert.New(t)

	// create new db
	db, closeFunc := newWalletDB()
	defer closeFunc()
	manager, _ := Create(db)

	key := []byte("testdlc")
	dOrig := newDLC()
	err := manager.StoreContract(key, dOrig)
	assert.NoError(err)

	d, err := manager.RetrieveContract(key)
	assert.NoError(err)
	assert.NotNil(d)

	// TODO: serialize all information
	assert.Equal(d.Conds, dOrig.Conds)
}

func newWalletDB() (walletdb.DB, func()) {
	path := testDBPath()
	db, _ := walletdb.Create("bdb", path)
	closeFunc := func() {
		db.Close()
		os.RemoveAll(path)
	}
	return db, closeFunc
}

func testDBPath() string {
	dir, _ := ioutil.TempDir("", testDBPrefix)
	path := filepath.Join(dir, testDBName)
	return path
}

func newDLC() *dlc.DLC {
	ftime := fixingTime()
	famt1, _ := btcutil.NewAmount(1)
	famt2, _ := btcutil.NewAmount(1)
	feerate := btcutil.Amount(10)
	refundlc := uint32(1)
	deals := testDeals()

	conds, _ := dlc.NewConditions(
		ftime, famt1, famt2, feerate, feerate, refundlc, deals)

	return &dlc.DLC{Conds: conds}
}

func testDeals() []*dlc.Deal {
	deals := []*dlc.Deal{}

	total := 5
	nMsg := 3
	for i := 0; i < total+1; i++ {
		amt1 := btcutil.Amount(i)
		amt2 := btcutil.Amount(total - i)
		msgs := [][]byte{}
		for j := 0; j < nMsg; j++ {
			msgs = append(msgs, []byte{byte(i)})
		}
		deal := dlc.NewDeal(amt1, amt2, msgs)
		deals = append(deals, deal)
	}

	return deals
}

func fixingTime() time.Time {
	t := time.Now().AddDate(0, 0, 1)
	y, m, d := t.Date()
	return time.Date(y, m, d, 12, 0, 0, 0, time.UTC)
}