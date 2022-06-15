package keeper

import (
	"github.com/irclausen/nameservice/x/nameservice/types"
)

var _ types.QueryServer = Keeper{}
