import XCTest
@testable import ElyNativeCore
#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

final class TransportTests: XCTestCase {
    /// Opt-in because this makes real requests to the loopback Go example.
    func testRealGoLoadValidationSaveReplayAndCookieRetention() async throws {
        guard let origin = ProcessInfo.processInfo.environment["ELY_NATIVE_TEST_BASE_URL"] else {
            throw XCTSkip("Set ELY_NATIVE_TEST_BASE_URL=http://127.0.0.1:8094 with the Go example running.")
        }
        let url = try XCTUnwrap(URL(string: origin))
        let configuration = URLSessionConfiguration.ephemeral
        let transport = try HTTPProfileTransport(baseURL: url, allowLocalHTTP: true, configuration: configuration)
        var state = ProfileState()
        let load = try state.beginLoad()
        let ready = try await transport.send(load)
        try state.receive(ready)
        XCTAssertEqual(ready.status, .ready)
        XCTAssertEqual(state.drafts["displayName"], "Ada")
        state.edit(field: "displayName", value: "")
        let invalidRequest = try state.beginSave()
        let invalid = try await transport.send(invalidRequest)
        try state.receive(invalid)
        XCTAssertEqual(invalid.status, .invalid)
        XCTAssertEqual(state.drafts["displayName"], "")
        XCTAssertNotNil(state.fieldErrors["displayName"])
        XCTAssertEqual(invalid.sessionID, ready.sessionID)
        state.edit(field: "displayName", value: "  Swift Native  ")
        let save = try state.beginSave()
        let saved = try await transport.send(save)
        try state.receive(saved)
        XCTAssertEqual(saved.status, .saved)
        XCTAssertEqual(state.drafts["displayName"], "Swift Native")
        XCTAssertEqual(saved.dataVersion, 1)
        // Simulate retry after a lost success response. Exact bytes/ID must replay.
        let replay = try await transport.send(save)
        XCTAssertEqual(replay, saved)
        let reload = try state.beginLoad(discardingDrafts: true)
        let reloaded = try await transport.send(reload)
        try state.receive(reloaded)
        XCTAssertEqual(reloaded.sessionID, ready.sessionID)
        XCTAssertEqual(reloaded.fields["displayName"], "Swift Native")
        XCTAssertEqual(reloaded.dataVersion, 1)
    }

    func testChunkedResponseOverLimitIsRejectedBeforeDecode() async throws {
        let configuration = URLSessionConfiguration.ephemeral
        configuration.protocolClasses = [OversizedProtocol.self]
        let transport = try HTTPProfileTransport(baseURL: URL(string: "https://example.com")!, configuration: configuration)
        var state = ProfileState()
        do {
            _ = try await transport.send(state.beginLoad())
            XCTFail("An oversized chunked response must fail")
        } catch let error as ProtocolFailure {
            XCTAssertTrue(error.message.contains("1 MiB"))
        }
    }
}

private final class OversizedProtocol: URLProtocol, @unchecked Sendable {
    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }
    override func startLoading() {
        let response = HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: "HTTP/1.1", headerFields: ["Content-Type": "application/json"])!
        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: Data(repeating: 32, count: NativeProtocol.maximumResponseBytes))
        client?.urlProtocol(self, didLoad: Data([32]))
        client?.urlProtocolDidFinishLoading(self)
    }
    override func stopLoading() {}
}
