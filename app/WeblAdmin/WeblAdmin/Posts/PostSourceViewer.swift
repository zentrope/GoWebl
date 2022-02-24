//
//  PostSourceViewer.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct PostSourceViewer: View {

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

            VStack {
                WebPostPreview(document: "<p>\(post.text)</p>")
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .padding(.leading, 10)
        }
        .background(Color(nsColor: .textBackgroundColor))
    }
}
