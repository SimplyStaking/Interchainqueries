# Interchainqueries

The purpose of this repository is to demonstrate the capabilities and usage of the `x/icq` module. This module can create cross-chain queries where one would be able to retrieve the balance of an Account on another IBC connected chain or Spot-Price from Osmosis. This module only works with the following [GO Relayer](https://github.com/SimplyVC/relayer).

## DISCLAIMER

This code is unaudited and untested, it is only for demonstrational purposes to show how one can implement an end to end trustless query from one IBC connected blockchain to another. We do not assume any responsibility for any loss of funds due to the use of this code.

## Notes

ICQ can perform periodic queries from a source chain to a target chain. The query will only read a KV store of the target chain and not the normal GRPC endpoint as the GRPC endpoint cannot provide Proofs to validate our results. Read more about the module in the [DESIGN DOCS](./docs/design/ICQ.md).

## Get started

```
starport chain serve -c config-icq-1.yml
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.


## Example Periodic Query

An example periodic query which will query the balance of an account on the target chain.

```
  path := "/store/bank/key"
  clientId := "07-tendermint-0"
  chainId := "test-2"
  queryModuleAddress := "test_module_address"

  add1, _ := sdk.AccAddressFromBech32("icq17dtl0mjt3t77kpuhg2edqzjpszulwhgzp4faq3")
  icqPeriodic := types.PeriodicICQs{
    Path:                 path,
    TimeoutHeightPadding: uint64(20),
    TargetHeight:         uint64(0),
    ClientId:             clientId,
    Creator:              queryModuleAddress,
    ChainId:              chainId,
    QueryParameters:      append(banktypes.CreateAccountBalancesPrefix(add1), []byte("stake")...),
    BlockRepeat:          uint64(10),
    LastHeightExecuted:   uint64(ctx.BlockHeight()),
    MaxResults:         uint64(10),
  }
  k.AppendPeriodicICQs(ctx, icqPeriodic)
```

Once the result is received one can convert the raw data to the exact type as such:

```
var balance sdk.Coin
cdc.MustUnmarshal(data, &balance)
```

Which will be converted to a balance of _X_ stake tokens.


## Research Resources

- [IBC Queries Discussion](https://github.com/cosmos/ibc/discussions/605)
- [Cross-chain query spec draft](https://github.com/cosmos/ibc/pull/735)
- [Defund query module](https://github.com/defund-labs/defund/tree/main/x/query)
- [Defund relayer](https://github.com/defund-labs/relayer)
- [Quicksilver relayer](https://github.com/ingenuity-build/interchain-queries)
- [Quicksilver interchainquery module](https://github.com/ingenuity-build/quicksilver/tree/main/x/interchainquery)