import ArgumentParser
import Foundation

struct Package: Decodable {
    let dependencies: [String: String]?
    let peerDependencies: [String: String]?
    let optionalDependencies: [String: String]?
}

extension Data {
    func toJSON<T>(_ type: T.Type) throws -> T where T: Decodable {
        return try JSONDecoder().decode(T.self, from: self)
    }

    func toString() -> String {
        return String(decoding: self, as: UTF8.self)
    }
}

func fetch(url: String) async -> Package? {
    guard let urlObject = URL(string: url) else {
        print("Error parsing URL: \(url)")
        return nil
    }

    var data: Data
    do {
        (data, _) = try await URLSession.shared.data(from: urlObject)
    } catch {
        print("Error fetching dependencies")
        print(error)
        return nil
    }

    var package: Package
    do {
        package = try data.toJSON(Package.self)
    } catch {
        let errorVal = data.toString()
        if errorVal == "\"Not Found\"" {
            print("Package doesn't exist")
            return nil
        }

        print("Error decoding JSON")
        print(error)
        return nil
    }

    return package
}

func parseVersion(version: String) throws -> String {
    let versionRegex = #/([0-9]\.[0-9]\.[0-9])(-(alpha|beta|rc)\.[0-9]+)?/#
    if try versionRegex.wholeMatch(in: version) != nil {
        return version
    }

    let versionStart = version[version.startIndex]
    let versionNumber = version[version.index(version.startIndex, offsetBy: 1)...]
    if versionStart == "^" || versionStart == "~",
        try versionRegex.wholeMatch(in: versionNumber) != nil
    {
        return String(versionNumber)
    }

    if version == "next" { return version }

    return "latest"
}

func getDeps(packageName: String, skipPeer: Bool, skipOptional: Bool, version: String)
    async throws
    -> [String: String]?
{
    let version = try parseVersion(version: version)
    var deps: [String: String] = [:]
    guard
        let packageData = await fetch(
            url: "https://registry.npmjs.com/\(packageName)/\(version)")
    else {
        print("Error fetching dependencies for package \(packageName)@\(version)")
        return nil
    }

    if let dependencies = packageData.dependencies {
        for (dep, version) in dependencies {
            if !deps.contains(where: { (key: String, value: String) in
                key == dep && value == version
            }) {
                deps[dep] = version
            }
        }
    }
    if !skipPeer, let dependencies = packageData.peerDependencies {
        for (dep, version) in dependencies {
            if !deps.contains(where: { (key: String, value: String) in
                key == dep && value == version
            }) {
                deps[dep] = version
            }
        }
    }
    if !skipOptional, let dependencies = packageData.optionalDependencies {
        for (dep, version) in dependencies {
            if !deps.contains(where: { (key: String, value: String) in
                key == dep && value == version
            }) {
                deps[dep] = version
            }
        }
    }

    return deps
}

@main
struct lsdeps: AsyncParsableCommand {
    @Argument(help: "The npm package to count dependencies for.")
    var package: String

    @Flag(
        name: [.customShort("p", allowingJoined: true), .long],
        help: "Skip counting peer dependencies.")
    var skipPeer = false

    @Flag(
        name: [.customShort("o", allowingJoined: true), .long],
        help: "Skip counting optional dependencies.")
    var skipOptional = false

    @Flag(
        name: .long, help: "Hide the \"Fetching dependencies for...\" messages.")
    var silent = false

    @Option(name: .shortAndLong, help: "The version of the package being fetched.")
    var version = "latest"

    mutating func run() async throws {
        if !silent {
            print("Fetching dependencies for \(package)@\(version)")
        }

        guard
            var depSet = try await getDeps(
                packageName: package, skipPeer: skipPeer,
                skipOptional: skipOptional, version: version)
        else {
            return
        }

        var i = 0
        while i != depSet.count {
            let setPackage = Array(depSet.keys)[i]
            let setPackageVersion = depSet[setPackage]!

            if !silent {
                print("Fetching dependencies for \(setPackage)@\(setPackageVersion)")
            }

            guard
                let deps = try await getDeps(
                    packageName: setPackage, skipPeer: skipPeer,
                    skipOptional: skipOptional, version: setPackageVersion)
            else {
                continue
            }

            for (dep, version) in deps {
                if !depSet.contains(where: { (key: String, value: String) in
                    key == dep && value == version
                }) {
                    depSet[dep] = version
                }
            }
            i += 1
        }

        let depsCount = depSet.count

        print(
            """

            Name: \(package)
            URL: https://npmjs.com/package/\(package)/v/\(version)
            Dependency count: \(depsCount)

            """)
    }
}
