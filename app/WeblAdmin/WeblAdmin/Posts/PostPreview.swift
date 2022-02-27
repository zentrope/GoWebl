//
//  PostPreview.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import SwiftUI

struct PostPreview: View {
    var post: WebClient.Post

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
            WebPostPreview(document: post.text.markdownToHtml)
                .padding(.leading, 10)
        }
        .background(Color.textBackgroundColor)
    }
}
