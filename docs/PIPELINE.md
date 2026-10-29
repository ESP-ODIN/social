# CI pipeline

Every change must pass `go vet` and **100% of the unit tests**. This is enforced at
two levels.

| Level | File | When | What it blocks |
| --- | --- | --- | --- |
| Local | `.githooks/pre-push` | before each `git push` | the push itself |
| GitHub Actions | `.github/workflows/tests.yml` | each push on any branch, each pull request | the merge into `main` (with the ruleset below) |

Both run the same commands:

```sh
go vet ./...
go test -race -count=1 ./...
```

A single failing test (or a compilation error) makes the run fail.

## Local pre-push hook

Enable it once per clone:

```sh
make hooks
```

This runs `git config core.hooksPath .githooks`. From then on, `git push` runs the
checks first and is cancelled if any of them fails:

```
pre-push: go vet + tests unitaires...
--- FAIL: TestSomething
pre-push: des tests ont échoué, push annulé.
```

The hook can be skipped with `git push --no-verify`. Do not rely on it alone: the
GitHub check below is the one that cannot be bypassed.

## GitHub Actions workflow

The `Tests` workflow has one job, `unit-tests`, running on `ubuntu-latest` with the
Go version read from `go.mod`. Its result appears on each commit and pull request,
and the logs are in the repository's **Actions** tab.

## Protecting `main`

The workflow only reports a status; GitHub blocks the merge once `main` is protected.
A repository admin creates the rule in **Settings → Rules → Rulesets → New branch
ruleset**:

- Name `protect-main`, enforcement **Active**, empty bypass list
- Target: **Include default branch**
- Rules:
  - Restrict deletions
  - Block force pushes
  - Require a pull request before merging
  - Require status checks to pass → add `unit-tests` (and `commitlint`)

The `unit-tests` check only appears in the list after the workflow has run at least
once on the repository.

## Running the checks manually

```sh
make vet
make test
# or exactly what the pipeline runs:
sh .githooks/pre-push
```

The race detector (`-race`) requires CGO and a C compiler.
