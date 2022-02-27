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
    @State private var slugline = ""
    @State private var datePublished = Date.distantPast
    @State private var source = ""

    var body: some View {
        VStack(spacing: 0) {
            HStack(spacing: 20) {
                TextField("Title:", text: $slugline)
                    .textFieldStyle(.roundedBorder)
                Spacer()
                DatePicker("Published:", selection: $datePublished, displayedComponents: [.date, .hourAndMinute])
            }
            .lineLimit(1)
            .padding(10)
            .background(Color.windowBackgroundColor)
            .overlay(Divider(), alignment: .bottom)

            HStack(spacing: 3) {
                if isDirty {
                    Text("Changes have not been saved")
                        .foregroundColor(.red)
                } else {
                    Text("Saved")
                        .foregroundColor(.green)
                }
                Spacer()
                Text("\(source.words)")
                    .bold()
                Text("words ")
                Text("\(source.count)")
                    .bold()
                Text("characters ")
                Button {
                    showPreview.toggle()
                } label: {
                    Image(systemName: showPreview ? "rectangle" : "rectangle.trailinghalf.filled")
                }
                .buttonStyle(.borderless)
                .help(showPreview ? "Show preview" : "Hide preview")

            }
            .font(.callout)
            .padding([.horizontal], 10)
            .padding(.vertical, 5)
            .overlay(Divider(), alignment: .bottom)

            HStack(spacing: 0) {
                TextEditor(text: $source)
                    .lineSpacing(5)
                    .font(.body.monospaced())
                    .padding(.leading, 10)

                if showPreview {
                    Divider()
                    WebPostPreview(document: source.markdownToHtml)
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Color.textBackgroundColor)

            HStack {
                Spacer()
                Button("Cancel") {
                    dismiss()
                }
                .keyboardShortcut(.cancelAction)
                Button("Save") {
                    state.update(post: state.post, title: slugline, source: source, published: datePublished)
                }
                .disabled(!isDirty)
            }
            .controlSize(.small)
            .padding(10)
            .overlay(Divider(), alignment: .top)
            .background(Color.windowBackgroundColor)
        }

        .onAppear {
            state.setPost(toPostWithId: postId)        
            source = state.post.text
            slugline = state.post.slugline
            datePublished = state.post.datePublished//.customEditFormat
        }
    }
}

// MARK: - Helpers

extension PostSourceEditorView {

    var isDirty: Bool {
        return state.post.text != source || state.post.slugline != slugline || state.post.datePublished != datePublished
    }
}

// MARK: - Convenience Extensions

extension NSTextView {
    open override var frame: CGRect {
        didSet {
            self.isAutomaticQuoteSubstitutionEnabled = false
            self.isContinuousSpellCheckingEnabled = true
        }
    }
}
