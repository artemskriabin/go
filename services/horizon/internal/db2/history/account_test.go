package history

import (
	"testing"

	"github.com/stellar/go/services/horizon/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestAccountQueries(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	q := &Q{tt.HorizonSession()}

	// Test Accounts()
	acs := []Account{}
	err := q.Accounts().Select(&acs)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(acs, 4)
	}
}

func TestIsAuthRequired(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 1}
	tt.True(account.IsAuthRequired())

	account = AccountEntry{Flags: 0}
	tt.False(account.IsAuthRequired())
}

func TestIsAuthRevocable(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 2}
	tt.True(account.IsAuthRevocable())

	account = AccountEntry{Flags: 1}
	tt.False(account.IsAuthRevocable())
}
func TestIsAuthImmutable(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 4}
	tt.True(account.IsAuthImmutable())

	account = AccountEntry{Flags: 0}
	tt.False(account.IsAuthImmutable())
}

func assertAccountsContainAddresses(tt *test.T, accounts map[string]int64, addresses []string) {
	tt.Assert.Len(accounts, len(addresses))
	set := map[int64]bool{}
	for _, address := range addresses {
		accountID, ok := accounts[address]
		tt.Assert.True(ok)
		tt.Assert.False(set[accountID])
		set[accountID] = true
	}
}

func TestCreateAccounts(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetHorizonDB(t, tt.HorizonDB)
	q := &Q{tt.HorizonSession()}

	addresses := []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	}
	accounts, err := q.CreateAccounts(addresses, 1)
	tt.Assert.NoError(err)
	tt.Assert.Len(accounts, 2)
	assertAccountsContainAddresses(tt, accounts, addresses)

	dupAccounts, err := q.CreateAccounts([]string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	}, 2)
	tt.Assert.NoError(err)
	tt.Assert.Equal(accounts, dupAccounts)

	addresses = []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		"GCYVFGI3SEQJGBNQQG7YCMFWEYOHK3XPVOVPA6C566PXWN4SN7LILZSM",
		"GBYSBDAJZMHL5AMD7QXQ3JEP3Q4GLKADWIJURAAHQALNAWD6Z5XF2RAC",
	}
	accounts, err = q.CreateAccounts(addresses, 1)
	tt.Assert.NoError(err)
	assertAccountsContainAddresses(tt, accounts, addresses)
	for address, accountID := range dupAccounts {
		id, ok := accounts[address]
		tt.Assert.True(ok)
		tt.Assert.Equal(id, accountID)
	}
}
