//
//  ContentView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct ContentView: View {

    enum TestResult {
        case untested
        case succeeded
        case failed(String)
    }

    @State private var result: TestResult = .untested

    var body: some View {
        PostListView()
            .frame(minWidth: 800, minHeight: 600)
        
    }
}
