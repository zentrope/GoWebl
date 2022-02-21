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

    @Published var posts = [WebClient.Post]()
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
                let posts = try await client.viewerData()
                self.posts = posts.sorted(by: { $0.dateCreated > $1.dateCreated })
                log.debug("Downloaded \(self.posts.count) posts.")
                if let fp = self.posts.first {
                    log.debug("First \(String(describing: fp))")
                }
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
