// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "MetricCalendar",
    platforms: [
        .macOS(.v10_15),
        .iOS(.v13),
        .watchOS(.v6),
        .tvOS(.v13),
    ],
    products: [
        .library(
            name: "MetricCalendar",
            targets: ["MetricCalendar"]
        ),
    ],
    targets: [
        .target(
            name: "MetricCalendar",
            path: "Sources/MetricCalendar"
        ),
        .testTarget(
            name: "MetricCalendarTests",
            dependencies: ["MetricCalendar"],
            path: "Tests/MetricCalendarTests"
        ),
    ]
)
