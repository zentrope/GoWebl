//
//  PostListViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import Foundation
import OSLog

fileprivate let log = Logger(subsystem: "com.zentrope.WeblAdmin", category: "PostListViewState")

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

        //Task { await self.reload() }
    }

}

// MARK: - Public API

extension PostListViewState {

    func post(id: String?) -> WebClient.Post? {
        return posts.first(where: { $0.id == id })
    }

    func toggle(id: String, isPublished: Bool) {
        Task {
            do {
                let client = WebClient()
                try await client.togglePost(withId: id, isPublished: isPublished)
                reload(client)
            } catch (let e) {
                showAlert(error: e)
            }
        }
    }

    func refresh() {
        reload()
    }
}

// MARK: - Private Implementation Details

extension PostListViewState {

    private func reload(_ wc: WebClient? = nil) {
        Task {
            do {
                let client = wc ?? WebClient()
                let viewerData = try await client.viewerData()
                self.name = viewerData.name
                self.email = viewerData.email
                self.site = viewerData.site
                self.posts = viewerData.posts.sorted(by: { $0.dateCreated > $1.dateCreated })
            } catch (let e) {
                showAlert(error: e)
            }
        }
    }

    private func showAlert(error: Error) {
        log.error("\(String(describing: error))")
        self.error = error
        self.showAlert = true
    }
}
