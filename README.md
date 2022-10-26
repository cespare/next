# Next

This is a collection of Go packages for testing additions and changes to
standard library packages as well as some popular x/ packages to make use of
type parameters (generics).

This is a testbed for experimentation. I intend to follow the normal semantic
versioning rules but there may be a lot of churn. As these sorts of packages are
accepted into the standard library or x/ repos, I intend to deprecate and freeze
my versions. For any which prove useful but which seem unlikely to be headed for
acceptance into the Go repos, I may move them to more permanent locations.

## Why?

Generics are a nice language feature with some more-or-less obvious candidate
uses within the standard library. The Go team is rightly being very cautious
and methodical about introducing these changes. But I'm impatient! I want to use
them now.

For example, golang/go#45955 describes a proposal for a new `slices` package.
This was later made available for use as
[golang.org/x/exp/slices](https://pkg.go.dev/golang.org/x/exp/slices)
to let folks use the proposed API while it was still under discussion.

Similarly, golang/go#47331 is a discussion about a new `container/set` package.
But, as of October 2022, no `container/set` package has been provided in
`x/exp`. So that's where this repo comes in.

## Packages

TODO: describe

* `github.com/cespare/next/container/ordmap`
* `github.com/cespare/next/container/set`
* `github.com/cespare/next/container/heap`
* `github.com/cespare/next/sync/syncutil`
* `github.com/cespare/next/sync/atomicutil`
* `github.com/cespare/next/sync/singleflight`

## License

Packages adapted from existing code in the Go project are released under the Go
project's license (see [LICENSE-THIRD-PARTY.txt](/LICENSE-THIRD-PARTY.txt)). The
other packages are released under the MIT license ([LICENSE.txt](/LICENSE.txt)).
