//
//  String.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/25/22.
//

import Foundation

extension String {
    var markdownToHtml: String {
        HTMLEncoder(self).html()
    }

    var escapeHtml: String {
        self
            .replacingOccurrences(of: "<", with: "&lt;")
            .replacingOccurrences(of: ">", with: "&gt;")
    }
}
