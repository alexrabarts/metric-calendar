plugins {
    kotlin("jvm") version "1.9.23"
    `maven-publish`
}

group = "com.metricweek"
version = "1.0.0"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation("junit:junit:4.13.2")
}

java {
    sourceCompatibility = JavaVersion.VERSION_1_8
    targetCompatibility = JavaVersion.VERSION_1_8
}

kotlin {
    jvmToolchain(8)
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
                description.set("Metric Calendar date conversion — a precision decimal calendar system")
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
