import Foundation
import SQLite

extension MessageStore {
  static func detectAttributedBody(connection: Connection) -> Bool {
    do {
      let rows = try connection.prepare("PRAGMA table_info(message)")
      for row in rows {
        if let name = row[1] as? String,
          name.caseInsensitiveCompare("attributedBody") == .orderedSame
        {
          return true
        }
      }
    } catch {
      return false
    }
    return false
  }

  static func enhance(error: Error, path: String) -> Error {
    let message = String(describing: error).lowercased()
    if message.contains("out of memory (14)") || message.contains("authorization denied")
      || message.contains("unable to open database") || message.contains("cannot open")
    {
      return IMsgError.permissionDenied(path: path, underlying: error)
    }
    return error
  }

  func appleDate(from value: Int64?) -> Date {
    guard let value else { return Date(timeIntervalSince1970: MessageStore.appleEpochOffset) }
    return Date(
      timeIntervalSince1970: (Double(value) / 1_000_000_000) + MessageStore.appleEpochOffset)
  }

  func stringValue(_ binding: Binding?) -> String {
    return binding as? String ?? ""
  }

  func int64Value(_ binding: Binding?) -> Int64? {
    if let value = binding as? Int64 { return value }
    if let value = binding as? Int { return Int64(value) }
    if let value = binding as? Double { return Int64(value) }
    return nil
  }

  func intValue(_ binding: Binding?) -> Int? {
    if let value = binding as? Int { return value }
    if let value = binding as? Int64 { return Int(value) }
    if let value = binding as? Double { return Int(value) }
    return nil
  }

  func boolValue(_ binding: Binding?) -> Bool {
    if let value = binding as? Bool { return value }
    if let value = intValue(binding) { return value != 0 }
    return false
  }

  func dataValue(_ binding: Binding?) -> Data {
    if let blob = binding as? Blob {
      return Data(blob.bytes)
    }
    return Data()
  }
}
