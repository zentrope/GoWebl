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
            //.overlay(Divider(), alignment: .bottom)

            ScrollView {
                VStack {

                    Text(post.text)
                        .font(.body.monospaced())
                        .foregroundColor(.indigo)
                        .lineSpacing(5)
                    Spacer()
                }
            }

        }
        .padding()
        .background(Color(nsColor: .textBackgroundColor))
    }
}
