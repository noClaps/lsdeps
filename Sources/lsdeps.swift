import ArgumentParser
import Foundation

struct Package: Decodable {
  let dependencies: [String: String]?
  let peerDependencies: [String: String]?
  let optionalDependencies: [String: String]?
}

func fetch(url: String) async -> Package? {
  guard let urlObject = URL(string: url) else {
    print("Error parsing URL: \(url)")
    return nil
  }

  guard let (data, _) = try? await URLSession.shared.data(from: urlObject) else {
    print("Error parsing JSON")
    return nil
  }

  guard let package = try? JSONDecoder().decode(Package.self, from: data) else {
    print("Error decoding data: \(data)")
    return nil
  }

  return package
}

func getDeps(packageName: String, skipPeer: Bool, skipOptional: Bool) async -> [String]? {
  var deps: [String] = []
  guard let packageData = await fetch(url: "https://registry.npmjs.com/\(packageName)/latest")
  else {
    print("Error fetching dependencies for package \(packageName)")
    return nil
  }

  if let dependencies = packageData.dependencies {
    for dep in dependencies.keys {
      if !deps.contains(dep) {
        deps.append(dep)
      }
    }
  }
  if !skipPeer, let dependencies = packageData.peerDependencies {
    for dep in dependencies.keys {
      if !deps.contains(dep) {
        deps.append(dep)
      }
    }
  }
  if !skipOptional, let dependencies = packageData.optionalDependencies {
    for dep in dependencies.keys {
      if !deps.contains(dep) {
        deps.append(dep)
      }
    }
  }

  return deps
}

@main
struct lsdeps: AsyncParsableCommand {
  @Argument(help: "The npm package to count dependencies for")
  var package: String

  @Flag(
    name: [.customShort("p", allowingJoined: true), .long], help: "Skip counting peer dependencies")
  var skipPeer = false

  @Flag(
    name: [.customShort("o", allowingJoined: true), .long],
    help: "Skip counting optional dependencies")
  var skipOptional = false

  mutating func run() async throws {
    guard
      var depSet = await getDeps(
        packageName: package, skipPeer: skipPeer, skipOptional: skipOptional)
    else {
      return
    }

    var i = 0
    while i != depSet.count {
      let setPackage = depSet[i]

      print("Fetching dependencies for \(setPackage)")

      guard
        let deps = await getDeps(
          packageName: setPackage, skipPeer: skipPeer, skipOptional: skipOptional)
      else {
        continue
      }
      for dep in deps {
        if !depSet.contains(dep) {
          depSet.append(dep)
        }
      }
      i += 1
    }

    let depsCount = depSet.count

    print(
      """

      Name: \(package)
      URL: https://npmjs.com/package/\(package)
      Dependency count: \(depsCount)

      """)
  }
}
