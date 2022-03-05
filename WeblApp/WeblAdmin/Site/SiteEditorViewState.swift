//
//  SiteEditorViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/2/22.
//

import Foundation
import OSLog

@MainActor
final class SiteEditorViewState: NSObject, ObservableObject {

    private let log = Logger("SiteEditorViewState")

    @Published var siteTitle = ""
    @Published var siteDescription = ""
    @Published var siteBaseURL = ""

    @Published var accountName = ""
    @Published var accountEmail = ""

    @Published var showAlert = false
    @Published var error: Error?

    @Published var savingSite = false
    @Published var savingAccount = false

    override init() {
        super.init()

        Task {
            await reload()
        }
    }
}

// MARK: - Public API

extension SiteEditorViewState {

    func updateAccount(name: String, email: String) {
        self.savingAccount = true
        Task {
            do {
                let client = WebClient()
                let update = try await client.updateViewer(name: name, email: email)
                DataCache.shared.email = update.email
                DataCache.shared.name = update.name
            } catch (let e) {
                showAlert(error: e)
            }
            self.savingAccount = false
        }
    }

    func updateSite(title: String, description: String, baseURL: String) {
        self.savingSite = true
        Task {
            do {
                let client = WebClient()
                let site = try await client.updateSite(title: title, description: description, baseURL: baseURL)
                DataCache.shared.replace(site: site)
            } catch (let error) {
                showAlert(error: error)
            }
            self.savingSite = false
        }
    }

    var siteDirty: Bool {
        let current = DataCache.shared.site
        if siteTitle.isEmpty {
            return false
        }
        return current.title != siteTitle || current.description != siteDescription || current.baseUrl != siteBaseURL
    }

    var accountDirty: Bool {
        if accountName.isEmpty {
            return false
        }
        return accountName != DataCache.shared.name || accountEmail != DataCache.shared.email
    }
}

// MARK: - Implementation

extension SiteEditorViewState {

    private func reload() {
        Task {
            siteTitle = DataCache.shared.site.title
            siteDescription = DataCache.shared.site.description
            siteBaseURL = DataCache.shared.site.baseUrl

            accountName = DataCache.shared.name
            accountEmail = DataCache.shared.email
        }
    }

    private func showAlert(error: Error) {
        log.error("\(error.localizedDescription)")
        self.showAlert = true
        self.error = error
    }
}