![Cthulhu](https://assets.digitalocean.com/ghost/2017/10/DivingIntoTheDepths_blog_pat.png)


### What is this tool for ?
gta, which stands for “Go Test Auto” inspects the git history to determine which files changed between a reference and an other reference, and uses this information to determine which packages must be tested for a given build (including packages that import the changed package).

As an example, suppose a change is committed which modifies a package, do/teams/example/droplet. Suppose this package is imported by another package, do/teams/example/hypervisor. gta is used to inspect the git history and determine that both of these packages must be tested, although only the first package was changed.

If you think you've already read that, it's true. This excerpt comes from the blog of digitalocean in which they talk about [the organization of their monorepository](https://blog.digitalocean.com/cthulhu-organizing-go-code-in-a-scalable-repo/).

In this blog they talk about a tool that allows them to reduce the build time from 20 minutes to 2-3 minutes average.

That's exactly what I did. The tool `gta` is now open source, as wanted some comments below the blog.


### Commands

- `gta changes --previous-ref="master"` which show the changes between master and the staging (uncommitted changes).
- `gta test --previous-ref="master"` which run the test for the package which changes and the packages that import changed packages.

Each of theses commands have in common 2 flags:

- `--previous-ref or -p` which is the old reference it can be a branch, tag or commit hash. 
- `--current-ref or -c` which is the current reference it can be a branch, tag or commit hash. 

#### Global flags

- `--json or -J` output commands as json
- `--verbose or -V` output commands with more details that the default output
- `--module-name or -M` set the module name which is by default parsed from the the go.mod file
