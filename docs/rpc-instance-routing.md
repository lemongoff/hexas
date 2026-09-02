# RPC instance routing

Hexas RPC clients support an optional per-call route target while keeping service ownership and allocation policy outside the framework.

## Contract

Service discovery publishes a stable `InstanceID` independently from the network address. The etcd resolver attaches that identity to gRPC address attributes; direct and Kubernetes targets use their complete `host:port` endpoint as the instance identity.

The default `p2c_ewma` balancer reads route directives created by:

```go
ctx = zrpc.WithRouteTarget(ctx, instanceID, zrpc.RouteRequire)
response, err := client.Call(ctx, request)
```

- `RouteRequire` selects exactly the requested ready instance. A missing or unready target returns gRPC `Unavailable`; it never falls back to another instance.
- `RoutePrefer` selects the requested instance when ready and otherwise falls back to normal P2C selection.
- A call without a route directive uses normal P2C selection.

Resolvers publish the complete discovered endpoint set. They do not randomly truncate membership before the balancer sees it, because truncation can hide a required state owner. This means clients create SubConns for all discovered endpoints; very large deployments that need lazy connections require a separate explicit design rather than a correctness-breaking hidden subset.

Duplicate non-empty instance identities are treated as unroutable by identity. Instances must therefore publish unique IDs within one discovery key.

## Boundary

Hexas only carries endpoint identity and executes a routing directive. It does not select a state owner, store ownership records, understand player IDs, reserve capacity, migrate state or retry on a different owner. Consumer projects must make those decisions before issuing a `RouteRequire` call and must use their own revision/fencing contract.

The discovery value changed atomically from `{"Addr":...,"ServerName":...}` to `{"Addr":...,"InstanceID":...}`. Mixed formats are not supported; deployers must publish Hexas and all consumers as one coordinated change and remove stale etcd values before rollback or rollout.
