/**
A function that gets the total number of dependencies that a package has. This
includes all dependencies of dependencies as well. The output may take a while,
especially if it needs to work through a large dependency tree. The function
will error if a package is not found.

@param {string} packageName
The name of the package

@returns {Promise<Set<string>>}
The set of dependencies for the given package.
*/
async function getDeps(packageName) {
  let depsSet = new Set();

  const packageData = await fetch(
    `https://registry.npmjs.com/${packageName}/latest`,
  ).then((r) => r.json());

  if (packageData === "Not Found") {
    throw new Error("Package not found");
  }

  const { dependencies } = packageData;

  for (const dep in dependencies) {
    depsSet.add(dep);
  }

  return depsSet;
}

/**
The name of the package to look up.
@type {string}
*/
const packageName = process.argv[2];

console.write("Counting dependencies...");

let pkgDeps = await getDeps(packageName);
for (const dep of pkgDeps) {
  await getDeps(dep).then((d) => d.forEach(pkgDeps.add, pkgDeps));
}

/**
The total number of dependencies that a package has.
@type {number}
*/
const depCount = pkgDeps.size;

process.stdout.clearLine();
process.stdout.cursorTo(0);
console.log(
  `The "${packageName}" package has ${depCount} ${depCount === 1 ? "dependency" : "dependencies"}.`,
);
