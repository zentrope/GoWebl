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
    @State private var showEditor = false

    var body: some View {
        VStack {
            if showEditor,
               let id = selectedPost,
               let post = state.post(id: id) {
                PostSourceEditor(post: post, show: $showEditor)
            } else {
                PostListNavigator(showEditor: $showEditor, selectedPost: $selectedPost, state: state)
            }
        }
        .navigationSubtitle("\(state.site.title) — \(state.name) <\(state.email)>")
    }
}


struct PostListNavigator: View {

    @Binding var showEditor: Bool
    @Binding var selectedPost: String?
    @ObservedObject var state: PostListViewState

    var body: some View {
        HStack(spacing: 0) {
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
            .frame(width: 400)

            Divider()

            if let postId = selectedPost,
               let post = state.post(id: postId) {
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
                    WebPostPreview(document: post.text.markdownToHtml)
                        .padding(.leading, 10)
                }
                .background(Color(nsColor: .textBackgroundColor))
                .frame(minWidth: 400, idealWidth: 400)
            } else {
                VStack(alignment: .center, spacing: 20) {
                    Image(systemName: "doc.text")
                        .font(.system(size: 72).weight(.thin))
                        .foregroundColor(.secondary)
                    Text("No post selected")
                        .font(.title)
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .background(Color(nsColor: .textBackgroundColor))
                .frame(minWidth: 350, idealWidth: 350)
            }
        }

        .toolbar {

            ToolbarItem {
                if selectedPost != nil {
                    Button {
                        showEditor.toggle()
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

extension PostListNavigator {

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
