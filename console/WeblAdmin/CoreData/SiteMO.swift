//
//  SiteMO+CoreDataProperties.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/18/22.
//
//

import Foundation
import CoreData


// Custom

extension SiteMO {

    public static func retrieve() -> SiteMO {
        let context = CoreData.shared.container.viewContext
        let request = SiteMO.fetchRequest()
        if let site = try? context.fetch(request).first {
            return site
        }

        let mo = SiteMO(context: context)
        mo.id = UUID()
        mo.author = ""
        mo.title = "New Site"
        mo.subtitle = "A new blog for a new world"
        mo.baseUrl = ""
        return mo
    }

    public func save() throws {
        let context = CoreData.shared.container.viewContext
        if context.hasChanges {
            try context.save()
        }
    }

    public func rollback() {
        let context = CoreData.shared.container.viewContext
        context.rollback()
    }
}

// Generated

@objc(SiteMO)
public class SiteMO: NSManagedObject {

}

extension SiteMO {

    @nonobjc public class func fetchRequest() -> NSFetchRequest<SiteMO> {
        return NSFetchRequest<SiteMO>(entityName: "SiteMO")
    }

    @NSManaged public var id: UUID?
    @NSManaged public var title: String
    @NSManaged public var subtitle: String
    @NSManaged public var baseUrl: String
    @NSManaged public var author: String

}

extension SiteMO : Identifiable {

}
