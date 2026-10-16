package dev.ely.nativeprofile

import java.io.File
import kotlinx.serialization.json.*
import org.junit.Assume.assumeTrue
import kotlin.test.*

class ControlsTest {
    private fun fixture(name: String = "profile-controls-ready") = ProfileProtocol.decode(File(System.getProperty("ely.fixtures"), "$name.json").readBytes())
    @Test fun `multiple fields reconcile independently with exact retry`() {
        val ready = fixture()
        val pending = ProfileState().beginLoad(ready.requestID).acceptLoad(ready)
            .edit("displayName", "  Grace Hopper  ").edit("company", "  Example Labs  ")
            .edit("notifications", "true").beginSave("save-controls")
        val later = pending.edit("company", "Later company").networkFailed("save-controls").retrySave()
        assertEquals(pending.pending, later.pending)
        val saved = fixture("profile-controls-saved").copy(requestID = "save-controls", acknowledgedEdits = pending.pending!!.request.editRevisions)
        val final = later.acceptSave(saved)
        assertEquals("Later company", final.drafts["company"])
        assertEquals("Example Labs", final.baseline["company"])
        assertEquals("Grace Hopper", final.drafts["displayName"])
        assertEquals("true", final.drafts["notifications"])
    }
    @Test fun `invalid boolean and undeclared fields are rejected`() {
        val root = Json.parseToJsonElement(File(System.getProperty("ely.fixtures"), "profile-controls-ready.json").readText()).jsonObject
        for (fields in listOf(root.getValue("fields").jsonObject + ("notifications" to JsonPrimitive("yes")),
            root.getValue("fields").jsonObject + ("undeclared" to JsonPrimitive("secret")),
            root.getValue("fields").jsonObject - "company")) {
            val invalid = JsonObject(root + ("fields" to JsonObject(fields))).toString().toByteArray()
            assertFailsWith<ProtocolException> { ProfileProtocol.decode(invalid) }
        }
    }
    @Test fun `changed field set during save blocks reconciliation`() {
        val ready = fixture()
        val state = ProfileState().beginLoad(ready.requestID).acceptLoad(ready).beginSave("save-controls")
        val reply = fixture("profile-controls-saved").copy(requestID = "save-controls", fields = mapOf("displayName" to "Ada"))
        assertNotNull(state.acceptSave(reply).hardError)
    }
    @Test fun `real Go multiple fields save and reload`() {
        val base = System.getenv("ELY_NATIVE_CONTROLS_BASE_URL")
        assumeTrue("Controls Go server required", !base.isNullOrBlank())
        val client = ProfileHttpClient(base!!.trimEnd('/') + "/native/profile", allowDemoHttp = true)
        val ready = client.load(ProfileProtocol.newRequestID())
        var state = ProfileState().beginLoad(ready.requestID).acceptLoad(ready)
            .edit("company", "  Android Company  ").edit("notifications", "true").beginSave()
        val pending = state.pending!!
        state = state.acceptSave(client.save(pending))
        assertEquals("Android Company", state.drafts["company"])
        assertEquals("true", state.drafts["notifications"])
        assertEquals(state.baseline, client.load(ProfileProtocol.newRequestID()).fields)
    }
}
