//
//  AccountPreferencesViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/2/22.
//

import CoreData
import Foundation
import OSLog

@MainActor
class AccountPreferencesViewState: NSObject, ObservableObject {

    private let log = Logger("AccountPreferencesViewState")

    @Published var accounts = [Account]()
    @Published var account = Account()
    @Published var result: TestResult = .untested
    @Published var showAlert = false
    @Published var error: Error?

    enum TestResult {
        case untested
        case succeeded
        case failed(String)
    }

    private lazy var cursor: NSFetchedResultsController<AccountMO> = {
        let fetcher = AccountMO.fetchRequest()
        fetcher.sortDescriptors = [ NSSortDescriptor(key: "name", ascending: true)]
        let cursor = NSFetchedResultsController(fetchRequest: fetcher, managedObjectContext: CoreData.shared.container.viewContext, sectionNameKeyPath: nil, cacheName: nil)
        cursor.delegate = self
        return cursor
    }()

    override init() {
        super.init()
        Task {
            await reload()
        }
    }

    func createNewAccount() {
        Task {
            do {
                try await CoreData.shared.createNewAccount()
            } catch (let e) {
                alert(e)
            }
        }
    }

    func deleteAccount(id: UUID) {
        Task {
            do {
                try await CoreData.shared.deleteAccount(id: id)
            } catch (let e) {
                alert(e)
            }
        }
    }

    func save() {
        Task {
            do {
                try await CoreData.shared.updateAccount(
                    id: account.id,
                    name: account.name,
                    user: account.user,
                    pass: account.password,
                    host: account.host
                )
            } catch {
                alert(error)
            }
        }
    }

    func test() {
        Task {
            do {
                let client = WebClient()
                let rc = try await client.test(user: account.user, pass: account.password, host: account.host)
                self.result = rc ? .succeeded : .failed("Unable to acquire an auth token.")
            } catch {
                self.result = .failed(error.localizedDescription)
            }
        }
    }

    func selectAccount(id: UUID?) {
        self.result = .untested
        guard let accountId = id else {
            self.account = Account()
            return
        }

        do {
            let mo = try CoreData.shared.findAccount(id: accountId)
            self.account = Account(mo)
        } catch {
            alert(error)
        }
    }
}

// MARK: - Implementation

extension AccountPreferencesViewState {

    private func reload() {
        Task {
            do {
                try cursor.performFetch()
                self.accounts = (cursor.fetchedObjects ?? []).map { .init($0) }
            } catch (let error) {
                alert(error)
            }
        }
    }

    private func alert(_ error: Error) {
        log.error("\(error.localizedDescription)")
        self.error = error
        self.showAlert = true
    }
}

// MARK: - Presentation Objects

extension AccountPreferencesViewState {

    struct Account: Identifiable {
        var id = UUID()
        var name = ""
        var user = ""
        var password = ""
        var host = ""

        init() {

        }

        init(_ mo: AccountMO) {
            self.id = mo.id
            self.name = mo.name
            self.user = mo.user
            self.password = mo.password
            self.host = mo.host
        }
    }
}

// MARK: - Cursor Delegate

extension AccountPreferencesViewState: NSFetchedResultsControllerDelegate {
    func controllerDidChangeContent(_ controller: NSFetchedResultsController<NSFetchRequestResult>) {
        reload()
    }
}
