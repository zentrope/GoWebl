//
//  PostSourceEditor.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

class PostSourceEditorViewState: NSObject, ObservableObject {

    @Published var showAlert = false
    @Published var error: Error?

    func update(post: WebClient.Post, newText: String) {
        Task {
            do {
                let client = WebClient()
                let updatedPost = try await client.updatePost(uuid: post.id, slugline: post.slugline, text: newText, datePublished: post.datePublished)
                print("Updated: \(updatedPost.id)")
            } catch (let err) {
                showAlert(error: err)
            }
        }
    }

    func showAlert(error: Error) {
        self.showAlert = true
        self.error = error
    }
}

struct PostSourceEditor: View {
    @Environment(\.dismiss) private var dismiss

    var post: WebClient.Post

    @StateObject private var state = PostSourceEditorViewState()
    @State private var showPreview = false
    @State private var source = ""

    var body: some View {
        VStack(spacing: 0) {
            VStack(spacing: 10) {
                HStack {
                    Text(post.slugline)
                        .font(.headline)

                    Spacer()
                    DateView(date: post.datePublished, format: .dateTimeNameLong)
                        .font(.subheadline)
                }
                .lineLimit(1)
            }
            .padding(10)
            .overlay(Divider(), alignment: .bottom)

            HStack(spacing: 0) {
                TextEditor(text: $source)
                    .lineSpacing(5)
                    .font(.body.monospaced())
                    .foregroundColor(.indigo)
                    .padding(.leading, 10)

                if showPreview {
                    Divider()
                    WebPostPreview(document: source.markdownToHtml)
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)

            HStack {
                Button("Save") {
                    state.update(post: post, newText: source)
                }
                .disabled(post.text == source)
                Button("Cancel") {
                    dismiss()
                }
                .keyboardShortcut(.cancelAction)
                Spacer()
                Button(showPreview ? "Hide Preview" : "Show Preview") {
                    showPreview.toggle()
                }
            }
            .controlSize(.small)
            .padding(10)
            .overlay(Divider(), alignment: .top)
        }
        .background(Color.textBackgroundColor)
        .onAppear {
            source = post.text
        }
    }
}

extension NSTextView {
    open override var frame: CGRect {
        didSet {
            self.isAutomaticQuoteSubstitutionEnabled = false
            self.isContinuousSpellCheckingEnabled = true
        }
    }
}
