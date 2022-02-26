//
//  PostSourceEditor.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostSourceEditor: View {

    var post: WebClient.Post

    @Binding var show: Bool

    @State private var showPreview = false
    @State private var source = ""

    var body: some View {
        VStack {
            VStack(spacing: 10) {
                HStack {
                    Text(post.slugline)
                        .font(.headline)

                    Spacer()
                    DateView(date: post.datePublished, format: .dateTimeNameLong)
                        .font(.subheadline)
                }
                .lineLimit(1)

                Divider()
            }
            .padding([.horizontal, .top])

            HStack(spacing: 0) {
                    TextEditor(text: $source)
                        .lineSpacing(5)
                        .font(.body.monospaced())
                        .disableAutocorrection(true)
                        .foregroundColor(.indigo)
                        .padding(.leading, 10)
                if showPreview {
                    Divider()
                    WebPostPreview(document: source.markdownToHtml)
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
        .background(Color(nsColor: .textBackgroundColor))
        .onAppear {
            source = post.text
        }
        .onChange(of: post) { newPost in
            source = newPost.text
        }
        .toolbar {
            ToolbarItem {
                Button("Done") {
                    show.toggle()
                }
            }
            ToolbarItem {
                Button {
                    showPreview.toggle()
                } label: {
                    Image(systemName: "sidebar.right")
                }
            }

        }
    }
}

extension NSTextView {
    open override var frame: CGRect {
        didSet {
            self.isAutomaticQuoteSubstitutionEnabled = false
        }
    }
}
