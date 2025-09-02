# git flow support

## git flow support - List all support branches

**Description**

List all local support branches

**Synopsis**

```
git flow support [list] [-h] [-v]
```

**Options**

- `-h,--[no]help` show this help
- `-v,--[no]verbose` verbose (more) output

## git flow support start - Start a support branch

**Description**

Start a new support branch name <version>, optionally basing it on <base> instead of <master>

**Synopsis**

```
git flow support start [-h] [-F] <version> [<base>]
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-F,--[no]fetch` fetch from origin before performing start

## git flow support rebase - Perform a rebase

**Description**

Rebase <name> on <base_branch>

**Synopsis**

```
git flow support rebase [-h] [-i] [-p] [<name|nameprefix>]
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-i,--[no]interactive` do an interactive rebase
- `-p, --[no]preserve-merges` preserve merges