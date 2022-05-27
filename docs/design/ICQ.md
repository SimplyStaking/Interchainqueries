# ICQ Module

The Interchain Queries (ICQ) module leverages the existing relayer infrastructure, with some custom changes, to enable two connected chains to query data from eachothers' state. The module only needs to be imported by the chain that wishes to perform queries. The relayer is expected to parse queries from the requesting chain, perform them on the chain indicated in the query, and pass the result back to the requesting chain. This document will describe the Cosmos SDK module aspect of Interchain Queries. For the relayer aspect, refer to the [Relayer document](./Relayer.md).


- [ICQ Module](#icq-module)
  - [Concepts](#concepts)
  - [State](#state)
  - [Transactions](#transactions)
  - [Params](#params)
  - [BeginBlock](#beginblock)
  - [EndBlock](#endblock)
  - [Dependencies](#dependencies)
  - [Hooks](#hooks)
  - [Future](#future)

## Concepts

- `PeriodicICQ`: describes a periodic data request and serves as a template for instances of such periodic queries. The _periodic_ includes all details about what data is required and which chain it should be queried from. It also specifies how often the query should take place (i.e. the period length) and timeout parameters.
- `PendingICQInstance`: an instance of a periodic query that can be seen as a pending query. An instance is created at every period, as specified by the _periodic_, and this instance is deleted once a result is submitted or it timesout. These query instances are picked up by the custom relayer and fulfilled according to the details in the request.
- `PendingICQRequest`: This is the message that will be retrieved by the relayer. The data is a combination of the `PendingICQInstance` and `PeriodicICQ` to have the full message available to the relayer without storing duplicate data in the blockchain.
- `ICQResult`: stores the `last_result_id` of the last result submitted by the relayer in response to a specific periodic query. There is one `ICQResult` per `PeriodicICQ` and it stores the `periodic_id`. The `last_result_id` points to `DataPointResult` which stores the last result.
- `DataPointResult`: Stores a result of a serviced `PendingICQInstance`.
- `ICQTimeouts`: this is created per `PeriodicICQ` we increment the amount of timeouts that occur for that periodic query as well as the last source chain height that the timeout occured at. If the [BeginBlock](#beginblock) notices that the timeout height was reached it will expire a `PendingICQInstance`s. If the relayer submits a late result, this is completely ignored.

## State

List of `PeriodicICQ`, each with:

- **Id** (`uint64`): basic unique identifier
- **Path** (`string`): query path, e.g. `/store/bank/key`
- **TimeoutHeightPadding** (`uint64`): timeout height will be current height plus this value
- **TargetHeight** (`uint64`):The height of the chain we are querying, this is usually set to `0` to always query the heighest height of the target chain.
- **ClientId** (`string`): IBC client Id of the target chain e.g `07-tendermint-0`
- **Creator** (`string`): string identifying the creator of this query (not necessarily an address), e.g. a module's name
- **ChainId** (`string`): chain ID of the data source where data should be queried from
- **QueryParameters** (`[]byte`): query parameters, e.g. `append(banktypes.CreateAccountBalancesPrefix(add1), []byte("stake")...)`
- **BlockRepeat** (`uint64`): the period length, e.g. 10 blocks
- **LastHeightExecuted** (`uint64`): the last local height at which this interchain query was invoked
- **MaxResults** (`uint64`): the total number of concurrent results this `PeriodicICQ` should store. The results will be overwritten once exceeded. E.G `30` meaning the `31st` result will overwrite the `1st` result.

List of `PendingICQInstance`, each based on an `PeriodicICQ` instance:

- **Id** (`uint64`): basic unique identifier
- **TimeoutHeight** (`uint64`): height at which a result for this query will be considered invalid, calculated by adding _periodic_'s timeout height padding to the current height
- **TargetHeight** (`uint64`): from _periodic_
- **PeriodicId** (`uint64`): ID of _periodic_

List of `PendingICQRequest`, each based on an `PeriodicICQ` and `PendingICQInstance`:

- **Id** (`uint64`): basic unique identifier
- **Path** (`uint64`): from _periodic_
- **TimeoutHeight** (`uint64`): height at which a result for this query will be considered invalid, calculated by adding _periodic_'s timeout height padding to the current height
- **TargetHeight** (`uint64`): from _periodic_
- **ClientId** (`uint64`): from _periodic_
- **Creator** (`uint64`): from _periodic_
- **QueryParameters** (`uint64`): from _periodic_
- **PeriodicId** (`uint64`): ID of _periodic_

List of `ICQResult`, each based on an `PeriodicICQ` instance:

- **PeriodicId** (`uint64`): matches the ID of the _periodic_ query
- **LastResultId** (`string`): matches the ID of the last _datapointresult_

List of `DataPointResult`, each based on an `PendingICQInstance` and a `MsgSubmitICQResult` instance:


- **Id** (`string`): basic unique identifier
- **LocalHeight** (`uint64`): local height at which result was recorded, i.e. upon _MsgSubmitICQResult_ handling
- **TargetHeight** (`uint64`): the height at which the result was retrieved from the target chain
- **Data** (`[]byte`): encoded query result from _MsgSubmitICQResult_
- **PrevDataPointId**: (`string`): id of the previous _datapointresult_ linking them all together

List of `ICQTimeouts`, each based on an `PeriodicICQ` instance:

- **PeriodicId** (`uint64`): identifier which points to the _periodic_ query
- **Timeouts** (`uint64`): number of times pending query instances timed out for the _periodic_ query
- **LastTimeoutHeight** (`uint64`): last local height that a _periodic_ query has timedout

For all the above values, a count is also stored.

## Transactions

`MsgSubmitICQResult`: used by relayer to submit an interchain query result. If the result arrives too late (i.e. timeout) it is ignored. `ICQResult` is updated as well as `DataPointResult` is created or updated depending on if the max results limit is reached.

- **QueryId** (`uint64`): matches the ID of _query_ that this result is for
- **Result** (`[]byte`): encoded query result
- **Height** (` ibc.core.client.v1.Height`): the height at which the result was retrieved from the target chain
- **FromAddress** (`string`): relayer that submitted this result, i.e. the one that signed this message
- **Proof** (`*crypto.ProofOps`): data validity proof
- **PeriodicId** (`uint64`): ID of _periodic_ that this result corresponds to

## Params

n/a

## BeginBlock

- `ProcessBeginBlockPendingQueries`: iterates over all `PendingICQInstance`s and if the timeout height has been exceeded, the particular instance is timed-out. This means that a `ICQTimeouts` instance is updated and the original `PendingICQInstance` is deleted.

## EndBlock

- `ExecutePeriodicQueries`: iterates over all `PeriodicICQ`s and creates `PendingICQInstance`s if enough blocks have elapsed since the last time the query was invoked. This is the case when the sum of _LastHeightExecuted_ (e.g. 1,000,000) and _BlockRepeat_ (e.g. 5) is equal to the current block height (e.g. 1,000,005).

## Dependencies

No special dependencies

## Hooks

n/a

## Future

- Programmatic relayer incentivisation/subsidisation