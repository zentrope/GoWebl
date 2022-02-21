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
        //HStack(spacing: 0) {
            List(selection: $selectedPost) {

                ForEach(state.posts, id: \.id) { post in
                    HStack {
                        Text("\(post.wordCount)")
                            .frame(width: 40, alignment: .trailing)
                            .foregroundColor(.secondary)
                        Text("\(post.slugline)")
                            .lineLimit(1)
                        Spacer()
                        Group {
                        DateView(date: post.dateCreated, format: .dateNameShort)
                            .frame(width: 80, alignment: .leading)
                            .font(.callout)
                            .foregroundColor(.secondary)
                        }
                    }
                }
            }
            .listStyle(.inset(alternatesRowBackgrounds: true))
            .frame(minWidth: 350, idealWidth: 350)

          //  Divider()

            if selectedPost != nil {
                ScrollView {
                    VStack {
                        TextEditor(text: $bodyText)
                            .disableAutocorrection(false)                            
                            .font(.body.monospaced())
                            .frame(maxWidth: .infinity, maxHeight: .infinity)
                        Spacer()
                    }
                    .padding()
                }
                .background(Color(nsColor: .textBackgroundColor))
                .frame(minWidth: 350)
            }
        }
        .onChange(of: selectedPost) { postId in
            let post = state.post(id: postId)
            bodyText = post?.text ?? ""
        }
        .toolbar {
            ToolbarItem {
                Button {
                    print("Not implemented.")
                } label: {
                    Image(systemName: "gear")
                }

            }
        }
    }
}
