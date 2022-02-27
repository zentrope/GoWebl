//
//  PostSourceEditorView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostSourceEditorView: View {

    var postId: String

    @Environment(\.dismiss) private var dismiss

    @StateObject private var state = PostSourceEditorViewState()

    @State private var showPreview = false
    @State private var source = ""

    var body: some View {
        VStack(spacing: 0) {
            VStack(spacing: 10) {
                HStack {
                    Text(state.post.slugline)
                        .font(.headline)
                    Spacer()
                    DateView(date: state.post.datePublished, format: .dateTimeNameLong)
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
                    state.update(post: state.post, newText: source)
                }
                .disabled(state.post.text == source)
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
            state.setPost(toPostWithId: postId)        
            source = state.post.text
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
