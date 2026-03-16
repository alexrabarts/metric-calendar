plugins {
    java
    `maven-publish`
}

group = "com.metricweek"
version = "1.0.0"

repositories {
    mavenCentral()
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

dependencies {
    testImplementation("junit:junit:4.13.2")
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            groupId = "com.metricweek"
            artifactId = "metric-calendar"
            version = "1.0.0"
            from(components["java"])
            pom {
                name.set("Metric Calendar")
                description.set("Metric Calendar date conversion — a rational decimal calendar system")
                url.set("https://metricweek.com")
                licenses {
                    license {
                        name.set("MIT License")
                        url.set("https://opensource.org/licenses/MIT")
                    }
                }
            }
        }
    }
}
