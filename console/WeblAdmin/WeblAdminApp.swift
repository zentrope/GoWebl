//
//  WeblAdminApp.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI

@main
struct WeblAdminApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, CoreData.viewContext)
                .frame(minWidth: 500, minHeight: 300)
        }
        .windowToolbarStyle(.unified(showsTitle: false))
        .commands {
            SidebarCommands()
        }

        Settings {
            SettingsView()
                .environment(\.managedObjectContext, CoreData.viewContext)
        }
    }
}
