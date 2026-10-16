package dev.ely.nativeprofile

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.test.platform.app.InstrumentationRegistry
import org.junit.Rule
import org.junit.Test

class ProfileScreenTest {
    @get:Rule val compose = createComposeRule()
    private fun ready(): Envelope = InstrumentationRegistry.getInstrumentation().context.assets
        .open("profile-ready.json").use { ProfileProtocol.decode(it.readBytes()) }
    private fun findInput(node: Node): Node? = if (node.type == "input") node else node.children.firstNotNullOfOrNull(::findInput)
    private fun findButton(node: Node): Node? = if (node.type == "button") node else node.children.firstNotNullOfOrNull(::findButton)

    @Test fun nativeFieldRetainsNewEditsAndFocusAfterFullRegionSave() {
        val ready = ready()
        val state = mutableStateOf(ProfileState().beginLoad(ready.requestID).acceptLoad(ready))
        compose.setContent {
            MaterialTheme {
                ProfileScreen(state.value, { field, text -> state.value = state.value.edit(field, text) },
                    { state.value = state.value.beginSave("save-1") }, {}, {})
            }
        }
        val inputID = findInput(ready.region)!!.id
        val buttonID = findButton(ready.region)!!.id
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextReplacement("Alice")
        compose.onNodeWithTag(buttonID).performScrollTo().performClick()
        compose.onNodeWithTag(buttonID).assertIsNotEnabled()
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextReplacement("Alice Cooper")
        compose.runOnIdle {
            val pending = state.value.pending!!.request
            state.value = state.value.acceptSave(ready.copy(requestID = pending.requestID, status = "saved", dataVersion = ready.dataVersion + 1,
                fields = mapOf("displayName" to "Alice"), acknowledgedEdits = pending.editRevisions))
        }
        compose.onNodeWithTag(inputID).assertTextContains("Alice Cooper").assertIsFocused()
        compose.onNodeWithTag(buttonID).assertIsEnabled()
    }

    @Test fun nativeValidationClearsAfterUserCorrectsField() {
        val ready = ready()
        val state = mutableStateOf(ProfileState().beginLoad(ready.requestID).acceptLoad(ready))
        compose.setContent {
            MaterialTheme {
                ProfileScreen(state.value, { field, text -> state.value = state.value.edit(field, text) },
                    { state.value = state.value.beginSave("save-1") }, {}, {})
            }
        }
        val inputID = findInput(ready.region)!!.id
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextClearance()
        compose.onNodeWithTag(findButton(ready.region)!!.id).performScrollTo().performClick()
        compose.runOnIdle {
            val pending = state.value.pending!!.request
            state.value = state.value.acceptSave(ready.copy(requestID = pending.requestID, status = "invalid",
                acknowledgedEdits = pending.editRevisions, errors = mapOf("displayName" to "Enter a display name.")))
        }
        compose.onNodeWithText("Enter a display name.").assertExists()
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextReplacement("Ada")
        compose.onNodeWithText("Enter a display name.").assertDoesNotExist()
    }
}
