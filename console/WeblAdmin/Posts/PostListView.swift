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
            ToolbarItem(placement: .navigation) {
                Button {
                    showSiteEditor.toggle()
                } label: {
                    Image(systemName: "gearshape")
                }
                .help("Update metadata about the site")
            }

            ToolbarItem {
                Button {
                    state.createNewPost()
                } label: {
                    Image(systemName: "plus")
                }
                .help("Create a new post")
            }

            ToolbarItem {
                Button {
                    showEditor.toggle()
                } label: {
                    Image(systemName: "square.and.pencil")
                }
                .disabled(selectedPost == nil)
                .help("Edit the currently selected post")
            }

            ToolbarItem {
                Button {
                    state.refresh()
                } label: {
                    Image(systemName: "arrow.clockwise")
                }
                .help("Refresh posts from the published site")
            }

            ToolbarItem {
                Button {
                    if let baseUrl = URL(string: DataCache.shared.site.baseUrl) {
                        openURL(baseUrl)
                    }
                } label: {
                    Image(systemName: "link")
                }
                .help("Visit the published site")
            }
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

    @ViewBuilder
    private func Item(post: WebClient.Post) -> some View {
        HStack(alignment: .center) {
            StatusIcon(post.status)
                .font(.callout)
                .frame(width: 15)
            Text("\(post.slugline)")
                .lineLimit(1)
            Spacer()
            DateView(date: post.datePublished, format: .dateDense)
                .frame(width: 110, alignment: .leading)
                .font(.caption)
                .foregroundColor(.secondary)
                .help("Date published")
            Text("\(post.wordCount)")
                .font(.caption.monospacedDigit())
                .frame(width: 40, alignment: .trailing)
                .help("Word count")
        }
        .opacity(post.status == .draft ? 0.5 : 1)
    }

    @ViewBuilder
    private func StatusIcon(_ status: WebClient.Post.Status) -> some View {
        switch status {
            case .draft:
                Image(systemName: "square.dashed")
                    .help(status.rawValue)
            case .published:
                Image(systemName: "p.square")
                    .foregroundColor(.secondary)
                    .help(status.rawValue)
        }
    }
}
