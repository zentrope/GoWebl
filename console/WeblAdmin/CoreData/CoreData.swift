//
//  CoreData.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/19/22.
//

import Foundation
import CoreData
import OSLog

struct CoreData {

    static let shared = CoreData()

    let log = Logger("CoreData")

    var container: NSPersistentContainer

    init() {
        let container = NSPersistentContainer(name: "WeblData")
        container.loadPersistentStores { storeDescription, error in
            container.viewContext.automaticallyMergesChangesFromParent = true
            container.viewContext.mergePolicy = NSMergePolicy.mergeByPropertyObjectTrump
            if let error = error as NSError? {
                fatalError("unresolved error \(error), \(error.userInfo)")
            }
        }
        self.container = container
    }

    static var viewContext: NSManagedObjectContext {
        shared.container.viewContext
    }
}
