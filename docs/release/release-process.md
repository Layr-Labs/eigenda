# Release Management Process

## Table of Contents

1. [Feature Freeze & Release Branch Creation](#1-feature-freeze--release-branch-creation)
2. [Changes to a Release Branch](#2-changes-to-a-release-branch)
   - [Change Policy](#change-policy)
   - [Change Process](#change-process)
3. [Tagging a Release](#3-tagging-a-release)
4. [Github Release](#4-github-release)
   - [Creating the Release](#creating-the-release)
   - [Release Notes](#release-notes)

---

### 1. **Feature Freeze & Release Branch Creation**

Enacting a feature freeze helps to ensure that the code we publish to production environments is well tested and
mature. The start of a feature freeze is marked by the creation of a release branch, which allows development
against `master` to continue uninterrupted while the release is prepared.

#### Plan Feature Freeze

- A feature freeze may be tied either to a date scheduled in advance, or to the completion of a key feature.
- As a general rule of thumb, a feature freeze should be planned such that there are two weeks between the freeze
and the release on testnet.
- The team should be notified of an upcoming feature freeze as soon as it has been planned.

#### Enact Feature Freeze

A feature freeze is officially marked by the creation of a release branch:

- From latest `master` commit:
  - `git checkout master && git pull`
  - `git checkout -b release/0.<MINOR>`
    - **Note:** there is no patch number in the branch name. The same branch is used across multiple patch versions.
  - Example: `release/0.10`
- Push the branch:
  - `git push origin release/0.10`
  - GitHub policies are configured to automatically protect a branch prefixed with 'release', to prevent it from being
  directly pushed to or deleted.

Note: The current branch naming scheme is `release/0.<MINOR>`, so that a user can checkout and pull the release branch
without necessarily being aware of what the latest patch release is. Once we release the first major semver version,
the branch naming format will be changed to `release/<MAJOR>`, to enable a similar user flow (checking out the major
version release branch, and pulling without needing to know the latest minor or patch versions).

---

### 2. **Changes to a Release Branch**

#### Change Policy

- **High bar for inclusion**: Only critical bugfixes or business-critical features
  - Even bugfixes should not be reflexively included: only high-severity issues
- **Team consensus required**: Single engineer cannot make the decision
- **Public visibility**: Must have team discussion (e.g. Slack thread) before proceeding. Alternatively, management
may sign-off that a feature should be included after a feature freeze has been enacted. Note that even with management
sign-off, a PR targeting a release branch must still go through the standard peer-review process.

#### Change Process

- **If change is also needed on `master`:**
  1. Submit PR and merge into `master` first
  2. Cherry-pick the squashed commit into the release branch
- **If change is release-only:**
  - Submit PR directly against the release branch
- **⚠️ NEVER push directly to the release branch**

---

### 3. **Tagging a Release**

**⚠️ Tags are immutable:** NEVER force-push a tag to a different commit

#### Release Candidate Tags

- **Cut release candidate tags for initial testing** (e.g. preprod environments):
  - Tag format: `v<MAJOR>.<MINOR>.<PATCH>-rc.<NUMBER>`
  - Example: `v0.10.0-rc.1`
  - Release candidates enable iterative testing without causing the patch version to increase
  - Release candidate tags clearly indicate to operators and users that a release is **not production-ready**
  - Commands:
    - `git checkout release/0.10`
    - `git tag v0.10.0-rc.1`
    - `git push origin v0.10.0-rc.1`
- **Tag additional release candidates** with incremented RC number (e.g. `v0.10.0-rc.2`, `v0.10.0-rc.3`)

#### Production Release Tags

- **Tag first production release** when ready to deploy to testnet:
  - Tag format: `v<MAJOR>.<MINOR>.<PATCH>`
  - Example: `v0.10.0`
  - Commands:
    - `git checkout release/0.10`
    - `git tag v0.10.0`
    - `git push origin v0.10.0`
- **Additional release candidate tags** may be cut even after the first production release has been tagged
  - Do this when testing of a production release reveals that additional iterations are necessary

See the [Release Example](release-example.md) document for a step-by-step release procedure example.

---

### 4. **Github Release**

#### Creating the Release

- **When ready to make the release public:**
  - If necessary, tag final patch version from release branch HEAD
  - Create GitHub release via UI, targeting the most recent tag
  - **Note**: Release will likely have non-zero patch version

#### Release Notes

- Follow the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format
