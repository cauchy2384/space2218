# Questions

## What instrumentation this service would need to ensure its observability and operational transparency?

* Metrics:
    * Hardware resources (cpu/ram etc)
    * Go runtime (goroutines/gc etc)
    * Web (rps/latency/error codes etc)
* SLO:
    * Availability
    * Latency
* Alerting
* Log collector
* Tracing

## Why throttling is useful (if it is)? How would you implement it here?

Throttling may be useful because location calculation is a CPU-bound code (might become with service evolution).

There are multiple ways to implement it. It might be one or a combination of the following methods:
* Worker pool for business logic
* Throttling middleware in HTTP server
* Rate limiting by a sidecar service

## What we have to change to make ​DNS ​be able to service several sectors at the same time?

Service should be configured with an array of sector IDs instead of a single one.

If I understand the user story correctly an API should be extended to accept a new input parameter: sector ID.

Another option is to serve multiple navigation services on different ports: in a single application or running multiple applications with a different configuration.

## Our CEO wants to establish B2B integration with Mom's Friendly Robot Company by allowing cargo ships of MomCorp to use D​NS​. The only issue is - MomCorp software expects l​oc ​value in the location ​ field, but math stays the same. How would you approach this? What’s would be your implementation strategy?

I would recommend implementing a different service serving another version of API that is approved by both parties. This service should serve as an adapter to the existing API.

A more simple, but less flexible way is to prove a separate API in the same service and use an adapter pattern to call the original business logic.

## Atlas Corp mathematicians made another breakthrough and now our navigation math is even better and more accurate, so we started producing a new drone model, based on new math. How would you enable a scenario where ​DNS ​can serve both types of clients?

I would increment an API version and serve both old and new ones. The new drone model would call the API v2, while the old models will be using v1.

Another option is to differentiate models by some property (header an example) and make business logic work with both new and old models. Yet this solution is less flexible as both models are bound to the same API schema.

## In general, how would you separate technical decision ​to ​deploy something from business decision ​to​ release ​something?

The business decision is based on business requirements and should take known technical issues into consideration (nothing is ever perfect).

The technical decision is based on technical readiness. It should pass all technical requirements and known technical risks should be weighted.