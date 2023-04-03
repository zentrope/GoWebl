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
        case resources
        case templates
    }

    @State private var selectedMenuItem = MenuItem.local

    var body: some View {
        NavigationSplitView {
            List(selection: $selectedMenuItem) {
                Section("Posts") {
                    Label("Local", systemImage: "internaldrive")
                        .tag(MenuItem.local)
                    Label("Remote", systemImage: "cylinder.split.1x2")
                        .tag(MenuItem.remote)
                }

                Section("Site") {
                    Label("Resources", systemImage: "curlybraces")
                        .tag(MenuItem.resources)
                    Label("Templates", systemImage: "doc.on.doc")
                        .tag(MenuItem.templates)
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
