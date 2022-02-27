//
//  UnselectedView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import SwiftUI

struct UnselectedView: View {
    var message: String
    var systemName = "doc.text"

    var body: some View {
        VStack(alignment: .center, spacing: 20) {
            Image(systemName: systemName)
                .font(.system(size: 72).weight(.thin))
                .foregroundColor(.secondary)
            Text(message)
                .font(.title)
                .foregroundColor(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(Color.textBackgroundColor)
    }
}
