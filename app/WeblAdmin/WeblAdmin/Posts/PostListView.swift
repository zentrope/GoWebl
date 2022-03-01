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
        HStack(spacing: 0) {
            List(selection: $selectedPost) {
                ForEach(state.posts, id: \.id) { post in
                    Item(post: post)
                        .tag(post.id)
                        .padding(.vertical, 2)
                        .contextMenu {
                            ContextMenu(post: post)
                        }
                }
            }
            .listStyle(.inset(alternatesRowBackgrounds: true))
            .frame(width: 400)

            Divider()

            if let postId = selectedPost, let post = state.post(id: postId) {
                PostPreview(post: post)
                    .frame(minWidth: 400, idealWidth: 400)
            } else {
                UnselectedView(message: "No Post Selected")
                    .frame(minWidth: 350, idealWidth: 350)
            }
        }
        .navigationSubtitle("\(state.site.title) — \(state.name) <\(state.email)>")
        .sheet(isPresented: $showEditor, content: {
            if let postId = selectedPost {
                PostSourceEditorView(postId: postId)
                    .frame(minWidth: 800, minHeight: 600)
            } else {
                UnselectedView(message: "Selected Post Disappeared")
            }
        })
        .toolbar {
            ToolbarItem {
                Button {
                    Task { self.selectedPost = await state.newPost() }
                } label: {
                    Image(systemName: "plus")
                }
            }
            ToolbarItem {
                Button {
                    showEditor.toggle()
                } label: {
                    Image(systemName: "square.and.pencil")
                }
                .disabled(selectedPost == nil)
            }
            ToolbarItem(placement: .navigation) {
                Button {
                    state.refresh()
                } label: {
                    Image(systemName: "arrow.clockwise")
                }
            }
        }
    }
}

// MARK: - Supplemental views

extension PostListView {

    @ViewBuilder
    private func ContextMenu(post: WebClient.Post) -> some View {
        Button("Edit this post") {
            selectedPost = post.id
            showEditor.toggle()
        }
        switch post.status {
            case .draft:
                Button("Publish this post") {
                    state.toggle(id: post.id, isPublished: true)
                }
            case .published:
                Button("Unpublish this post (return to draft mode)") {
                    state.toggle(id: post.id, isPublished: false)
                }
        }
    }

    @ViewBuilder
    private func Item(post: WebClient.Post) -> some View {
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
    }

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
