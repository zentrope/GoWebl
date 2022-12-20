//
//  AccountPreferences.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct AccountPreferences: View {

    @StateObject private var state = AccountPreferencesViewState()

    @State private var formName = ""
    @State private var formHost = ""
    @State private var formUser = ""
    @State private var formPassword = ""

    var body: some View {
        HStack(spacing: 10) {

            VStack(alignment: .leading, spacing: 10) {

                // Using a table so that we can set the selection value after a slight delay.
                Table(state.accounts, selection: $state.selectedAccount) {
                    TableColumn("Account") { account in
                        HStack(alignment: .center) {
                            Text(account.name)
                            Spacer()
                            if state.isDefault(id: account.id) {
                                Image(systemName: "checkmark")
                            }
                        }
                    }
                }
                .tableStyle(.bordered(alternatesRowBackgrounds: false))
                .labelsHidden()
                .onChange(of: state.selectedAccount) { accountId in                    
                    state.selectAccount(id: accountId)
                }

                HStack(spacing: 5) {
                    Button {
                        state.createNewAccount()
                    } label: {
                        Image(systemName: "plus")
                    }
                    .controlSize(.small)

                    Button {
                        if let id = state.selectedAccount {
                            state.deleteAccount(id: id)
                        }
                    } label: {
                        Image(systemName: "minus")
                    }
                    .disabled(state.selectedAccount == nil || state.accounts.isEmpty)
                    .controlSize(.small)
                    Spacer()
                }
            }
            .frame(width: 175)
            .alert(state.error?.localizedDescription ?? "Error", isPresented: $state.showAlert) {}

            if state.selectedAccount == nil {
                NoSelection()
            } else {
                AccountForm()
            }

        }
        .padding()
        .frame(width: 600, height: 200)
        .fixedSize()
        .navigationTitle("Account Settings")
    }
}

// MARK: - Supplemental Views

extension AccountPreferences {

    @ViewBuilder private func AccountForm() -> some View {
        VStack {
            Form {
                // Form
                TextField("Name:", text: $state.account.name)
                TextField("Host:", text: $state.account.host)
                TextField("User:", text: $state.account.user)
                SecureField("Password:", text: $state.account.password)

                // Test message result
                Group {
                    switch state.result {
                        case .untested:
                            Text(" ")
                        case .succeeded:
                            Text("Success")
                                .foregroundColor(.green)
                        case .failed(let msg):
                            Text(msg)
                                .foregroundColor(.red)
                    }
                }
                .font(.callout)

                // Control buttons

                HStack(alignment: .center, spacing: 20) {
                    Button("Test") {
                        state.test()
                    }
                    Spacer()
                    Button("Use") {
                        state.setAsDefault()
                    }
                    .disabled(state.isDefault)
                    Button("Save") {
                        state.save()
                    }
                }
            }
            Spacer()
        }
        .padding()
    }

    @ViewBuilder private func NoSelection() -> some View {
        VStack {
            if state.accounts.isEmpty {
                Text("Click [+] to create an account")
            } else {
                Text("No account selected")
            }
        }
        .font(.title3)
        .foregroundColor(.secondary)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}
