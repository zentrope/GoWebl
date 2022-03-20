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
        view.navigationDelegate = context.coordinator
        return view
    }

    func updateNSView(_ view: WKWebView, context: Context) {

        if let resourceDir = Bundle.main.resourcePath {
            let dir = URL(fileURLWithPath: resourceDir + "/", isDirectory: true)
            view.loadHTMLString(document, baseURL: dir)
        } else {
            view.loadHTMLString("<p>Unable to load document.</p>", baseURL: nil)
        }
    }

    func makeCoordinator() -> Coordinator {
        Coordinator(self)
    }

    class Coordinator: NSObject, WKNavigationDelegate {

        private var parent: WebPostPreview

        init(_ parent: WebPostPreview) {
            self.parent = parent
        }

        func webView(_ webView: WKWebView, decidePolicyFor navigationAction: WKNavigationAction,
                     decisionHandler: @escaping (WKNavigationActionPolicy) -> Void) {
            // Only allow the user to open links in their default browser, rather than
            // within the webview itself.
            if let url = navigationAction.request.url {
                if url.absoluteString.hasPrefix("http") {
                    NSWorkspace.shared.open(url)
                    decisionHandler(WKNavigationActionPolicy.cancel)
                    return
                }
            }
            decisionHandler(WKNavigationActionPolicy.allow)
        }
    }
}
