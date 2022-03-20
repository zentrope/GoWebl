//
//  AccountPreferencesViewState.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 3/2/22.
//

import CoreData
import OSLog
import SwiftUI

@MainActor
class AccountPreferencesViewState: NSObject, ObservableObject {

    private let log = Logger("AccountPreferencesViewState")

    @Published var accounts = [Account]()
    @Published var account = Account()
    @Published var selectedAccount: UUID?
    @Published var result: TestResult = .untested
    @Published var showAlert = false
    @Published var error: Error?

    @AppStorage("WMAccountID") private var savedAccountId = ""
    @AppStorage("WAAccountEmail") private var savedEmail = ""
    @AppStorage("WAAccountPassword") private var savedPassword = ""
    @AppStorage("WAAccountEndpoint") private var savedEndpoint = ""

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
            await setInitialSelection()
        }
    }

    var isDefault: Bool {
        savedAccountId == account.id.uuidString
    }

    func isDefault(id: UUID) -> Bool {
        id.uuidString == savedAccountId
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

    func setAsDefault() {
        savedAccountId = account.id.uuidString
        savedEmail = account.user
        savedPassword = account.password
        savedEndpoint = account.host
    }

    func test() {
        self.result = .untested
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

    private func setInitialSelection() {
        Task {
            // HACK: Delaying assignment seems to be necessary. My guess is that List and Table don't finish initializing until some time after they're actually rendered thus overwriting whatever selection we set. Or it's a bug.
            try await Task.sleep(nanoseconds: 100_000)
            let accountId = UUID(uuidString: savedAccountId) ?? self.accounts.first?.id
            selectedAccount = accountId
        }
    }

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
