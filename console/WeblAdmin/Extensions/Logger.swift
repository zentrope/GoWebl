//
//  Logger.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import Foundation
import OSLog

extension Logger {
    init(_ category: String) {
        self.init(subsystem: "com.zentrope.WeblAdmin", category: category)
    }
}

