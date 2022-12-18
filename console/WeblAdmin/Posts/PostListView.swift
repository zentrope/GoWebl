//
//  PostListView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostListView: View {

    @Environment(\.openURL) var openURL

    @StateObject private var state = PostListViewState()

    @State private var selectedPost: String?
    @State private var showEditor = false
    @State private var showSiteEditor = false

    var body: some View {
        NavigationView/*(spacing: 0)*/ {
            List(selection: $selectedPost) {
                ForEach(state.posts, id: \.id) { post in
                    Item(post: post)
                        .tag(post.id)
                        .padding(.vertical, 5)
                        .contextMenu {
                            ContextMenu(post: post)
                        }
                }
            }
            .listStyle(.inset(alternatesRowBackgrounds: true))
            .background(.background)
            .frame(minWidth: 250, maxWidth: .infinity)
            if let postId = selectedPost, let post = state.post(id: postId) {
                PostPreview(post: post)
                    .frame(minWidth: 400, idealWidth: 400, maxWidth: .infinity)
            } else {
                UnselectedView(message: "No Post Selected")
                    .frame(minWidth: 350, idealWidth: 350)
            }
        }
        .alert(state.error?.localizedDescription ?? "Error: check logs.", isPresented: $state.showAlert, actions: {})
        .sheet(isPresented: $showEditor, content: {
            if let postId = selectedPost {
                PostSourceEditorView(postId: postId)
                    .frame(minWidth: 800, minHeight: 600)
            } else {
                UnselectedView(message: "Selected Post Disappeared")
            }
        })
        .sheet(isPresented: $showSiteEditor, content: {
            SiteEditorView()
                .frame(width: 600)
                .fixedSize(horizontal: false, vertical: true)
        })
        .toolbar {
            Button {
                state.createNewPost()
            } label: {
                Image(systemName: "plus")
            }
            .help("Create a new post")

            Button {
                showEditor.toggle()
            } label: {
                Image(systemName: "square.and.pencil")
            }
            .disabled(selectedPost == nil)
            .help("Edit the currently selected post")

            Spacer()

            Button {
                state.refresh()
            } label: {
                Image(systemName: "arrow.clockwise")
            }
            .help("Refresh posts from the published site")

            Button {
                if let baseUrl = URL(string: DataCache.shared.site.baseUrl) {
                    openURL(baseUrl)
                }
            } label: {
                Image(systemName: "link")
            }
            .help("Visit the published site")

            Button {
                showSiteEditor.toggle()
            } label: {
                Image(systemName: "gearshape")
            }
            .help("Update metadata about the site")
        }
    }
}

// MARK: - Supplemental views

extension PostListView {

    @ViewBuilder
    private func ContextMenu(post: WebClient.Post) -> some View {
        Button("Edit") {
            selectedPost = post.id
            showEditor.toggle()
        }
        switch post.status {
            case .draft:
                Button("Publish") {
                    state.toggle(id: post.id, isPublished: true)
                }
            case .published:
                Button("Unpublish (return to draft mode)") {
                    state.toggle(id: post.id, isPublished: false)
                }
        }
        Button("Delete") {
            state.deletePost(withId: post.id)
        }
    }
}

struct Item: View {

    var post: WebClient.Post

    var body: some View {

        VStack(alignment: .leading, spacing: 2) {

            HStack(alignment: .center) {
                Text("\(post.slugline)")
                    .font(.system(.title3, design: .rounded))
                    .fontWeight(.medium)
                    .lineLimit(1)
                Spacer()
            }

            HStack(alignment: .center) {
                if post.status == .draft {
                    Text("Unpublished")
                        .font(.callout.italic())
                } else {
                    DateView(date: post.datePublished, format: .dateNameShort)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .font(.callout)
                        .help("Date published")
                }
                Spacer()
                Text("\(post.wordCount)w")
                    .font(.caption.monospacedDigit())
            }
            .foregroundColor(.secondary)
        }
//        .foregroundColor(post.status == .draft ? .secondary : .primary)
    }
}

struct StatusIcon: View {
    var status: WebClient.Post.Status
    var body: some View {
        switch status {
            case .draft:
                Image(systemName: "square.dashed")
                    .help(status.rawValue)
            case .published:
                Image(systemName: "p.square")
                    .help(status.rawValue)
        }
    }
}
