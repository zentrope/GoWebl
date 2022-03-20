//
//  DataCache.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import Foundation

extension Notification.Name {
    static let WeblDataCacheDidChange = NSNotification.Name("WeblDataCacheDidChange")
}

@MainActor
final class DataCache {

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
        sendNotification()
    }

    func replace(site: WebClient.Site) {
        self.site = site
    }

    /// Clear out all data: does _not_ send a notification by default.
    /// - Parameter notify: Send a notification after everything is cleared.
    func clear(notify: Bool = false) {
        self.uuid = ""
        self.name = ""
        self.email = ""
        self.site = WebClient.Site(baseUrl: "", title: "", description: "")
        self.type = ""
        posts.removeAll()
        if notify {
            sendNotification()
        }
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
            sendNotification()
        }
    }

    private func sendNotification() {
        NotificationCenter.default.post(name: .WeblDataCacheDidChange, object: self)
    }
}
