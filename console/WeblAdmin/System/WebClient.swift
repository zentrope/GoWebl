//
//  WebClient.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI
import OSLog

fileprivate let log = Logger("WebClient")

// MARK: - WebClient

final class WebClient: NSObject {

    @AppStorage("WAAccountEmail") private var email = ""
    @AppStorage("WAAccountPassword") private var password = ""
    @AppStorage("WAAccountEndpoint") private var endpoint = ""

    private var token = ""

    override init() {
        super.init()
    }
}

// MARK: - Public API

extension WebClient {

    func createPost(post: Post) async throws -> Post {
        let token = try await login()
        let ql = """
            mutation
            CreatePost($slugline: String! $status: String! $text: String! $d: String! $token: String) {
              createPost(slugline: $slugline, status: $status, text: $text, datePublished: $d, token: $token) {
                uuid slugline status dateCreated dateUpdated datePublished text wordCount }}
        """
        let mutation = Query(query: ql, operationName: "CreatePost", variables: [
            "slugline": Param(post.slugline),
            "status": Param(post.status.rawValue),
            "text" : Param(post.text),
            "d": Param(post.datePublished)
        ])
        let result = try await doQuery(mutation, token: token)
        if let post = result.data.createPost {
            return post
        }
        throw GraphQlError.NoViewerData
    }

    func deletePost(postId id: String) async throws {
        let token = try await login()
        let ql = "mutation DeletePost($uuid: String!) { deletePost(uuid: $uuid) } "
        let mutation = Query(query: ql, operationName: "DeletePost", variables: ["uuid": Param(id)])
        let result = try await doQuery(mutation, token: token)
        log.debug("\(String(describing: result))")
    }

    func test() async throws -> Bool {
        let result = try await login()
        return !result.isEmpty
    }

    func test(user: String, pass: String, host: String) async throws -> Bool {
        let result = try await login(user: user, pass: pass, host: host)
        return !result.isEmpty
    }

    func setPassword(toNewPassword password: String) async throws {
        let token = try await login()
        let ql = "mutation UpdateViewerPassword($p: String!) { updateViewerPassword(password: $p) { id }}"
        let mutation = Query(query: ql, operationName: "UpdateViewerPassword", variables: ["p": Param(password)])
        let result = try await doQuery(mutation, token: token)
        if let viewerId = result.data.updateViewerPassword {
            log.debug("\(String(describing: viewerId))")
            return
        }
        throw GraphQlError.UnableToResetPassword
    }

    func togglePost(withId uuid: String, isPublished: Bool) async throws -> Post {
        let token = try await login()
        let ql = """
        mutation
            SetPostStatus($uuid: String!, $isPublished: Boolean!) {
              setPostStatus(uuid: $uuid, isPublished: $isPublished) {
                uuid slugline status dateCreated dateUpdated datePublished wordCount text }}
        """
        let mutation = Query(
            query: ql,
            operationName: "SetPostStatus",
            variables: ["uuid" : Param(uuid), "isPublished" : Param(isPublished)]
        )

        let result = try await doQuery(mutation, token: token)
        if let post = result.data.setPostStatus {
            return post
        }
        log.debug("\(String(describing: result))")
        throw GraphQlError.NoViewerData
    }

    func updatePost(uuid: String, slugline: String, text: String, datePublished: Date) async throws -> Post {
        let token = try await login()
        let ql = """
            mutation
            UpdatePost($u: String! $s: String! $t: String! $d: String!) {
              updatePost(uuid: $u slugline: $s text: $t datePublished: $d) {
                uuid slugline status dateCreated dateUpdated datePublished text wordCount}}
        """
        let mutation = Query(query: ql, operationName: "UpdatePost", variables: [
            "u" : Param(uuid), "s": Param(slugline), "t" : Param(text), "d": Param(datePublished)
        ])
        let result = try await doQuery(mutation, token: token)

        if let post = result.data.updatePost {
            return post
        }
        log.debug("\(String(describing: result))")
        throw GraphQlError.NoViewerData
    }

    func updateSite(title: String, description: String, baseURL: String) async throws -> Site {
        let token = try await login()
        let ql = """
        mutation
            UpdateSite($t: String! $d: String! $b: String!) {
              updateSite(title: $t description: $d baseUrl: $b) {
                title description baseUrl }}
        """
        let mutation = Query(query: ql, operationName: "UpdateSite", variables: [
            "t" : Param(title),
            "d" : Param(description),
            "b" : Param(baseURL)
        ])
        let result = try await doQuery(mutation, token: token)
        if let site = result.data.updateSite {
            return site
        }
        throw GraphQlError.NoViewerData
    }

    func updateViewer(name: String, email: String) async throws -> GData.ViewerUpdate {
        let token = try await login()
        let ql = "mutation UpdateViewer($n: String! $e: String!) { updateViewer(name: $n, email: $e) { name email }}"
        let mutation = Query(query: ql, operationName: "UpdateViewer", variables: ["n":Param(name), "e":Param(email)])
        let result = try await doQuery(mutation, token: token)
        if let update = result.data.updateViewer {
            return update
        }
        throw GraphQlError.NoViewerData
    }

