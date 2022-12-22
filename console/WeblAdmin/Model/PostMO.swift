//
//  PostMO+CoreDataProperties.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/20/22.
//
//

import Foundation
import CoreData
import OSLog

// Custom

extension PostMO {

    private static let log = Logger("PostMO")

    typealias Transaction = NSManagedObjectContext

    static func withTransaction(handler: (Transaction) throws -> ()) {
        let tx = CoreData.newBackgroundContext()
        do {
            try handler(tx)
            log.info("COMMIT")
            try tx.commit()
        } catch {
            log.error("\(error.localizedDescription, privacy: .public)")
            log.info("ROLLBACK")
            tx.rollback()
        }
    }

    static func upsert(id: UUID, status: String, title: String, dateCreated: Date, dateUpdated: Date, datePublished: Date, text: String, context: NSManagedObjectContext = CoreData.viewContext) {
        let post = findOrNew(id: id, context: context)
        post.status = status
        post.title = title
        post.dateCreated = dateCreated
        post.dateUpdated = dateUpdated
        post.datePublished = datePublished
        post.text = text
    }

    static func find(id: UUID, context: NSManagedObjectContext = CoreData.viewContext) -> PostMO? {
        let request = PostMO.fetchRequest()
        request.predicate = NSPredicate(format: "id == %@", id as CVarArg)
        return try? context.fetch(request).first
    }

    private static func findOrNew(id: UUID, context: NSManagedObjectContext = CoreData.viewContext) -> PostMO {
        if let post = find(id: id, context: context) {
            return post
        }
        let post = PostMO(context: context)
        post.id = id
        return post
    }

    var wordCount: Int {
        text.components(separatedBy: NSCharacterSet.whitespaces).count
    }

}

// Generated

@objc(PostMO)
public class PostMO: NSManagedObject {

}

extension PostMO {

    @nonobjc public class func fetchRequest() -> NSFetchRequest<PostMO> {
        return NSFetchRequest<PostMO>(entityName: "PostMO")
    }

    @NSManaged public var id: UUID
    @NSManaged public var status: String
    @NSManaged public var title: String
    @NSManaged public var dateCreated: Date
    @NSManaged public var dateUpdated: Date
    @NSManaged public var datePublished: Date
    @NSManaged public var text: String

}

extension PostMO : Identifiable {

}
