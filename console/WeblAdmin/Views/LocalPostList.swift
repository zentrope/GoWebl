//
//  LocalPostList.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/20/22.
//

import SwiftUI

struct LocalPostList: View {

    @FetchRequest(sortDescriptors: [SortDescriptor(\.datePublished, order: SortOrder.reverse)])
    private var posts: FetchedResults<PostMO>

    @Binding var selection: PostMO.ID?

    var body: some View {
        List(selection: $selection) {
            ForEach(posts, id: \.id) { post in
                VStack(alignment: .leading, spacing: 2) {
                    Text(post.title)
                        .font(.system(.title3, design: .rounded))
                        .fontWeight(.medium)
                        .lineLimit(1)
                    HStack(alignment: .center) {
                        if post.status == "draft" {
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
                .padding(.vertical, 5)
            }
        }
        .listStyle(.inset(alternatesRowBackgrounds: true))
        .frame(minWidth: 250, maxWidth: .infinity)
//        .toolbar {
//            Button("Sync") {
//                importRemotePosts()
//            }
//        }
    }

    private func importRemotePosts() {
        Task {
            let client = WebClient()
            do {
                let data = try await client.viewerData()
                let posts = data.posts

                PostMO.withTransaction { tx in
                    for post in posts {
                        guard let id = UUID(uuidString: post.id) else {
                            continue
                        }
                        PostMO.upsert(id: id, status: post.status.rawValue, title: post.slugline, dateCreated: post.dateCreated, dateUpdated: post.dateUpdated, datePublished: post.datePublished, text: post.text, context: tx)
                        print(post.slugline)
                    }
                }

            } catch {
                print("ERROR: \(error)")
            }
        }
        
    }
}
