import SwiftUI
import ElyNativeCore
import ElyNativeUI

@main
struct ElyProfileApp: App {
    private let transport: Result<HTTPProfileTransport, Error> = Result {
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
    init(transport: HTTPProfileTransport) {
        _model = StateObject(wrappedValue: ProfileViewModel(transport: transport))
    }
    var body: some View { NativeProfileView(model: model) }
}
