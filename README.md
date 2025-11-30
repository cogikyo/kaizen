# kaizen

> [!NOTE] _change for the better_
>
> Personal code kata practice manager

## Install

```bash
go install ./cmd/kz
```

## Usage

```bash
kz init                       # initialize repo
kz init algorithms            # create section

kz new two-sum                # create problem (interactive)
kz solve two-sum              # open in $EDITOR
kz done two-sum               # run tests, record result

kz                            # show status + suggest next (history for last 6 months)
kz history                    # practice history (6 months)
kz history 06-01              # history starting from yy-mm
kz random                     # pick random (prioritizes due)
kz review                     # problems due for review

kz list                       # all problems
kz stats                      # detailed breakdown
kz test two-sum               # run tests
```

## Structure

```
algorithms/
  two-sum/
    solution.go
    solution_test.go
    README.md           # frontmatter with kyu, tags, source
.kaizen/
  kaizen.db             # SQLite (commit to git for sync)
```

## Kyu Levels

| Level | Name      |
| ----- | --------- |
| 1     | elite     |
| 2     | difficult |
| 3     | hard      |
| 4     | medium    |
| 5     | normal    |
| 6     | easy      |
