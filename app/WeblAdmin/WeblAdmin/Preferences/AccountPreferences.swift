//
//  AccountPreferences.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI


@MainActor
class AccountPreferencesViewState: NSObject, ObservableObject {

    @Published var result: TestResult = .untested

    enum TestResult {
        case untested
        case succeeded
        case failed(String)
    }

    func testConnection(email: String, password: String) {
        Task {
            do {
                let client = WebClient()
                let rc = try await client.test()
                self.result = rc ? .succeeded : .failed("Unable to acquire an auth token.")
            } catch (let e) {
                self.result = .failed(e.localizedDescription)
            }
        }
    }
}

struct AccountPreferences: View {

    @StateObject private var state = AccountPreferencesViewState()
    @AppStorage("WAAccountEmail") private var email = ""
    @AppStorage("WAAccountPassword") private var password = ""

    var body: some View {
        VStack {
            Form {
                TextField("Email", text: $email)
                SecureField("Password", text: $password)
            }
            HStack(alignment: .center, spacing: 20) {
                Button("Test") {
                    state.testConnection(email: email, password: password)
                }
                switch state.result {
                    case .untested:
                        EmptyView()
                    case .succeeded:
                        Text("Success")
                            .foregroundColor(.green)
                    case .failed(let msg):
                        Text(msg)
                            .foregroundColor(.red)
                }
                Spacer()
            }
        }
        .padding()
        .frame(width: 400)
        .fixedSize()
    }
}
