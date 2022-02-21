//
//  PostListView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostListView: View {

    @StateObject private var state = PostListViewState()

    @State private var selectedPost: String?
    @State private var bodyText = ""

    var body: some View {
        HSplitView {
            List(selection: $selectedPost) {

                ForEach(state.posts, id: \.id) { post in
                    HStack(alignment: .center) {
                        StatusIcon(post.status)
                            .font(.callout)
                            .frame(width: 15)
                        Text("\(post.wordCount)")
                            .frame(width: 40, alignment: .trailing)
                            .help("Word count")
                        Text("\(post.slugline)")
                            .lineLimit(1)
                        Spacer()
                        DateView(date: post.datePublished, format: .dateDense)
                            .frame(width: 130, alignment: .leading)
                            .font(.caption)
                            .foregroundColor(.secondary)
                            .help("Date published")
                    }
                    .padding(.vertical, 2)
                }
            }
            .listStyle(.inset(alternatesRowBackgrounds: true))
            .frame(minWidth: 350, idealWidth: 350)

            if let postId = selectedPost,
               let post = state.post(id: postId) {
                PostSourceViewer(post: post)
                    .frame(minWidth: 350)
            }
        }
        .navigationSubtitle("\(state.site.title) — \(state.name) <\(state.email)>")
        .onChange(of: selectedPost) { postId in
            let post = state.post(id: postId)
            bodyText = post?.text ?? ""
        }
        .toolbar {

            ToolbarItem {
                if selectedPost != nil {
                    Button {
                        print("Not implemented.")
                    } label: {
                        Image(systemName: "square.and.pencil")
                    }
                }
            }

            ToolbarItem {
                Button {
                    print("Not implemented.")
                } label: {
                    Image(systemName: "gearshape")
                }

            }
        }
    }
}

// MARK: - Supplemental views

extension PostListView {

    @ViewBuilder
    private func StatusIcon(_ status: WebClient.Post.Status) -> some View {
        switch status {
            case .draft:
                Image(systemName: "icloud.slash")
                    .foregroundColor(.secondary)
                    .help(status.rawValue)
            case .published:
                Image(systemName: "icloud.fill")
                    .foregroundColor(.mint)
                    .help(status.rawValue)

        }
    }
}
