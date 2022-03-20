//
//  AccountMO+CoreDataProperties.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/19/22.
//
//

import Foundation
import CoreData


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
