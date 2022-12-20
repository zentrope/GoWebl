//
//  SiteSettingsView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/19/22.
//

import SwiftUI

struct SiteSettingsView: View {

    @ObservedObject private var site = SiteMO.retrieve()

    var body: some View {
        VStack(spacing: 20) {
            Form {
                TextField("Title:", text: $site.title)
                TextField("Subtitle:", text: $site.subtitle)
                TextField("Author:", text: $site.author)
                TextField("BaseURL:", text: $site.baseUrl, prompt: Text("https://example.com"))
            }

            HStack {
                Spacer()
                Button {
                    site.rollback()
                } label: {
                    Text("Cancel")
                        .frame(width: 60)
                }
                .disabled(!site.hasPersistentChangedValues)
                .keyboardShortcut(.cancelAction)

                Button {
                    save()
                } label: {
                    Text("Save")
                        .frame(width: 60)
                }
                .disabled(!site.hasPersistentChangedValues)
                .keyboardShortcut(.return)
            }
        }
        .padding()
        .frame(width: 400)
        .fixedSize(horizontal: true, vertical: true)
        .navigationTitle("Site Settings")
    }

    private func save() {
        do {
            site.author = site.author.trimmingCharacters(in: .whitespacesAndNewlines)
            site.title = site.title.trimmingCharacters(in: .whitespacesAndNewlines)
            site.subtitle = site.subtitle.trimmingCharacters(in: .whitespacesAndNewlines)
            site.baseUrl = site.baseUrl.trimmingCharacters(in: .whitespacesAndNewlines)
            try site.save()
        } catch {
            print("ERROR SAVE: \(error)")
        }
    }
}
