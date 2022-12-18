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
                        .font(.system(.title, design: .rounded))
                        .fontWeight(.semibold)

                    Spacer()
                    DateView(date: post.datePublished, format: .dateTimeNameLong)
                        .font(.body)
                }
                .lineLimit(1)
            }
            .padding([.horizontal, .top])
            WebPostPreview(document: post.text.markdownToHtml)
                .padding(.leading, 10)
        }
        .background(Color.textBackgroundColor)
    }
}
