package integration

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/dgarage/dlc/internal/dlc"
	"github.com/dgarage/dlc/internal/oracle"
	"github.com/stretchr/testify/assert"
)

func TestOracleCommitAndSign(t *testing.T) {
	// Given an oracle "Olivia" who provides weather information
	//   weather info contains "weather" "temperature" "windspeed"
	olivia, _ := newOracle("Olivia", 3)

	// And a contractor "Alice"
	alice, _ := newContractor("Alice")

	// And Alice bet on "weather" and "temprature" at a future time
	fixingTime := contractorBetOnWeatherAndTemperature(t, alice)

	// Alice asks Olivia to fix weather info at the fixing time
	contractorAsksOracleToCommit(t, alice, olivia)

	// Olivia fixes weather info
	fixedWeather := oracleFixesWeather(t, olivia, fixingTime)

	// When Alice fixes a deal using Olivia's sign and singed weather info
	contractorFixesDeal(t, alice, olivia)

	// Then The fixed deal should be a subset of the fixed weather info
	shouldFixedDealSameWithFixedWeather(t, alice, fixedWeather)
}

func newWeather(weather string, temp int, windSpeed int) [][]byte {
	return [][]byte{
		[]byte(weather),
		[]byte(strconv.Itoa(temp)),
		[]byte(strconv.Itoa(windSpeed)),
	}
}

func contractorBetOnWeatherAndTemperature(t *testing.T, c *Contractor) time.Time {
	deal1 := dlc.NewDeal(2, 0, newWeather("fine", 20, 0)[:2])
	deal2 := dlc.NewDeal(1, 1, newWeather("fine", 10, 0)[:2])
	deal3 := dlc.NewDeal(1, 1, newWeather("rain", 20, 0)[:2])
	deal4 := dlc.NewDeal(0, 2, newWeather("rain", 10, 0)[:2])
	deals := []*dlc.Deal{deal1, deal2, deal3, deal4}
	fixingTime := time.Now().AddDate(0, 0, 1)
	conds, err := dlc.NewConditions(fixingTime, 1, 1, 1, 1, 1, deals)
	assert.NoError(t, err)
	c.createDLCBuilder(conds, dlc.FirstParty)
	return fixingTime
}

func contractorAsksOracleToCommit(
	t *testing.T, c *Contractor, o *oracle.Oracle) {
	ftime := c.DLCBuilder.DLC().Conds.FixingTime

	pubkeySet, err := o.PubkeySet(ftime)
	assert.NoError(t, err)

	c.DLCBuilder.SetOraclePubkeySet(&pubkeySet)
}

func oracleFixesWeather(
	t *testing.T, o *oracle.Oracle, ftime time.Time) [][]byte {
	msgs := [][][]byte{
		newWeather("fine", 20, 0),
		newWeather("fine", 20, 5),
		newWeather("fine", 10, 0),
		newWeather("fine", 10, 5),
		newWeather("rain", 20, 0),
		newWeather("rain", 20, 5),
		newWeather("rain", 10, 0),
		newWeather("rain", 10, 5),
	}

	fixingMsg := randomMsg(msgs)

	err := o.FixMsgs(ftime, fixingMsg)
	assert.NoError(t, err)

	return fixingMsg
}

func randomMsg(msgs [][][]byte) [][]byte {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	idx := r.Intn(len(msgs))
	return msgs[idx]
}

func contractorFixesDeal(t *testing.T, c *Contractor, o *oracle.Oracle) {
	ftime := c.DLCBuilder.DLC().Conds.FixingTime

	// receive signset
	signSet, err := o.SignSet(ftime)
	assert.NoError(t, err)

	// fix deal with the signset
	idxs := []int{0, 1} // use only weather and temperature
	err = c.DLCBuilder.FixDeal(&signSet, idxs)
	assert.NoError(t, err)
}

func shouldFixedDealSameWithFixedWeather(t *testing.T, c *Contractor, fixedWeather [][]byte) {
	_, fixedDeal, err := c.DLCBuilder.DLC().FixedDeal()
	assert.NoError(t, err)
	assert.Equal(t, fixedWeather[:2], fixedDeal.Msgs)
}
