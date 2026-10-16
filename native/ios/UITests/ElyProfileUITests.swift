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
        attachScreenshot(app, named: "01 - Native profile loaded")
        replaceText(input, with: "", app: app)
        app.buttons["save-profile"].tap()
        XCTAssertTrue(app.staticTexts["error-displayName"].waitForExistence(timeout: 10))
        attachScreenshot(app, named: "02 - Native validation error")
        replaceText(input, with: "Native Ada", app: app)
        XCTAssertFalse(app.staticTexts["error-displayName"].exists)
        app.buttons["save-profile"].tap()
        let saved = NSPredicate(format: "label CONTAINS[c] %@", "saved")
        expectation(for: saved, evaluatedWith: app.staticTexts["profile-status"])
        waitForExpectations(timeout: 10)
        XCTAssertEqual(input.value as? String, "Native Ada")
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
        if app.keyboards.buttons["Done"].exists { app.keyboards.buttons["Done"].tap() }
    }
}
