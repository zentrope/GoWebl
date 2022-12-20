//
//  AccountMO+CoreDataProperties.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/19/22.
//
//

import Foundation
import CoreData

// Custom

extension AccountMO {

    public static func new() -> AccountMO {
        let context = CoreData.viewContext
        let account = AccountMO(context: context)
        account.name = "New Account"
        account.id = UUID()
        account.host = ""
        account.user = ""
        account.password = ""
        return account
    }

    public static func all() -> [AccountMO] {
        let context = CoreData.viewContext
        let request = AccountMO.fetchRequest()
        request.sortDescriptors = [
            NSSortDescriptor(key: "name", ascending: true)
        ]
        let results = try? context.fetch(request)
        return results ?? []
    }

    public static func find(id: AccountMO.ID) throws -> AccountMO? {
        let context = CoreData.viewContext
        let request = AccountMO.fetchRequest()
        request.predicate = NSPredicate(format: "id == %@", id as CVarArg)
        let result = try context.fetch(request).first
        return result
    }
}

// Generated

@objc(AccountMO)
public class AccountMO: NSManagedObject {

}

extension AccountMO {

    @nonobjc public class func fetchRequest() -> NSFetchRequest<AccountMO> {
        return NSFetchRequest<AccountMO>(entityName: "AccountMO")
    }

    @NSManaged public var host: String
    @NSManaged public var id: UUID
    @NSManaged public var name: String
    @NSManaged public var password: String
    @NSManaged public var user: String

}

extension AccountMO : Identifiable {

}
