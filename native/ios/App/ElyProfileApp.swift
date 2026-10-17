import SwiftUI
import ElyNativeCore
import ElyNativeUI

@main
struct ElyProfileApp: App {
    private let transport: Result<any ProfileTransport, Error> = Result {
        var origin = Bundle.main.object(forInfoDictionaryKey: "ElyAPIBaseURL") as? String ?? ""
        #if DEBUG
        if let override = ProcessInfo.processInfo.environment["ELY_API_BASE_URL"] { origin = override }
        let allowLocalHTTP = true
        #else
        let allowLocalHTTP = false
        #endif
        guard let url = URL(string: origin), !origin.isEmpty else {
            throw ProtocolFailure("Set ELY_API_BASE_URL to your HTTPS API origin in the Release build configuration.")
        }
        #if DEBUG
        // UI tests opt in explicitly; never bundle fixtures or replace the Go backend.
        if ProcessInfo.processInfo.arguments.contains("--selection-ui-test"),
           let scenario = SelectionUITestTransport.Scenario(rawValue: ProcessInfo.processInfo.environment["ELY_SELECTION_UI_SCENARIO"] ?? "") {
            let live = try HTTPProfileTransport(baseURL: url, allowLocalHTTP: allowLocalHTTP, configuration: .ephemeral)
            return SelectionUITestTransport(live: live, scenario: scenario)
        }
        #endif
        return try HTTPProfileTransport(baseURL: url, allowLocalHTTP: allowLocalHTTP)
    }

    var body: some Scene {
        WindowGroup {
            switch transport {
            case .success(let transport): ProfileAppContent(transport: transport)
            case .failure(let error):
                ContentUnavailableView("API configuration required", systemImage: "network", description: Text(error.localizedDescription))
            }
        }
    }
}

private struct ProfileAppContent: View {
    @StateObject private var model: ProfileViewModel
    init(transport: any ProfileTransport) {
        _model = StateObject(wrappedValue: ProfileViewModel(transport: transport))
    }
    var body: some View { NativeProfileView(model: model) }
}

#if DEBUG
/// Test-only fault injection around actual Go requests. Release contains no hooks.
private actor SelectionUITestTransport: ProfileTransport {
    enum Scenario: String, Sendable { case disabled, interruptedSave = "interrupted-save" }
    let live: HTTPProfileTransport
    let scenario: Scenario
    private var interruptedRequest: PendingRequest?
    private var didInterrupt = false

    init(live: HTTPProfileTransport, scenario: Scenario) {
        self.live = live; self.scenario = scenario
    }

    func send(_ request: PendingRequest) async throws -> ProfileEnvelope {
        if let interruptedRequest, case .save = request.kind {
            guard request == interruptedRequest else {
                throw ProtocolFailure("The UI test retry changed the original request snapshot.")
            }
        }
        var response = try await live.send(request)
        if scenario == .disabled {
            func disableSelections(_ node: inout NativeNode) {
                if ["checkbox", "select"].contains(node.type) { node.props.disabled = true }
                for index in node.children.indices { disableSelections(&node.children[index]) }
            }
            disableSelections(&response.region)
        }
        if scenario == .interruptedSave, !didInterrupt, case .save = request.kind {
            didInterrupt = true
            interruptedRequest = request
            // The server has committed, but the app sees an uncertain network result.
            // Leave enough time for actual native controls to create newer local edits.
            try await Task.sleep(nanoseconds: 15_000_000_000)
            throw TransportFailure("UI test interrupted a completed save. Retry the same request.", retryable: true)
        }
        if request == interruptedRequest { interruptedRequest = nil }
        return response
    }
}
#endif
