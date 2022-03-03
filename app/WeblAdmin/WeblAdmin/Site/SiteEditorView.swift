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
                TextField("Title:", text: $state.siteTitle)
                TextField("Description:", text: $state.siteDescription)
                TextField("Base URL:", text: $state.siteBaseURL)
                Button("Apply") {
                    state.updateSite(title: state.siteTitle, description: state.siteDescription, baseURL: state.siteBaseURL)
                }
                .disabled(!state.siteDirty)
                .controlSize(.small)
                Divider()
                TextField("Name:", text: $state.accountName)
                TextField("Email:", text: $state.accountEmail)
                Button("Apply") {

                }
                .controlSize(.small)
            }
            HStack {
                if state.working {
                    ProgressView()
                        .progressViewStyle(.circular)
                        .controlSize(.small)
                } else {
                    Text(state.siteDirty ? "Unsaved" : state.message)
                        .foregroundColor(state.siteDirty ? .red : .green)
                }
                Spacer()
                Button("Done") {
                    dismiss()
                }
                .keyboardShortcut(.cancelAction)
            }
        }
        .alert(state.error?.localizedDescription ?? "Error", isPresented: $state.showAlert, actions: {})
        .padding()
    }
}


struct Previews_SiteEditorView_Previews: PreviewProvider {
    static var previews: some View {
        SiteEditorView()
            .frame(width: 500)
    }
}
