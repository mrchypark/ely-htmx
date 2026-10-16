package dev.ely.nativeprofile

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.*
import androidx.compose.ui.test.*
import androidx.compose.ui.focus.FocusManager
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.test.espresso.Espresso
import org.junit.Rule
import org.junit.Test

/** Real ViewModel + HTTP + Compose controls against the Go multi-field demo. */
class ControlsScreenTest {
    @get:Rule val compose = createComposeRule()
    @Test fun multipleFieldsAndNativeSwitchPersistThroughGo() {
        lateinit var focusManager: FocusManager
        compose.setContent {
            focusManager = LocalFocusManager.current
            val model = remember { ProfileViewModel("http://10.0.2.2:8095/native/profile") }
            val state by model.state.collectAsState()
            MaterialTheme { ProfileScreen(state, model::edit, model::save, model::retry, model::reload) }
        }
        compose.waitUntil(timeoutMillis = 90_000) { compose.onAllNodesWithTag("company").fetchSemanticsNodes().isNotEmpty() }
        compose.onNodeWithTag("company").performScrollTo().performClick().performTextReplacement("  Android Controls  ")
        compose.runOnIdle { focusManager.clearFocus(force = true) }
        Espresso.closeSoftKeyboard()
        compose.onNodeWithTag("notifications").performScrollTo().assertIsOff().performClick().assertIsOn()
        compose.onNodeWithTag("save-profile").performScrollTo().assertIsEnabled().performClick()
        compose.waitUntil(timeoutMillis = 30_000) { compose.onAllNodesWithText("Profile saved.").fetchSemanticsNodes().isNotEmpty() }
        compose.onNodeWithTag("company").performScrollTo().assertTextContains("Android Controls")
        compose.onNodeWithTag("notifications").performScrollTo().assertIsOn()
        compose.onNodeWithTag("reload-profile").performScrollTo().performClick()
        compose.onNodeWithText("Discard draft and reload").performClick()
        compose.waitUntil(timeoutMillis = 30_000) { compose.onAllNodesWithText("Edit your display name.").fetchSemanticsNodes().isNotEmpty() }
        compose.onNodeWithTag("company").performScrollTo().assertTextContains("Android Controls")
        compose.onNodeWithTag("notifications").performScrollTo().assertIsOn()
    }
}
