# git flow hotfix

## git flow hotfix - Lists all hotfix branches

**Description**

Lists all local hotfix branches

**Synopsis**

```
git flow hotfix [list] [-h] [-v]
```

**Options**

- `-h,--[no]help` show this help
- `-v,--[no]verbose` verbose (more) output

## git flow hotfix start - Start a hotfix branch

**Description**

Start new hotfix branch named <version>, optionally base it on <base> instead of the <master> branch

**Synopsis**

```
git flow hotfix start [-h] [-F] <version> [<base>]
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-F,--[no]fetch` fetch from origin before performing finish

## git flow hotfix finish - Finish a hotfix branch

**Description**

Finish hotfix branch <version>

**Synopsis**

```
git flow hotfix finish [-h] [-F] [-s] [-u] [-m | -f ] [-p] [-k] [-n] [-b] <version>
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-F,--[no]fetch` fetch from origin before performing finish
- `-s,--[no]sign` sign the release tag cryptographically
- `-u,--signingkey` use the given GPG-key for the digital signature (implies -s)
- `-m,--message` use the given tag message
- `-f,--messagefile` use the contents of the given file as tag message
- `-p,--[no]push` push to origin after performing finish
- `-k,--[no]keep` keep branch after performing finish
- `--[no]keepremote` keep the remote branch
- `--[no]keeplocal` keep the local branch
- `-n,--[no]notag` don't tag this release
- `-T,--tagname` Use given tag name
- `-b,--[no]nobackmerge` don't back-merge master, or tag if applicable, in develop

## git flow hotfix publish - Publish the hotfix branch

**Description**

Start sharing hotfix <name> on $ORIGIN

**Synopsis**

```
git flow hotfix publish [-h] <name>
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them

## git flow hotfix delete - Delete a hotfix branch

**Description**

Deletes a given hotfix branch

**Synopsis**

```
git flow hotfix delete [-h] [-f] [-r] <name>
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-f,--[no]force` force deletion
- `-r,--[no]remote` delete remote branch

## git flow hotfix rebase - Perform a rebase

**Description**

Rebase <name> on <base_branch>

**Synopsis**

```
git flow hotfix rebase [-h] [-i] [-p] [<name|nameprefix>]
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them
- `-i,--[no]interactive` do an interactive rebase
- `-p, --[no]preserve-merges` preserve merges

## git flow hotfix rename - Rename a hotfix branch

**Description**

Rename branch <name> to <new_name>

**Synopsis**

```
git flow hotfix rename [-h] <new_name> [<name>]
```

**Options**

- `-h,--[no]help` show this help
- `--showcommands` Show git commands while executing them