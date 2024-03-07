/**
A function that gets the total number of dependencies that a package has. This
includes all dependencies of dependencies as well. As it is a recursive
function, the output may take a while, especially if it needs to work through a
large dependency tree. The function will error if a package is not found.

@param {string} packageName
The name of the package

@param {string} version
The version of the package. Defaults to "latest". If the version of the package
starts with ^, the version will resolve to "latest".

@returns {Promise<number>}
The total number of dependencies.
*/
async function getDeps(packageName, version = "latest") {
  let total = 0;
  if (version.startsWith("^")) {
    version = "latest";
  }

  const packageData = await fetch(
    `https://registry.npmjs.com/${packageName}/${version}`,
  ).then((r) => r.json());

  if (packageData === "Not Found") {
    throw new Error("Package not found");
  }

  const { dependencies } = packageData;

  for (const dep in dependencies) {
    total += 1;
    total += await getDeps(dep, dependencies[dep]);
  }

  return total;
}

async function main() {
  /**
@type {string}
The name of the package to look up
*/
  const packageName = process.argv[2] ?? prompt("Enter a package name:");

  console.write("Counting dependencies...");

  /**
@type {number}
The number of dependencies that the package has.
*/
  const depCount = await getDeps(packageName);

  process.stdout.clearLine();
  process.stdout.cursorTo(0);
  console.log(
    `The "${packageName}" package has ${depCount} ${depCount === 1 ? "dependency" : "dependencies"}.`,
  );
}

main();
