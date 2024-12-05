# A note about experimental/WIP APIs

There are a number of APIs that are currently under active development. These APIs can be fully ignored.
All such APIs will have comments in the form

```
/////////////////////////////////////////////////////////////////////////////////////
// Experimental: the following definitions are experimental and subject to change. //
/////////////////////////////////////////////////////////////////////////////////////
```

The majority of the WIP APIs are for a project we are calling internally `EigenDA v2 Architecture`.
More on that below.

## Q: Which APIs are currently experimental?

The following APIs are currently experimental:
- `disperser/v2/*`
- `node/v2/*`
- `relay/*`

## Q: are APIs not marked with "Experimental" stable?

Yes. We are commited to maintaining backwards compatibility for all APIs that are not marked as experimental,
and any breaking changes will be made only after a long deprecation period and active communication with
all stakeholders. Furthermore, breaking API changes are expected to be rare.

## Q: Should I use experimental APIs?

No. No experimental APIs are currently deployed to any public environments. In general, assume
that experimental APIs are not functional absent messaging from the EigenDA team declaring otherwise.

## Q: Are experimental APIs stable?

No, although they will become more and more stable as they reach maturity.

## Q: What is "v2"?

The EigenDA v2 Architecture is a fundamental redesign of the protocol. The v2 Architecture improves robustness,
efficiency, and paves the way for upcoming features such as permissionless disperser instances
and data availability sampling.

We intend on publishing a more detailed roadmap and design overview in the near future, stay tuned!