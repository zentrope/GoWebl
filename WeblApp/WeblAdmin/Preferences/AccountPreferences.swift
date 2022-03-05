//
//  AccountPreferences.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct AccountPreferences: View {

    @StateObject private var state = AccountPreferencesViewState()

    @AppStorage("WAAccountEmail") private var email = ""
    @AppStorage("WAAccountPassword") private var password = ""

    var body: some View {
        VStack {
            Form {
                TextField("Email:", text: $email)
                SecureField("Password:", text: $password)
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
