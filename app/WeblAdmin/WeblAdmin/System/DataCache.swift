//
//  DataCache.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import Foundation

@MainActor
final class DataCache {

    static let DataCacheDidChange = NSNotification.Name("DataCacheDidChange")
    static let shared = DataCache()

    var uuid = ""
    var name = ""
    var email = ""
    var site = WebClient.Site(baseUrl: "", title: "", description: "")
    var type = ""
    var posts = [WebClient.Post]()

    private var cachedData = [String:WebClient.Post]()

    func replaceAll(viewer: WebClient.Viewer) {
        self.uuid = viewer.id;
        self.name = viewer.name;
        self.email = viewer.email
        self.site = viewer.site
        self.type = viewer.type
        replaceAll(posts: viewer.posts)
        notify()
    }

    private func replaceAll(posts: [WebClient.Post]) {
        cachedData.removeAll()
        for post in posts {
            cachedData[post.id] = post
        }
        self.posts = posts.sorted(by: { $0.dateCreated > $1.dateCreated })
    }

    subscript(index: String) -> WebClient.Post? {
        get {
            return cachedData[index]
        }

        set(newValue) {
            cachedData[index] = newValue
            self.posts = cachedData.values.sorted(by: { $0.dateCreated > $1.dateCreated })
            notify()
        }
    }

    private func notify() {
        NotificationCenter.default.post(name: DataCache.DataCacheDidChange, object: self)
    }
}
