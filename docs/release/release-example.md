# Release Example

This file is a visual example of the release process outlined in [Release Process](release-process.md) document.

1. Initial state

    <img src="images/01-initial.svg" alt="initial master branch" />

2. Cut release branch `release/0.10`

    <img src="images/02-release-branch.svg" alt="cut release branch release/0.10" />

3. Commit `bugfix 3` to `master`

    <img src="images/03-bugfix.svg" alt="commit bugfix 3 to master" />

4. Cherry pick `bugfix 3` to `release/0.10`

    <img src="images/04-cherry-pick.svg" alt="cherry pick bugfix 3 to release/0.10" />

5. Create tag `v0.10.0-rc.1`

    <img src="images/05-rc-tag.svg" alt="create tag v0.10.0-rc.1" />

6. Commit `bugfix 4` to `master`

    <img src="images/06-bugfix.svg" alt="commit bugfix 4 to master" />

7. Cherry pick `bugfix 4` to `release/0.10`

    <img src="images/07-cherry-pick.svg" alt="cherry pick bugfix 4 to release/0.10" />

8. Create tag `v0.10.0-rc.2`

    <img src="images/08-rc-tag.svg" alt="create tag v0.10.0-rc.2" />

9. Create production tag `v0.10.0`

    <img src="images/09-production-tag.svg" alt="create production tag v0.10.0" />

10. Merge hotfix PR directly to `release/0.10`

    <img src="images/10-hotfix.svg" alt="merge hotfix to release/0.10" />

11. Create tag `v0.10.1-rc.1`. Since production tag `v0.10.0` has already been created, it is no longer permissible to create
any `v0.10.0-rc.X` tags

    <img src="images/11-rc-tag.svg" alt="create tag v0.10.1-rc.1" />

12. Create production tag `v0.10.1`

    <img src="images/12-production-tag.svg" alt="create production tag v0.10.1" />

(Note for document maintainers: the source diagrams can be found [here](https://link.excalidraw.com/l/1XPZRMVbRNH/32yMzzv0C50).
Please be sure to use consistent svg format by exporting from Excalidraw. Output svgs should be scaled down to 40% of the original
size, for the sake of consistency.)