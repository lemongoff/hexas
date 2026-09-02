# RPC instance routing

Hexas RPC clients support an optional per-call route target while keeping service ownership and allocation policy outside the framework.

## Contract

Service discovery publishes a stable `InstanceID` and explicit lifecycle `State` independently from the network address. `State` is either `ready` or `draining`. The etcd resolver attaches both values to gRPC address attributes; direct and Kubernetes targets use their complete `host:port` endpoint as the instance identity and have no framework-managed drain state.

The default `p2c_ewma` balancer reads route directives created by:

```go
ctx = zrpc.WithRouteTarget(ctx, instanceID, zrpc.RouteRequire)
response, err := client.Call(ctx, request)
```

- `RouteRequire` selects exactly the requested connected instance, including a `draining` instance. This exception lets a project finish explicitly fenced state transfer. A missing target returns gRPC `Unavailable`; it never falls back to another instance.
- `RoutePrefer` selects the requested ready instance and otherwise falls back to normal P2C selection; it never prefers a `draining` instance.
- A call without a route directive uses normal P2C selection and excludes every `draining` instance.

Resolvers publish the complete discovered endpoint set. They do not randomly truncate membership before the balancer sees it, because truncation can hide a required state owner. This means clients create SubConns for all discovered endpoints; very large deployments that need lazy connections require a separate explicit design rather than a correctness-breaking hidden subset.

Duplicate non-empty instance identities are treated as unroutable by identity. Instances must therefore publish unique IDs within one discovery key.

## Lifecycle

An etcd-backed `RpcServer` follows one ordered lifecycle:

1. Bind the listener and start gRPC serving.
2. Publish `ready` with the active lease, then expose healthy/ready probes.
3. `BeginDrain` changes the value on the same lease to `draining` and removes the endpoint from ordinary P2C selection.
4. `Stop` performs a gRPC graceful stop and synchronously revokes the registration.

Ready publication and stop are serialized, so a concurrent shutdown cannot re-register an already stopped server. `BeginDrain` and `Stop` are idempotent. Applications that need a propagation or migration window call `BeginDrain`, wait according to their own workload policy, then call `Stop`; Hexas does not invent a workload-specific delay.

`discov.Subscriber.Instances` returns only valid lifecycle records and isolates malformed values. `zrpc.Client.Close` releases its gRPC connection and resolver subscription; application composition roots should own and close long-lived clients instead of creating them per handler.

## Boundary

Hexas only carries endpoint identity and executes a routing directive. It does not select a state owner, store ownership records, understand player IDs, reserve capacity, migrate state or retry on a different owner. Consumer projects must make those decisions before issuing a `RouteRequire` call and must use their own revision/fencing contract.

The discovery value is `{"Addr":...,"InstanceID":...,"State":"ready|draining"}`. Values without a valid address, instance identity, or state are rejected. Mixed formats are not supported; deployers must publish Hexas and all consumers as one coordinated change and remove stale etcd values before rollback or rollout.
