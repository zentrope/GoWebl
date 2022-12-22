//
//  LocalPostDetail.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/20/22.
//

import SwiftUI

struct LocalPostDetail: View {

    var post: PostMO

    var body: some View {
        VStack(alignment: .leading, spacing: 2) {
            VStack(alignment: .leading, spacing: 2) {
                Text("Title").bold()
                Text(post.title)
                Text("Published").bold()
                DateView(date: post.datePublished, format: .dateTimeNameLong)
                Text("Body").bold()
            }
            .padding()

            ScrollView {
                Text(post.text)
                    .lineSpacing(1.5)
                    .monospaced()
                    .padding()
                    .frame(maxWidth: .infinity, alignment: .leading)
            }
        }
        .background(.background)
    }
}
