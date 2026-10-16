import XCTest

/// Run against the real Go example on the simulator host, never a bundled web view.
final class ElyProfileUITests: XCTestCase {
    @MainActor
    func testNativeLoadValidationSaveAndExplicitReload() throws {
        let app = XCUIApplication()
        app.launchEnvironment["ELY_API_BASE_URL"] = "http://127.0.0.1:8094"
        app.launch()
        let input = app.textFields["input-displayName"]
        XCTAssertTrue(input.waitForExistence(timeout: 15))
        XCTAssertEqual(app.webViews.count, 0)
        assertVisible(input, in: app)
        assertVisible(app.buttons["save-profile"], in: app)
        XCTAssertEqual(app.staticTexts.matching(NSPredicate(format: "label == %@", "Edit your display name.")).count, 1)
        attachScreenshot(app, named: "01 - Native profile loaded")
        replaceText(input, with: "", app: app)
        app.buttons["save-profile"].tap()
        assertVisible(app.staticTexts["error-displayName"], in: app)
        assertVisible(app.buttons["save-profile"], in: app)
        attachScreenshot(app, named: "02 - Native validation error")
        replaceText(input, with: "Native Ada", app: app)
        XCTAssertFalse(app.staticTexts["error-displayName"].exists)
        app.buttons["save-profile"].tap()
        let saved = NSPredicate(format: "label CONTAINS[c] %@", "saved")
        expectation(for: saved, evaluatedWith: app.staticTexts["profile-status"])
        waitForExpectations(timeout: 10)
        XCTAssertEqual(input.value as? String, "Native Ada")
        assertVisible(app.staticTexts["profile-status"], in: app)
        assertVisible(app.buttons["save-profile"], in: app)
        attachScreenshot(app, named: "03 - Native profile saved")
        replaceText(input, with: "Unsaved draft", app: app)
        app.buttons["reload-profile"].tap()
        let keepEditing = app.alerts.buttons["Keep editing"]
        XCTAssertTrue(keepEditing.waitForExistence(timeout: 5))
        keepEditing.tap()
        XCTAssertEqual(input.value as? String, "Unsaved draft")
        app.buttons["reload-profile"].tap()
        let discardAndReload = app.alerts.buttons["Discard edits and reload"]
        XCTAssertTrue(discardAndReload.waitForExistence(timeout: 5))
        discardAndReload.tap()
        let restored = NSPredicate(format: "value == %@", "Native Ada")
        expectation(for: restored, evaluatedWith: input)
        waitForExpectations(timeout: 10)
    }

    @MainActor
    private func assertVisible(_ element: XCUIElement, in app: XCUIApplication,
                               file: StaticString = #filePath, line: UInt = #line) {
        XCTAssertTrue(element.waitForExistence(timeout: 10), file: file, line: line)
        let hittable = XCTNSPredicateExpectation(predicate: NSPredicate(format: "hittable == true"), object: element)
        XCTAssertEqual(XCTWaiter.wait(for: [hittable], timeout: 5), .completed, file: file, line: line)
        XCTAssertTrue(app.frame.contains(element.frame), "Expected the complete element inside the visible screen.", file: file, line: line)
    }

    @MainActor
    private func attachScreenshot(_ app: XCUIApplication, named name: String) {
        let attachment = XCTAttachment(screenshot: app.screenshot())
        attachment.name = name
        attachment.lifetime = .keepAlways
        add(attachment)
    }

    @MainActor
    private func replaceText(_ field: XCUIElement, with text: String, app: XCUIApplication) {
        field.tap()
        let old = field.value as? String ?? ""
        field.typeText(String(repeating: XCUIKeyboardKey.delete.rawValue, count: old.count) + text)
        // The return key may be a checkmark rather than a localized "Done" button.
        // Send Return to the actual control so UITextFieldDelegate handles dismissal.
        field.typeText("\n")
        let dismissed = XCTNSPredicateExpectation(
            predicate: NSPredicate(format: "exists == false"), object: app.keyboards.firstMatch
        )
        XCTAssertEqual(XCTWaiter.wait(for: [dismissed], timeout: 5), .completed,
                       "The native Return key must dismiss the keyboard before saving or taking screenshots.")
    }
}
