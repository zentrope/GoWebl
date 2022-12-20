//
//  SettingsView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/18/22.
//

import SwiftUI

struct SettingsView: View {

    enum MenuItem: Hashable {
        case site
        case account
    }

    @State private var selection = MenuItem.site

    var body: some View {
        TabView {
            SiteSettingsView()
                .tabItem {
                    Label("Site", systemImage: "doc.on.doc")
                }
            AccountPreferences()
                .tabItem {
                    Label("Accounts", systemImage: "person.2")
                }
        }
    }
}
