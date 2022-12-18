//
//  PostListViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import Foundation
import OSLog

fileprivate let log = Logger("PostListViewState")

@MainActor
final class PostListViewState: NSObject, ObservableObject {

    @Published var name = ""
    @Published var email = ""
    @Published var posts = [WebClient.Post]()
    @Published var site = WebClient.Site(baseUrl: "…", title: "…", description: "…")

    @Published var showAlert = false
    @Published var error: Error?

    override init() {
        super.init()
        self.reload() // If we've just opened a new window, populate it before going back to the server.
        self.refresh() // Go back to the server and re-populate the data cache.

        NotificationCenter.default.addObserver(forName: .WeblDataCacheDidChange, object: DataCache.shared, queue: .main) { _ in
            Task {
                await self.reload()
            }
        }

        NotificationCenter.default.addObserver(forName: .WeblAccountPreferenceDidChange, object: nil, queue: .main) { _ in
            Task {
                await DataCache.shared.clear() // clear the cache
                await self.reload() // clear the view
                await self.refresh() // retrieve new data from server based on changed account
            }
        }
    }
}

// MARK: - Public API

extension PostListViewState {

    func deletePost(withId id: String) {
        Task {
            do {
                let client = WebClient()
                try await client.deletePost(postId: id)
                DataCache.shared[id] = nil
            } catch {
                showAlert(error: error)
            }
        }
    }

    func createNewPost() {
        Task {
            do {
                let slugline = DataCache.shared.getNewPostName()
                let body = "\n# \(slugline)\n\nThis is where you type something. I mean, compose.\n\n"
                let post = WebClient.Post(
                    id: UUID().uuidString,
                    status: .draft,
                    slugline: slugline,
                    dateCreated: Date(),
                    dateUpdated: Date(),
                    datePublished: Date(),
                    wordCount: body.words,
                    text: body
                )

                let client = WebClient()
                let newPost = try await client.createPost(post: post)
                DataCache.shared[newPost.id] = newPost
            } catch {
                showAlert(error: error)
            }
        }
    }

    func post(id: String?) -> WebClient.Post? {
        return posts.first(where: { $0.id == id })
    }

    func refresh() {
        Task {
            do {
                log.debug("Reloading from server.")
                let client = WebClient()
                let viewerData = try await client.viewerData()
                DataCache.shared.replaceAll(viewer: viewerData)
                reload()
            } catch (let e) {
                showAlert(error: e)
            }
        }
    }

    func toggle(id: String, isPublished: Bool) {
        log.debug("Setting post \(id) to isPublished: \(isPublished).")
        Task {
            do {
                let client = WebClient()
                let post = try await client.togglePost(withId: id, isPublished: isPublished)
                DataCache.shared[post.id] = post                
            } catch (let e) {
                showAlert(error: e)
            }
        }
    }
}

// MARK: - Private Implementation Details

extension PostListViewState {

    private func reload() {
        log.debug("Reloading from cache.")
        self.name = DataCache.shared.name
        self.email = DataCache.shared.email
        self.site = DataCache.shared.site
        self.posts = DataCache.shared.posts
    }

    private func showAlert(error: Error) {
        log.error("\(String(describing: error))")
        self.error = error
        self.showAlert = true
    }
}
