//
//  AccountPreferences.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct AccountPreferences: View {

    @StateObject private var state = AccountPreferencesViewState()

    @AppStorage("WAAccountEmail") private var savedEmail = ""
    @AppStorage("WAAccountPassword") private var savedPassword = ""
    @AppStorage("WAAccountEndpoint") private var savedEndpoint = ""

    @State private var formName = ""
    @State private var formHost = ""
    @State private var formUser = ""
    @State private var formPassword = ""

    @State private var selectedAccount: UUID?

    var body: some View {
        HStack(spacing: 10) {

            VStack(alignment: .leading, spacing: 10) {
                List(selection: $selectedAccount) {
                    ForEach(state.accounts, id: \.id) { account in
                        Label(account.name, systemImage: "cylinder.split.1x2")
                            .tint(.blue)
                            .tag(account.id)
                    }
                }
                .listStyle(.inset)
                .onChange(of: selectedAccount) { accountId in
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
                        if let id = selectedAccount {
                            state.deleteAccount(id: id)
                        }
                    } label: {
                        Image(systemName: "minus")
                    }
                    .disabled(selectedAccount == nil || state.accounts.isEmpty)
                    .controlSize(.small)
                    Spacer()
                }
            }
            .frame(width: 150)
            .alert(state.error?.localizedDescription ?? "Error", isPresented: $state.showAlert) {}

            if selectedAccount == nil {
                NoSelection()
            } else {
                AccountForm()
            }
        }
        .padding()
        .frame(width: 600, height: 200)
        .fixedSize()
        .onAppear {
            selectedAccount = state.accounts.first?.id
        }
    }

    @ViewBuilder
    private func AccountForm() -> some View {
        VStack {
            Form {
                TextField("Name:", text: $state.account.name)
                TextField("Host:", text: $state.account.host)
                TextField("User:", text: $state.account.user)
                SecureField("Password:", text: $state.account.password)
            }

            Spacer()

            HStack(alignment: .center, spacing: 20) {
                Button("Test") {
                    state.test()
                    //state.testConnection(email: email, password: password)
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
                Button("Use") {
                    savedEmail = state.account.user
                    savedPassword = state.account.password
                    savedEndpoint = state.account.host
                }
                .disabled(savedEmail == state.account.user && savedPassword == state.account.password && savedEndpoint == state.account.host)
                Button("Save") {
                    state.save()
                }
            }
            .padding(.top)
        }
        .padding()
    }

    @ViewBuilder
    private func NoSelection() -> some View {
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
