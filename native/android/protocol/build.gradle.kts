plugins { kotlin("jvm") }
kotlin { jvmToolchain(21) }
dependencies {
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.8.1")
    testImplementation(kotlin("test-junit"))
}
tasks.test {
    systemProperty("ely.fixtures", rootProject.file("../fixtures").absolutePath)
    inputs.property("liveServer", providers.environmentVariable("ELY_NATIVE_TEST_BASE_URL").orElse(""))
}
