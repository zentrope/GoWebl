//
//  LocalPostView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/21/22.
//

import SwiftUI

struct LocalPostView: View {

    @State private var selectedPost: PostMO.ID?

    var body: some View {
        HSplitView {
            LocalPostList(selection: $selectedPost)

            if let selectedPost, let post = PostMO.find(id: selectedPost) {
                LocalPostDetail(post: post)
                    .frame(minWidth: 200)
                    .layoutPriority(1)
            } else {
                UnselectedView(message: "Select a post")
                    .frame(minWidth: 200)
                    .layoutPriority(1)
            }
        }
    }
}
