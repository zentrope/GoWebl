//
//  PostSourceEditorViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/26/22.
//

import Foundation
import OSLog

@MainActor
class PostSourceEditorViewState: NSObject, ObservableObject {

    @Published var post = WebClient.Post()
    @Published var showAlert = false
    @Published var error: Error?

    let log = Logger("PostSourceEditorViewState")

    func update(post: WebClient.Post, newText: String) {
        Task {
            do {
                let client = WebClient()
                let updatedPost = try await client.updatePost(uuid: post.id, slugline: post.slugline, text: newText, datePublished: post.datePublished)
                DataCache.shared[updatedPost.id] = updatedPost
                self.post = updatedPost
                log.debug("Updated: \(updatedPost.id)")
            } catch (let err) {
                showAlert(error: err)
            }
        }
    }

    func setPost(toPostWithId id: String) {
        self.post = DataCache.shared[id] ?? WebClient.Post()
    }

    func showAlert(error: Error) {
        self.showAlert = true
        self.error = error
    }
}
