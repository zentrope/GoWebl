//
//  AccountSettingsView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/19/22.
//

import SwiftUI

struct AccountSettingsView: View {

    @AppStorage("WMAccountID") private var savedAccountId = ""
    @AppStorage("WAAccountEmail") private var savedEmail = ""
    @AppStorage("WAAccountPassword") private var savedPassword = ""
    @AppStorage("WAAccountEndpoint") private var savedEndpoint = ""

    @State private var selection: AccountMO.ID = UUID()

    @State private var id = AccountMO.ID()
    @State private var name = ""
    @State private var host = ""
    @State private var user = ""
    @State private var pass = ""

    @State private var testResult = " "

    var body: some View {
        HStack(spacing: 20) {
            List(selection: $selection) {
                ForEach(AccountMO.all(), id: \.id) { account in
                    HStack {
                        Label(account.name, systemImage: "person.fill")
                        Spacer()
                        if account.id.uuidString == savedAccountId {
                            Image(systemName: "hand.thumbsup.fill")
                        }
                    }
                    .contextMenu {
                        Button("Set Default") {
                            changeDefault()
                        }
                        .disabled(account.id.uuidString == savedAccountId)
                    }
                }
            }
            .clipShape(RoundedRectangle(cornerRadius: 7))

            VStack(spacing: 20) {

                Form(content: {
                    TextField("Name:", text: $name)
                    TextField("Host:", text: $host)
                    TextField("Username:", text: $user)
                    SecureField("Password:", text: $pass)
                })

                Divider()

                    .overlay(alignment: .leading, content: { Text("\(testResult)")
                            .lineLimit(1)
                            .frame(alignment: .leading)
                            .font(.caption)
                            .offset(x: 4, y: 8)
                            .foregroundColor(.accentColor)
                    })


                HStack {
                    Button {
                        test()
                    } label: {
                        Text("Test")
                            .frame(width: 60)
                    }

                    Spacer()
                    Button {
                    } label: {
                        Text("Delete")
                            .frame(width: 60)
                    }
                    .disabled(true)
                    Button {
                    } label: {
                        Text("Cancel")
                            .frame(width: 60)
                    }
                    .disabled(true)
                    Button {
                    } label: {
                        Text("Save")
                            .frame(width: 60)
                    }
                    .disabled(true)
                }
            }
            .frame(width: 350)

        }
        .padding()
        .frame(width: 600)
        .fixedSize(horizontal: true, vertical: true)

        .onChange(of: selection) { accountId in
            resetState()
            if let account = try? AccountMO.find(id: accountId) {
                populateState(account)
            }
        }
        .onAppear {
            selection = AccountMO.all().first?.id ?? UUID()
        }
    }

    private func test() {
        Task {
            do {
                let client = WebClient()
                let rc = try await client.test(user: user, pass: pass, host: host)
                testResult = rc ? "Success" : "Failed"
            } catch {
                testResult = error.localizedDescription
            }

            try? await Task.sleep(for: .seconds(5))
            testResult = ""
        }
    }

    private func changeDefault() {
        savedAccountId = id.uuidString
        savedEmail = user
        savedPassword = pass
        savedEndpoint = host
    }

    private func populateState(_ account: AccountMO) {
        id = account.id
        name = account.name
        host = account.host
        user = account.user
        pass = account.password
    }

    private func resetState() {
        id = UUID()
        name = ""
        host = ""
        user = ""
        pass = ""
    }
}
