import XCTest
@testable import ElyNativeCore

final class ControlsTests: XCTestCase {
    func fixture(_ name: String = "profile-controls-ready") throws -> ProfileEnvelope {
        try EnvelopeDecoder.decode(Data(contentsOf: ProfileStateTests.fixtureDirectory.appendingPathComponent("\(name).json")))
    }
    func testMultipleFieldsReconcileIndependentlyAndRetryExactSnapshot() throws {
        let ready = try fixture()
        var state = ProfileState()
        try state.beginLoad(requestID: ready.requestID)
        try state.receive(ready)
        state.edit(field: "displayName", value: "  Grace Hopper  ")
        state.edit(field: "company", value: "  Example Labs  ")
        state.edit(field: "notifications", value: "true")
        let pending = try state.beginSave()
        state.edit(field: "company", value: "Later company")
        state.fail(requestID: pending.requestID, message: "Offline", retryable: true)
        XCTAssertEqual(try state.retry(), pending)
        var saved = try fixture("profile-controls-saved")
        saved.requestID = pending.requestID
        if case .save(let save) = pending.kind { saved.acknowledgedEdits = save.editRevisions }
        try state.receive(saved)
        XCTAssertEqual(state.drafts["company"], "Later company")
        XCTAssertEqual(state.baseline["company"], "Example Labs")
        XCTAssertEqual(state.drafts["displayName"], "Grace Hopper")
        XCTAssertEqual(state.drafts["notifications"], "true")
    }
    func testUnknownMissingDuplicateAndInvalidBooleanFieldsFailClosed() throws {
        var ready = try fixture()
        ready.fields["undeclared"] = "secret"
        XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
        ready = try fixture(); ready.fields["notifications"] = "yes"
        XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
        ready = try fixture(); ready.fields.removeValue(forKey: "company")
        XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
        ready = try fixture()
        var duplicate = ready.region.children[0].children[2].children[0]
        duplicate.id = "other-input"
        ready.region.children[0].children[2].children.append(duplicate)
        XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
    }
    func testTrailingLineBreakFieldNamesAndNodeIDsAreRejected() throws {
        for suffix in ["\n", "\r", "\r\n", "\u{2028}"] {
            var ready = try fixture()
            let value = ready.fields.removeValue(forKey: "company")!
            ready.fields["company" + suffix] = value
            ready.region.children[0].children[2].children[1].props.name = "company" + suffix
            XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
            ready = try fixture()
            ready.region.children[0].children[2].children[1].id += suffix
            XCTAssertThrowsError(try EnvelopeDecoder.validate(ready))
        }
    }
    func testRealGoControlsSaveAndReload() async throws {
        guard let base = ProcessInfo.processInfo.environment["ELY_NATIVE_CONTROLS_BASE_URL"] else {
            throw XCTSkip("Set ELY_NATIVE_CONTROLS_BASE_URL with the controls server running.")
        }
        let transport = try HTTPProfileTransport(baseURL: XCTUnwrap(URL(string: base)), allowLocalHTTP: true)
        var state = ProfileState()
        let load = try state.beginLoad()
        let ready = try await transport.send(load)
        try state.receive(ready)
        XCTAssertEqual(state.drafts.count, 3)
        state.edit(field: "company", value: "  Native Company  ")
        state.edit(field: "notifications", value: "true")
        let save = try state.beginSave()
        let saved = try await transport.send(save)
        try state.receive(saved)
        XCTAssertEqual(state.drafts["company"], "Native Company")
        XCTAssertEqual(state.drafts["notifications"], "true")
        let reload = try state.beginLoad(discardingDrafts: true)
        let reloaded = try await transport.send(reload)
        try state.receive(reloaded)
        XCTAssertEqual(state.drafts["notifications"], "true")
    }
}
