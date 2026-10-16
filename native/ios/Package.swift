// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "ElyNative",
    platforms: [.iOS(.v17), .macOS(.v14)],
    products: [
        .library(name: "ElyNativeCore", targets: ["ElyNativeCore"]),
        .library(name: "ElyNativeUI", targets: ["ElyNativeUI"])
    ],
    targets: [
        .target(name: "ElyNativeCore"),
        .target(name: "ElyNativeUI", dependencies: ["ElyNativeCore"]),
        .testTarget(name: "ElyNativeCoreTests", dependencies: ["ElyNativeCore"])
    ]
)
