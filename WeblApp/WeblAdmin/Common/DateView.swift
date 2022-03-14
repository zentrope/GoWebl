//
//  DateView.swift
//  WeblAdmin
//
//  Created by Keith Irwin on 12/30/21.
//

import SwiftUI

struct DateView: View {

    enum Format: String {
        /// "MM-dd-yyyy hh:mm a"
        case dateTimeMonthFirst = "MM-dd-yyyy hh:mm a"
        /// "yyyy-MM-dd"
        case dateShort = "yyyy-MM-dd"
        /// "yyyy-MM-dd hh:mm a"
        case dateTime = "yyyy-MM-dd hh:mm:ss a"

        /// "yyyy-MM-dd hh:mm"
        case dateTimeLog = "yyyy-MM-dd HH:mm:ss"

        /// "yyyy-MM-dd'T'HH:mm:ssZ"
        case iso8601 = "yyyy-MM-dd'T'HH:mm:ssZ"
        /// "MMMM d, yyyy @ hh:mm a"
        case dateTimeNameLong = "MMMM d, yyyy @ hh:mm a"
        /// "MMMM d, yyyy"
        case dateNameMedium = "MMMM d, yyyy"
        /// "MMM dd, yyyy"
        case dateNameShort = "MMM dd, yyyy"
        /// "EEEE, MMMM d, yyyy' - 'h:mm a"
        case dateDayTimeNameLong = "EEEE, MMMM d, yyyy' - 'h:mm a"
        /// Tue, Dec 21, 2021
        case nameMonthDayYear = "EEE, MMM d, yyyy"
        /// "MMMM yyyy"
        case dateMonthYear = "MMMM yyyy"

        /// 15 Dec 21 - 14:37 PM
        case dateDense = "dd MMM yy hh:mm a"


        case overdueSince = "formatOverdue()"
    }

    var date: Date?
    var format: Format
    var ifNil: String = "…"

    var body: some View {
        Text(fmt(date, using: format))
    }

    private func fmt(_ date: Date?, using format: Format) -> String {
        guard let date = date else { return ifNil }
        if format == .overdueSince {
            return formatOverdue(date)
        }
        let formatter = DateFormatter()
        formatter.dateFormat = format.rawValue
        let result = formatter.string(from: date)
        return result
    }

    private func formatOverdue(_ date: Date) -> String {
        let interval = abs(date.timeIntervalSince(Date()))

        let hours = interval / 60 / 60

        if hours > (3 * 24) {
            return String(format: "Overdue %.0f days", hours / 24)
        }

        let rounding = hours < 10 ? "Overdue %.2f hours" : "Overdue %.0f hours"
        return String(format: rounding, interval / 60 / 60)
    }

}
