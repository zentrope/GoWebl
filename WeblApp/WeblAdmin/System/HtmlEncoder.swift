//
//  HtmlEncoder.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 2/23/22.
//

import Foundation
import Markdown

final class HTMLEncoder {

    var markdown: String

    init(_ markdown: String) {
        self.markdown = markdown
    }

    func html() -> String {
        let doc = Document(parsing: markdown)
        var decoder = MarkdownDecoder()
        decoder.visit(doc)
        return decoder.elements.joined(separator: "")
    }
}

struct MarkdownDecoder: MarkupWalker {

    var elements = [String]()

    mutating func visitHeading(_ heading: Heading) -> () {
        let tag = "h\(heading.level)"
        elements.append("<\(tag)>")
        descendInto(heading)
        elements.append("</\(tag)>")
    }

    // BLOCKS

    mutating func visitDocument(_ document: Document) -> () {
        elements.append(contentsOf: [
            "<!doctype html>",
            #"<html lang="en">"#,
            "<head>",
            #"<meta charset="UTF-8">"#,
            #"<link href="style.css" type="text/css" rel="stylesheet"/>"#,
            "</head>",
            "<body>"
        ])
        descendInto(document)
        elements.append("</body></html>")
    }

    mutating func visitParagraph(_ paragraph: Paragraph) -> () {
        elements.append("<p>")
        descendInto(paragraph)
        elements.append("</p>")
    }

    mutating func visitThematicBreak(_ thematicBreak: ThematicBreak) -> () {
        elements.append("<hr/>")
    }

    mutating func visitBlockQuote(_ blockQuote: BlockQuote) -> () {
        elements.append("<blockquote>")
        descendInto(blockQuote)
        elements.append("</blockquote>")
    }

    // TEXT

    mutating func visitText(_ text: Text) -> () {
        elements.append(text.string)
        descendInto(text)
    }

    // LINK

    mutating func visitLink(_ link: Link) -> () {
        // NOTE: This library does not support titles (the text in quotes after the link).
        let text = link.plainText
        let href = link.destination ?? "#"
        let link = #"<a href="\#(href)">\#(text)</a>"#
        elements.append(link)
    }

    // STYLE

    mutating func visitInlineCode(_ inlineCode: InlineCode) -> () {
        elements.append(contentsOf: [
            "<code>",
            inlineCode.code.escapeHtml,
            "</code>"
        ])
    }

    mutating func visitEmphasis(_ emphasis: Emphasis) -> () {
        elements.append("<i>")
        descendInto(emphasis)
        elements.append("</i>")
    }

    mutating func visitStrong(_ strong: Strong) -> () {
        elements.append("<b>")
        descendInto(strong)
        elements.append("</b>")
    }

    // CODE

    mutating func visitCodeBlock(_ codeBlock: CodeBlock) -> () {
        elements.append("<pre>")
        elements.append(codeBlock.code.escapeHtml)
        elements.append("</pre>")
    }

    mutating func visitInlineHTML(_ inlineHTML: InlineHTML) -> () {
        elements.append(inlineHTML.rawHTML)
    }

    mutating func visitHTMLBlock(_ html: HTMLBlock) -> () {        
        elements.append(html.rawHTML)
    }

//    mutating func defaultVisit(_ markup: Markup) -> () {
//        print("DEFAULT VISIT: \(markup)")
//        descendInto(markup)
//    }

    // LISTS

    mutating func visitUnorderedList(_ unorderedList: UnorderedList) -> () {
        elements.append("<ul>")
        descendInto(unorderedList)
        elements.append("</ul>")
    }

    mutating func visitListItem(_ listItem: ListItem) -> () {
        elements.append("<li>")
        descendInto(listItem)
        elements.append("</li>")
    }
}
