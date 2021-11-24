describe("Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Adding a firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/1/add-firm");
        });

        describe("Successfully adding a firm", () => {
            // it("has a add a firm page with expected fields", () => {
            //     cy.get(":nth-child(2) > .govuk-label").should("contain", "Title (required)")
            //     cy.get(".govuk-character-count > .govuk-form-group > .govuk-label").should("contain", "Note (required)")
            //     cy.get("#note-info").should("contain", "You have 1000 characters remaining")
            //     cy.get(".govuk-button").should("contain", "Save note")
            //     cy.get(".govuk-link").should("contain", "Cancel")
            // })
        })

        it("shows error message when submitting invalid data", () => {
            cy.setCookie("fail-route", "firm")
            cy.get("#add-firm-form").submit();
            cy.get(".govuk-error-summary__title").should("contain", "There is a problem");
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should("contain", "The title must be 255 characters or fewer");
                cy.get("li:last").should("contain", "Enter a note");
            })
        });
    })
});