    func viewerData() async throws -> Viewer {
        let token = try await login()

        let ql = """
        query {
          viewer { id name type email
            site { baseUrl title description }
            posts { uuid status slugline dateCreated dateUpdated datePublished text wordCount }}}
        """
        let query = Query(query: ql, operationName: "", variables: [:])
        let result = try await doQuery(query, token: token)
        if let viewer = result.data.viewer {
            return viewer
        }
        throw GraphQlError.NoViewerData
    }
}

// MARK: - Private Implementation Details

extension WebClient {

    private func doQuery(_ query: Query, token: String? = nil, host: String? = nil) async throws -> GQLResponse {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        let payload = try encoder.encode(query)

        let path = "\(host ?? self.endpoint)/query"
        guard let url = URLComponents(string: path) else { throw URLError(.badURL) }

        var request = URLRequest(url: url.url!)
        request.httpMethod = "POST"
        request.httpBody = payload
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue("application/json", forHTTPHeaderField: "Accept")
        if let auth = token {
            request.addValue("Bearer \(auth)", forHTTPHeaderField: "Authorization")
        }

        let (data, _) = try await URLSession.shared.data(for: request)

//        if !query.operationName.isEmpty {
//            guard let doc = String(data: data, encoding: .utf8) else {
//                log.error("cannot decode raw data into UTF8 string")
//                throw URLError(.cannotDecodeRawData)
//            }
//
//            log.debug("\n\ndecoding: \(doc)\n\n")
//        }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601

        let gdoc = try decoder.decode(GQLResponse.self, from: data)

        if gdoc.hasError() {
            let msg = gdoc.errors?.first?.message ?? "Unable to decode error, check log"
            log.error("\(msg)")
            throw GraphQlError.Error(msg)
        }
        return gdoc
    }

    private func login(user: String, pass: String, host: String) async throws -> String {
        if !token.isEmpty {
            return token
        }
        let gql = """
        query Authenticate($email: String! $pass: String!) {
          authenticate(email: $email pass: $pass) {
            token
          }
        }
        """
        let query = Query(query: gql, operationName: "Authenticate", variables: ["email" : Param(user), "pass" : Param(pass)])

        let result = try await doQuery(query, host: host)

        let token = result.data.authenticate?.token ?? ""
        self.token = token
        return token
    }

    private func login() async throws -> String {
        return try await login(user: self.email, pass: self.password, host: self.endpoint)
    }
}

// MARK: - Custom Errors

extension WebClient {

    enum GraphQlError: Error, LocalizedError {
        case Error(String)
        case NoViewerData
        case UnableToResetPassword

        var errorDescription: String? {
            switch self {
                case let .Error(msg): return "gql: \(msg)"
                case .NoViewerData: return "unable to retrieve data, check account preferences"
                case .UnableToResetPassword: return "unable to reset password"
            }
        }
    }
}

// MARK: - Post Objects

extension WebClient {

    /// Provides a wrapper around arbitrary values used as parameter lists for GraphQL queries and mutations.
    struct Param<T>: Encodable {

        var value: T

        init(_ value: T) {
            self.value = value
        }

        func encode(to encoder: Encoder) throws {
            var container = encoder.singleValueContainer()
            switch self.value {
                case let value as Bool:
                    try container.encode(value)
                case let value as String:
                    try container.encode(value)
                case let value as Int:
                    try container.encode(value)
                case let value as Date:
                    try container.encode(value)
                default:
                    let context = EncodingError.Context(
                        codingPath: container.codingPath,
                        debugDescription: "Param value cannot be encoded")
                                throw EncodingError.invalidValue(value, context)
            }
        }
    }

    struct Query: Encodable {
        var query: String
        var operationName: String
        var variables: [String:Param<Any>]
    }

    struct GQLResponse: Decodable {
        var data: GData
        var errors: [GQLError]?

        func hasError() -> Bool {
            errors != nil
        }
    }

    struct GData: Decodable {
        var authenticate: Token?
        var viewer: Viewer?
        var setPostStatus: Post?
        var updatePost: Post?
        var createPost: Post?
        var updateSite: Site?
        var updateViewer: ViewerUpdate?
        var updateViewerPassword: ViewerId?

        struct ViewerId: Decodable {
            var id: String
        }

        struct ViewerUpdate: Decodable {
            var name: String
            var email: String
        }
    }

    struct Viewer: Decodable {
        var id: String
        var name: String
        var type: String
        var email: String
        var site: Site
        var posts: [Post]
    }

    struct Site: Decodable {
        var baseUrl: String
        var title: String
        var description: String
    }

    struct Post: Decodable, Identifiable, Equatable {
        var id: String = ""
        var status: Status = Status.draft
        var slugline: String = ""
        var dateCreated: Date = Date.distantPast
        var dateUpdated: Date = Date.distantPast
        var datePublished: Date = Date.distantPast
        var wordCount: Int = 0
        var text: String = ""

        private enum CodingKeys: String, CodingKey {
            case id = "uuid", status, slugline, dateCreated, dateUpdated, datePublished, wordCount, text
        }

        enum Status: String, Codable {
            case published = "published"
            case draft = "draft"
        }
    }

    struct GQLError: Decodable {
        var message: String
        var path: [String]?
        var locations: [Location]?
    }

    struct Location: Decodable {
        var line: Int
        var column: Int
    }

    struct Authenticate: Decodable {
        var authenticate: Token?
    }

    struct Token: Decodable {
        var token: String
    }
}
