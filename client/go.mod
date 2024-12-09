// We use a separate module for the client to allow dependencies to import it without importing all of proxy's main module's dependencies.
// This follows the recommendation in: https://go.dev/wiki/Modules#should-i-have-multiple-modules-in-a-single-repository
//
// Two example scenarios where it can make sense to have more than one go.mod in a repository:
// 1. [omitted]
// 2. if you have a repository with a complex set of dependencies, but you have a client API with a smaller set of dependencies.
//    In some cases, it might make sense to have an api or clientapi or similar directory with its own go.mod, or to separate out that clientapi into its own repository.
module github.com/Layr-Labs/eigenda-proxy/client

go 1.21.0
