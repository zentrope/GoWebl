//
//  SiteEditorView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/2/22.
//

import SwiftUI

struct SiteEditorView: View {

    @Environment(\.dismiss) private var dismiss

    @StateObject private var state = SiteEditorViewState()

    var body: some View {
        VStack(spacing: 20) {
            HStack {
                Image(systemName: "network")
                    .foregroundColor(.accentColor)
                Text("Update site details")
                Spacer()
            }
            .font(.title2)
            Form {
                Group {
                    TextField("Title:", text: $state.siteTitle)
                    TextField("Description:", text: $state.siteDescription)
                    TextField("Base URL:", text: $state.siteBaseURL)
                    HStack(alignment: .center) {
                        Button("Apply") {
                            state.updateSite(title: state.siteTitle, description: state.siteDescription, baseURL: state.siteBaseURL)
                        }
                        .disabled(!state.siteDirty)
                        .controlSize(.small)
                        StatusMessage(isWorking: state.savingSite, isDirty: state.siteDirty)
                    }
                    .padding(.vertical, 4)
                }

                Group {
                    Divider()
                    TextField("Name:", text: $state.accountName)
                    TextField("Email:", text: $state.accountEmail)
                    HStack {
                        Button("Apply") {
                            state.updateAccount(name: state.accountName, email: state.accountEmail)
                        }
                        .disabled(!state.accountDirty)
                        .controlSize(.small)
                        StatusMessage(isWorking: state.savingAccount, isDirty: state.accountDirty)
                    }
                    .padding(.vertical, 4)
                }

                Group {
                    Divider()
                    SecureField("New password:", text: $state.newPassword)
                    SecureField("Confirm password:", text: $state.confirmPassword)
                    HStack {
                        Button("Save") {
                            state.updatePassword(toNewPassword: state.newPassword)
                        }
                        .disabled(state.newPassword != state.confirmPassword)
                        .disabled(state.newPassword.count < 8)
                        .controlSize(.small)
                    }
                    .padding(.vertical, 4)
                }
            }
            HStack {
                Spacer()
                Button {
                    //
                } label: {
                    Text("Cancel").frame(width: 60)
                }
                .keyboardShortcut(.cancelAction)
                Button {
                    dismiss()
                } label: {
                    Text("OK").frame(width: 60)
                }
                .keyboardShortcut(.defaultAction)
            }
        }
        .alert(state.error?.localizedDescription ?? "Error", isPresented: $state.showAlert, actions: {})
        .padding()
    }

    @ViewBuilder
    private func StatusMessage(isWorking: Bool, isDirty: Bool) -> some View {
        if isWorking {
            ProgressView()
                .progressViewStyle(.circular)
                .controlSize(.mini)
        } else {
            Image(systemName: isDirty ? "x.circle.fill" : "checkmark.circle.fill")
                .foregroundColor(isDirty ? .red : .green)
                .font(.callout)
        }
    }
}


struct Previews_SiteEditorView_Previews: PreviewProvider {
    static var previews: some View {
        SiteEditorView()
            .frame(width: 500)
    }
}
