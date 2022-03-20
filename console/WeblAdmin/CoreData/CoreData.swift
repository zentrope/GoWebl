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
}

// MARK: - Mutators

extension CoreData {

    func createNewAccount() async throws {
        let context = newContext()
        try await context.perform {
            let account = AccountMO(context: context)
            account.id = UUID()
            account.name = "New Account"
            account.user = ""
            account.password = ""
            account.host = ""
            try context.commit()
        }
    }

    func deleteAccount(id: UUID) async throws {
        let context = newContext()
        try await context.perform {
            let request = AccountMO.fetchRequest()
            request.predicate = NSPredicate(format: "id = %@", id as CVarArg)
            request.includesPropertyValues = false
            let accounts = try context.fetch(request)
            if let account = accounts.first {
                context.delete(account)
            }
            try context.commit()
        }
    }

    func updateAccount(id: UUID, name: String, user: String, pass: String, host: String) async throws {
        let context = newContext()
        try await context.perform {
            let account = try findAccount(id: id, context: context)
            account.name = name
            account.user = user
            account.password = pass
            account.host = host
            try context.commit()
        }
    }

    private func newContext() -> NSManagedObjectContext {
        let context = container.newBackgroundContext()
        context.mergePolicy = NSMergePolicy.mergeByPropertyObjectTrump
        context.automaticallyMergesChangesFromParent = true
        return context
    }
}

// MARK: - Queries

extension CoreData {

    func findAccount(id: UUID, context: NSManagedObjectContext = CoreData.shared.container.viewContext) throws -> AccountMO {
        let request = AccountMO.fetchRequest()
        request.predicate = NSPredicate(format: "id = %@", id as CVarArg)
        request.fetchLimit = 1
        guard let account = try context.fetch(request).first else {
            throw CDError.FetchFailed
        }
        return account
    }
    
}

// MARK: - Errors

extension CoreData {

    enum CDError: Error, LocalizedError {
        case FetchFailed

        var errorDescription: String? {
            switch self {
                case .FetchFailed: return "unable to complete fetch"
            }
        }
    }
}

// MARK: - Other Extensions

extension NSManagedObjectContext {
    func commit() throws {
        if hasChanges {
            try save()
        }
    }
}
