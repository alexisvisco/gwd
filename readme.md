### What is this tool for ?

gwd stands for go workspace diff

**This is tool must be used with a go workspace file (go.work) introduced in the 1.18 release of go.**

It allows you to easily track :
- when a module change
- when packages of a module change and packages that import them

As an example, suppose a change is committed which modifies a package, libs/hypervisor. 
Suppose this package is imported by another package, services/controller.
gwd is used to inspect your files based on a list of changes and determine that both of these packages must be tested, although only the first package was changed.

If you think you've already read that, it's true. This excerpt comes from the blog of digitalocean in which they talk about [the organization of their monorepository](https://blog.digitalocean.com/cthulhu-organizing-go-code-in-a-scalable-repo/).

### Install

```
go install github.com/alexisvisco/gwd@latest 
```

### Commands

Do `gwd --help` to see the list of commands and their usage.

#### gwd --stdin 

This command is used to determine which modules changed based on a list of file change provided in stdin.

!Example:
```bash
$ cat go.work                 
go 1.18

use (
        ./services/competition-vacuum
        ./services/dummy
        ./services/vote
)

$ git diff 2022.0728.1850 --name-only | cat
.gitlab-ci.yml
.gitlab-ci/get-previous-ref.sh
Makefile
Makefile.common
services/competition-vacuum/cmd/competition-vacuum/main.go
services/competition-vacuum/internal/clients/kkbb/client.go
services/competition-vacuum/internal/clients/kkbb/types.go
services/competition-vacuum/internal/clients/ulule/client.go
services/competition-vacuum/internal/services/competition_vacuum/vacuum.go
tools/deploy.sh

$ git diff 2022.0728.1850 --name-only | gwd --stdin
services/competition-vacuum

$ git diff 2022.0728.1850 --name-only | gwd --stdin -v
"miimosa.com/services/competition-vacuum":
 - "services/competition-vacuum/internal/clients/ulule/client.go"
   imported by:
   ∟ "miimosa.com/services/competition-vacuum/cmd/competition-vacuum" (1 times) module "miimosa.com/services/competition-vacuum"
 - "services/competition-vacuum/internal/services/competition_vacuum/vacuum.go"
   imported by:
   ∟ "miimosa.com/services/competition-vacuum/cmd/competition-vacuum" (1 times) module "miimosa.com/services/competition-vacuum"
 - "services/competition-vacuum/cmd/competition-vacuum/main.go"
 - "services/competition-vacuum/internal/clients/kkbb/client.go"
   imported by:
   ∟ "miimosa.com/services/competition-vacuum/cmd/competition-vacuum" (1 times) module "miimosa.com/services/competition-vacuum"
```

As you can see only "miimosa.com/services/competition-vacuum" has been changed.

This command can take one argument, if it is specified it will show only diff for packages in the module specified.

#### check

This command is used to check if a module from the go workspace has change.
This command takes one argument, the name of the module, it can be a path or a module name.

### Flags 
Each of theses commands have in common 2 flags:

- `--stdin` read stdin, used in conjuncture with git diff --name-only for example.
- `--file` read a file, each line must be a file path

#### Global flags

- `--json or -j` output commands as json
- `--verbose or -v` output commands with more details that the default output
- `--go-work or -w` set the go workspace file name which is by default parsed from the the go.mod file

 
#### Why using a pkg folder ?

I don't use internal for open source because maybe you will use some packages for your usage.
