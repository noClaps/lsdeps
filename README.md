# lsdeps

This is a CLI tool to show the total number of dependencies that any package has on the npm registry. This includes both direct dependencies, as well as indirect dependencies (dependencies of dependencies, for example). It currently includes packages listed in the `"dependencies"`, `"peerDependencies"` and `"optionalDependencies"` keys of the `package.json` for the package.

## Installation

You can install it using Homebrew on macOS/Linux:

```sh
brew install noclaps/tap/lsdeps
```

or you can build from source:

```sh
git clone https://github.com/noClaps/lsdeps.git && cd lsdeps
go build
```

## Usage

```
USAGE: lsdeps <package> [--skip-optional] [--skip-peer] [--version <version>]

ARGUMENTS:
  <package>              The npm package to count dependencies for.

OPTIONS:
  --skip-optional, -o    Skip counting optional dependencies.
  --skip-peer, -p        Skip counting peer dependencies.
  --version <version>    The version of the package being fetched.
  --help, -h             Display this help message and exit.
```

You can use the tool simply by running:

```sh
lsdeps astro # or any other package on npm
```

You can also skip counting peer dependencies with the `-p` or `--skip-peer` flag:

```sh
lsdeps astro -p
lsdeps astro --skip-peer
```

and optional dependencies with the `-o` or `--skip-optional` flag:

```sh
lsdeps astro -o
lsdeps astro --skip-optional
```

If you want to skip both, you can chain them together:

```sh
# All of these do the same thing
lsdeps astro -o -p
lsdeps astro -p -o
lsdeps astro --skip-peer --skip-optional
lsdeps astro --skip-optional --skip-peer
```

You can fetch dependencies for a specific package version with the `-v` or `--version` option:

```sh
lsdeps astro -v 4.0.0
```

The valid values for version are `latest`, `next` and a specific semantic version like `1.0.0` or `1.5.0-beta.1`. If the version cannot be parsed, `latest` is used instead.

**NOTE**: The version parsing is not perfect and likely needs to be rewritten in a much better way. However, it works in the meantime.

You can view the help by using `-h` or `--help`:

```sh
lsdeps -h
lsdeps --help
```

## Motivation

I've noticed that a lot of packages and frameworks that I use have hundreds of dependencies. I'm sure most of these are necessary, but I feel like the JavaScript ecosystem as a whole has decended into dependency hell, with some of the most basic packages having dependencies I don't think they really need.

I hadn't yet learned programming when leftpad happened, nor was I on the internet at that point, so I never really felt the impact of it. I didn't even learn about it until a couple years ago! But it kind of seems like the community never fully learned their lesson. We're still completely fine with installing hundreds upon hundreds of dependencies for sometimes the most basic of packages. In fact, we've made our package managers optimise for that specifically: PNPM and Bun both cache dependencies to avoid having to install them over and over again.

This, for me, is a flaw in how I think software development on the web _should_ work. I'm not one of those people who thinks having dependencies at all is bad. You can't do _everything_ yourself, obviously. However, having that many dependencies creates an entire graph of potential failure, of bugs, of security issues. If any one of those dependencies breaks, your entire application comes crashing down with it.

I experienced this a few months ago, admittedly in a very low-stakes situation. I updated one of the Rehype dependencies I was using for displaying math on my blog, and suddenly the blog stopped building. Something somewhere in the dependency made it incompatible with whatever Astro was doing, and the changelog for the Rehype dependency also didn't help me understand how to fix it. I've also had issues in the past with some Astro plugins like `astro-compress`. No hate to either of those packages' developers, I'm sure they're trying their best, but the point is that this kind of stuff happens.

Anyway, to demonstrate the number of dependencies that different packages have, I built `lsdeps`. There might be other tools like it, but I couldn't find any. It's currently a command line tool, although if there's interest in the project I might build a GUI around it as well. It's not very optimised at the moment, but I'll try to fix that over time.
