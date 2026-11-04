package precedence

# Precedence

Which source outranks which, highest first. Composition settles a conflict between tiers on its
own; anything the tiers cannot settle goes to the owner (Sajmani) and is recorded in
[`decisions.md`](decisions.md).

The tiers below were carried over unchanged from specman's source manifest (`spec/sources.md`
at [`dc25e3b`](https://github.com/Sajmani/birdsync/blob/dc25e3b6756bd59c2de931e6d080d47318617da4/spec/sources.md)).

| Tier | Sources | Why this tier |
| --- | --- | --- |
| **Mandatory** | [`inatterms`](imported/inaturalist.org/inatterms/inatterms.md), [`mlterms`](imported/ebird.org/mlterms/mlterms.md) | Terms of service and media terms. Breaching them can cost the user their account, or put another person's copyrighted media in it |
| **Governing** | [`inatapi`](imported/inaturalist.org/inatapi/inatapi.md) | The API operator's recommended practices. iNaturalist can throttle or block a client that ignores them. Only iNaturalist can grant an exception |
| **Local: tech** | [`tech`](local/tech/tech.md) | The project's own constraints. Safety invariants outrank convenience |
| **Local: product** | [`product`](local/product/) | The project's own product requirements |
| **Advisory** | `go-practice` (absorbed into `tech`; see [`imports.md`](imports.md)), `uuid` (its license; [decisions.uuidAdopted](decisions.md#uuidadopted)) | Loses to a local requirement, and fills in where they say nothing |

Specman's manifest ranked no source within the Local tier. `tech` is placed above `product` here
because birdsync's `AGENTS.md` and `tech` both say the safety invariants outrank convenience.
The owner confirmed this ranking on 2026-11-04 ([decisions.localTechAboveProduct](decisions.md#localtechaboveproduct)).
