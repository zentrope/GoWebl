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

        Task { await self.reload() }
    }

    func post(id: String?) -> WebClient.Post? {
        return posts.first(where: { $0.id == id })
    }

    private func reload() {
        Task {
            do {
                let client = WebClient()
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
