//
//  ResourceMO.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/21/22.
//
//

import Foundation
import CoreData

// MARK: - Implementation

extension ResourceMO {

    enum Types: String {
        case css = "text/css"
        case javascript = "application/javascript"
        case json = "application/json"
        case png = "image/png"
        case svg = "image/svg+xml"
    }

}

// MARK: - Generated

@objc(ResourceMO)
public class ResourceMO: NSManagedObject {

}

extension ResourceMO {

    @nonobjc public class func fetchRequest() -> NSFetchRequest<ResourceMO> {
        return NSFetchRequest<ResourceMO>(entityName: "ResourceMO")
    }

    @NSManaged public var id: UUID
    @NSManaged public var mimeType: String
    @NSManaged public var name: String
    @NSManaged public var text: String?
    @NSManaged public var data: Data?

}

extension ResourceMO : Identifiable {

}
