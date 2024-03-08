# lsdeps

This is a CLI tool to show the total number of dependencies that any package has on the npm registry. This includes both direct dependencies, as well as indirect dependencies (dependencies of dependencies, for example). It currently only includes packages listed in the `"dependencies"` key of the `package.json` for the package.

## Usage

1. Clone the repository

   ```sh
   git clone https://github.com/noClaps/lsdeps.git
   cd lsdeps
   ```

2. Run the `lsdeps.go` file.

   ```sh
   go run lsdeps.go astro # You can give it a package name as a command line argument.
   ```

3. (Optional) Compile the script to a binary.
   ```sh
   go build lsdeps.go
   ```
   You can then place it in your PATH and call it from anywhere.

## Motivation

I've noticed that a lot of packages and frameworks that I use have hundreds of dependencies. I'm sure most of these are necessary, but I feel like the JavaScript ecosystem as a whole has decended into dependency hell, with some of the most basic packages having dependencies I don't think they really need.

I hadn't yet learned programming when leftpad happened, nor was I on the internet at that point, so I never really felt the impact of it. I didn't even learn about it until a couple years ago! But it kind of seems like the community never fully learned their lesson. We're still completely fine with installing hundreds upon hundreds of dependencies for sometimes the most basic of packages. In fact, we've made our package managers optimise for that specifically: PNPM and Bun both cache dependencies to avoid having to install them over and over again.

This, for me, is a flaw in how I think software development on the web _should_ work. I'm not one of those people who thinks having dependencies at all is bad. You can't do _everything_ yourself, obviously. However, having that many dependencies creates an entire graph of potential failure, of bugs, of security issues. If any one of those dependencies breaks, your entire application comes crashing down with it.

I experienced this a few months ago, admittedly in a very low-stakes situation. I updated one of the Rehype dependencies I was using for displaying math on my blog, and suddenly the blog stopped building. Something somewhere in the dependency made it incompatible with whatever Astro was doing, and the changelog for the Rehype dependency also didn't help me understand how to fix it. I've also had issues in the past with some Astro plugins like `astro-compress`. No hate to either of those packages' developers, I'm sure they're trying their best, but the point is that this kind of stuff happens.

Anyway, to demonstrate the number of dependencies that different packages have, I built `lsdeps`. There might be other tools like it, but I couldn't find any. It's currently a command line tool, although if there's interest in the project I might build a GUI around it as well. It's not very optimised at the moment, but I'll try to fix that over time. And, to prove the point I've been making, it has no dependencies at all (other than Go to actually run it, I suppose).
