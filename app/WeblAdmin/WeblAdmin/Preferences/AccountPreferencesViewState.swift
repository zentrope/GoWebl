//
//  AccountPreferencesViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/2/22.
//

import Foundation
import OSLog

@MainActor
class AccountPreferencesViewState: NSObject, ObservableObject {

    private let log = Logger("AccountPreferencesViewState")

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
                log.error("\(e.localizedDescription)")
                self.result = .failed(e.localizedDescription)
            }
        }
    }
}
