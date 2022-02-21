//
//  WebClient.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/20/22.
//

import SwiftUI
import OSLog

fileprivate let log = Logger(subsystem: "com.zentrop.WeblAdmin", category: "WebClient")

class WebClient: NSObject {

    @AppStorage("WAAccountEmail") private var email = ""
    @AppStorage("WAAccountPassword") private var password = ""

    private var token = ""

    override init() {
        super.init()
    }

    func test() async throws -> Bool {
        let result = try await login()
        return !result.isEmpty
    }

    func viewerData() async throws -> [Post] {
        let token = try await login()

        let ql = """
        query {
          viewer { id name type email
            site { baseUrl title description }
            posts { uuid status slugline dateCreated dateUpdated datePublished text wordCount }}}
        """
        let query = Query(query: ql, operationName: "", variables: [:])
        let result = try await doQuery(query, token: token)
        return result.data.viewer?.posts ?? []
        //print(result)
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

        guard let doc = String(data: data, encoding: .utf8) else {
            log.error("cannot decode raw data into UTF8 string")
            throw URLError(.cannotDecodeRawData)
        }

        log.debug("decoding: \(doc)")

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let gdoc = try decoder.decode(GQLResponse.self, from: data)
        log.debug("decoded: \(String(describing: gdoc))")

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
        let query = Query(query: gql, operationName: "Authenticate", variables: ["email" : self.email, "pass" : self.password])

        let result = try await doQuery(query)

        let token = result.data.authenticate?.token ?? ""
        self.token = token
        return token
    }
}

// MARK: - Custom Errors

extension WebClient {

    enum GraphQlError: Error, LocalizedError {
        case Error(String)

        var errorDescription: String? {
            switch self {
                case let .Error(msg): return "gql: \(msg)"
            }
        }
    }
}

// MARK: - Post Objects

extension WebClient {

    struct Query: Encodable {
        var query: String
        var operationName: String
        var variables: [String:String]
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

    struct Post: Decodable, Identifiable {
        var id: String
        var status: String
        var slugline: String
        var dateCreated: Date
        var dateUpdated: Date
        var datePublished: Date
        var wordCount: Int
        var text: String

        private enum CodingKeys: String, CodingKey {
            case id = "uuid", status, slugline, dateCreated, dateUpdated, datePublished, wordCount, text
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


extension URLRequest {

    init?(method: String = "GET", host: String, path: String, auth: String, params: [String:String] = [:]) {

        let resource = "https://\(host)\(path)"

        guard var url = URLComponents(string: resource) else {
            return nil
        }

        if !params.isEmpty {
            url.queryItems = params.map { URLQueryItem(name: $0, value: $1) }
        }
        self.init(url: url.url!)
        httpMethod = method
        addValue("application/json", forHTTPHeaderField: "Accept")
        addValue("application/json", forHTTPHeaderField: "Content-Type")
        addValue("1", forHTTPHeaderField: "X-PrettyPrint")
        addValue("Bearer \(auth)", forHTTPHeaderField: "Authorization")
    }
}
