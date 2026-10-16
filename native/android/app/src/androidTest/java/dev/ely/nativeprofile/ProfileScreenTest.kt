package dev.ely.nativeprofile

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.test.espresso.Espresso
import androidx.test.platform.app.InstrumentationRegistry
import org.junit.Assert.assertEquals
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
        finishComposition()
        compose.onNodeWithTag(buttonID).performScrollTo().assertIsEnabled().performClick()
        val pending = compose.runOnIdle {
            requireNotNull(state.value.pending) { "The enabled Save button must start a request" }.request.also {
                assertEquals("save-1", it.requestID)
                assertEquals(mapOf("displayName" to "Alice"), it.fields)
            }
        }
        compose.onNodeWithTag(buttonID).assertIsNotEnabled()
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextReplacement("Alice Cooper")
        finishComposition()
        compose.onNodeWithTag(inputID).assertTextContains("Alice Cooper").assertIsFocused()
        compose.runOnIdle {
            assertEquals("Alice Cooper", state.value.drafts["displayName"])
            assertEquals(pending, state.value.pending?.request)
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
        finishComposition()
        compose.onNodeWithTag(findButton(ready.region)!!.id).performScrollTo().assertIsEnabled().performClick()
        compose.runOnIdle {
            val pending = requireNotNull(state.value.pending) { "The enabled Save button must start a request" }.request
            assertEquals("save-1", pending.requestID)
            assertEquals(mapOf("displayName" to ""), pending.fields)
            state.value = state.value.acceptSave(ready.copy(requestID = pending.requestID, status = "invalid",
                acknowledgedEdits = pending.editRevisions, errors = mapOf("displayName" to "Enter a display name.")))
        }
        compose.onNodeWithText("Enter a display name.").assertExists()
        compose.onNodeWithTag(inputID).performScrollTo().performClick().performTextReplacement("Ada")
        compose.onNodeWithText("Enter a display name.").assertDoesNotExist()
    }

    private fun finishComposition() {
        // Hiding the real IME commits its composition without triggering Done/save.
        // A composition-disabled button alone must never count as an in-flight save.
        Espresso.closeSoftKeyboard()
        compose.waitForIdle()
    }
}
