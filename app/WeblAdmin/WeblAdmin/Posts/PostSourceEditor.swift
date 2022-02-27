//
//  PostSourceEditor.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostSourceEditor: View {
    @Environment(\.dismiss) private var dismiss

    var post: WebClient.Post

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

                }
                .disabled(true)
                Button("Cancel") {
                    dismiss()
                }
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
