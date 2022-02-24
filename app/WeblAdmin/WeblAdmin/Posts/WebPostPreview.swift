//
//  WebPostPreview.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/23/22.
//

import SwiftUI
import WebKit

struct WebPostPreview: NSViewRepresentable {
    var document: String

    func makeNSView(context: Context) -> WKWebView {
        let view = WKWebView()
        view.setValue(false, forKey: "drawsBackground")
        return view
    }

    func updateNSView(_ view: WKWebView, context: Context) {

        if let resourceDir = Bundle.main.resourcePath {
            let dir = URL(fileURLWithPath: resourceDir + "/", isDirectory: true)
            let doc = wrap(document)
            view.loadHTMLString(doc, baseURL: dir)
        } else {
            view.loadHTMLString("<p>Unable to load document.</p>", baseURL: nil)
        }
    }

    private func wrap(_ doc: String) -> String {
        let head = #"<html><head><link href="style.css" type="text/css" rel="stylesheet"/></head><body>"#
        let foot = #"</body></html>"#
        return "\(head)\(doc)\(foot)"
    }
}
