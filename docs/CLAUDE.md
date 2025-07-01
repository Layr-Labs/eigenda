# CLAUDE.md - EigenDA Documentation

The `docs` directory contains documentation files describing the EigenDA system. This CLAUDE.md file causes all doc files to be read into 
context immediately. This isn't a strategy that should be followed for every CLAUDE.md file, but it makes sense for docs specifically,
since they are comparatively small and greatly improve Claude's ability to understand the repository.

## File Imports

1. @spec/src/glossary.md contains a glossary with basic EigenDA terms
2. @spec/src/introduction.md contains a basic introduction to the EigenDA system
3. @spec/src/integration/spec.md describes basic integration data flows
4. @spec/src/integration/spec/1-apis.md describes important APIs relevant to integrations
5. @spec/src/integration/spec/2-rollup-payload-lifecycle.md describes how user data moves through the EigenDA system
6. @spec/src/integration/spec/6-secure-integration.md contains considerations relevant to secure integrations
7. @release/release-process.md describes the EigenDA team's release processes
