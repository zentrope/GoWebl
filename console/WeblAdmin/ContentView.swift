//
//  ContentView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

struct ContentView: View {

    enum MenuItem {
        case remote
        case local
        case stylesheets
    }

    @State private var selectedMenuItem = MenuItem.local

    var body: some View {
        NavigationSplitView {
            List(selection: $selectedMenuItem) {
                Section("Posts") {
                    Label("Remote", systemImage: "cylinder.split.1x2")
                        .tag(MenuItem.remote)
                    Label("Local", systemImage: "internaldrive")
                        .tag(MenuItem.local)
                }

                Section("Artifacts") {
                    Label("Stylesheets", systemImage: "curlybraces")
                        .tag(MenuItem.stylesheets)
                }
            }
            .navigationSplitViewColumnWidth(ideal: 150)
        } detail: {
            switch selectedMenuItem {
                case .remote:
                    PostListView()
                case .local:
                    LocalPostView()
                default:
                    UnselectedView(message: "Select area", systemName: "map")
            }
        }        
    }
}
