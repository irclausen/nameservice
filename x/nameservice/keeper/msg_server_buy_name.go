package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irclausen/nameservice/x/nameservice/types"
)

func (k msgServer) BuyName(goCtx context.Context, msg *types.MsgBuyName) (*types.MsgBuyNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Try getting a name from the store
	whois, isFound := k.GetWhois(ctx, msg.Name)
	// Set the price at which the name has to be bought if it didn't have and owner before
	minPrice := sdk.Coins{sdk.NewInt64Coin("token", 10)}
	// Convert price and bid strings into sdk.Coins
	price, _ := sdk.ParseCoinsNormalized(whois.Price)
	bid, _ := sdk.ParseCoinsNormalized(msg.Bid)
	// Convert the owner and buyer address strings to sdk.AccAddress
	owner, _ := sdk.AccAddressFromBech32(whois.Owner)
	buyer, _ := sdk.AccAddressFromBech32(msg.Creator)
	// If a nam is found in store
	if isFound {
		// If the current price is higher than the bid
		if price.IsAllGT(bid) {
			// Throw an error
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is not high enough")
		}
		// Otherwise (when the bid is higher), send tokens from the buyer to the owner
		k.bankKeeper.SendCoins(ctx, buyer, owner, bid)
	} else {
		// If the nams is not found in the store
		// If the minimum price is higher than the bid
		if minPrice.IsAllGT(bid) {
			// Throw an error
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is less than min amound")
		}
		// Otherwise (when the bid is higher), send tokesn from the buyesr's account to the module's account (as payment for the name)
		k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, bid)
	}
	// Create an updated whois record
	newWhois := types.Whois{
		Index: msg.Name,
		Name:  msg.Name,
		Value: whois.Value,
		Price: bid.String(),
		Owner: buyer.String(),
	}
	// Write whois information to the store
	k.SetWhois(ctx, newWhois)
	return &types.MsgBuyNameResponse{}, nil
}
