package dev.ely.nativeprofile

import android.graphics.Bitmap
import androidx.compose.ui.semantics.SemanticsProperties
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import java.io.File
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import org.junit.Assert.assertTrue

/** Real launcher, ViewModel, HTTP client, native controls, and Go server. No mocked state. */
@RunWith(AndroidJUnit4::class)
class LiveProfileAppTest {
    @get:Rule val compose = createAndroidComposeRule<MainActivity>()

    @Test fun realAppLoadsValidatesAndSavesThroughGo() {
        // The debug endpoint is the emulator bridge to the host's Go server.
        val directory = screenshotDirectory()
        listOf("01-ready.png", "02-invalid.png", "03-saved.png").forEach { File(directory, it).delete() }
        compose.waitUntil(timeoutMillis = 30_000) {
            compose.onAllNodesWithTag("display-name").fetchSemanticsNodes().isNotEmpty()
        }
        compose.onNodeWithTag("display-name").assertIsEnabled()
        compose.onNodeWithTag("save-profile").assertIsEnabled()
        screenshot("01-ready.png")

        compose.onNodeWithTag("display-name").performScrollTo().performClick().performTextClearance()
        compose.onNodeWithTag("save-profile").performScrollTo().performClick()
        val validation = "Use 2 to 60 characters for your display name."
        compose.waitUntil(timeoutMillis = 30_000) {
            compose.onAllNodesWithText(validation).fetchSemanticsNodes().isNotEmpty()
        }
        compose.onNodeWithText(validation).assertExists()
        compose.onNodeWithTag("request-status").assertTextEquals("No changes saved.")
        compose.onNodeWithTag("display-name").performScrollTo()
        screenshot("02-invalid.png")

        compose.onNodeWithTag("display-name").performClick().performTextReplacement("  Android Emulator  ")
        compose.onNodeWithText(validation).assertDoesNotExist()
        compose.onNodeWithTag("save-profile").performScrollTo().performClick()
        compose.waitUntil(timeoutMillis = 30_000) {
            compose.onAllNodesWithText("Profile saved.").fetchSemanticsNodes().isNotEmpty()
        }
        compose.onNodeWithTag("request-status").assertTextEquals("Profile saved.")
        // Canonicalization proves a server reply was adopted by the actual app.
        compose.onNodeWithTag("display-name").assertTextContains("Android Emulator")
        compose.onNodeWithTag("display-name").assert(SemanticsMatcher.expectValue(
            SemanticsProperties.EditableText, AnnotatedString("Android Emulator")))
        compose.onNodeWithTag("save-profile").assertIsEnabled()
        compose.onNodeWithTag("request-status").performScrollTo()
        screenshot("03-saved.png")
    }

    private fun screenshotDirectory(): File {
        val target = InstrumentationRegistry.getInstrumentation().targetContext
        val external = requireNotNull(target.getExternalFilesDir(null)) { "Emulator external app files are unavailable" }
        val directory = File(external, "screenshots")
        assertTrue("Cannot create screenshot output directory", directory.isDirectory || directory.mkdirs())
        return directory
    }

    private fun screenshot(name: String) {
        compose.waitForIdle()
        val image = requireNotNull(InstrumentationRegistry.getInstrumentation().uiAutomation.takeScreenshot()) {
            "The emulator did not return a screenshot"
        }
        try {
            File(screenshotDirectory(), name).outputStream().use {
                assertTrue("Could not encode screenshot", image.compress(Bitmap.CompressFormat.PNG, 100, it))
            }
        } finally { image.recycle() }
    }
}
