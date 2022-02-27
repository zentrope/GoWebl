//
//  WebClient.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import AnyCodable
import SwiftUI
import OSLog

fileprivate let log = Logger("WebClient")

// MARK: - WebClient

final class WebClient: NSObject {

    @AppStorage("WAAccountEmail") private var email = ""
    @AppStorage("WAAccountPassword") private var password = ""

    private var token = ""

    override init() {
        super.init()
    }

    private func doQuery(_ query: Query, token: String? = nil) async throws -> GQLResponse {
        let encoder = JSONEncoder()
        let payload = try encoder.encode(query)

        let path = "https://devtrope.com/query"
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

//        guard let doc = String(data: data, encoding: .utf8) else {
//            log.error("cannot decode raw data into UTF8 string")
//            throw URLError(.cannotDecodeRawData)
//        }
//
//        log.debug("decoding: \(doc)")

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

    private func login() async throws -> String {
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
        let query = Query(query: gql, operationName: "Authenticate", variables: ["email" : AnyCodable(self.email), "pass" : AnyCodable(self.password)])

        let result = try await doQuery(query)

        let token = result.data.authenticate?.token ?? ""
        self.token = token
        return token
    }
}

// MARK: - Public API

extension WebClient {

    func test() async throws -> Bool {
        let result = try await login()
        return !result.isEmpty
    }

    func togglePost(withId uuid: String, isPublished: Bool) async throws {
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
            variables: ["uuid" : AnyCodable(uuid), "isPublished" : AnyCodable(isPublished)]
        )

        let result = try await doQuery(mutation, token: token)
        log.debug("\(String(describing: result))")

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

// MARK: - Custom Errors

extension WebClient {

    enum GraphQlError: Error, LocalizedError {
        case Error(String)
        case NoViewerData

        var errorDescription: String? {
            switch self {
                case let .Error(msg): return "gql: \(msg)"
                case .NoViewerData: return "unable to retrieve data, check account preferences"
            }
        }
    }
}

// MARK: - Post Objects

extension WebClient {

    struct Query: Encodable {
        var query: String
        var operationName: String
        var variables: [String:AnyCodable]
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
        var id: String
        var status: Status
        var slugline: String
        var dateCreated: Date
        var dateUpdated: Date
        var datePublished: Date
        var wordCount: Int
        var text: String

        private enum CodingKeys: String, CodingKey {
            case id = "uuid", status, slugline, dateCreated, dateUpdated, datePublished, wordCount, text
        }

        enum Status: String, Codable {
            case published
            case draft
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

