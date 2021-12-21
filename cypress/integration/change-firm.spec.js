describe("Change Firm", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Changing a firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/1/change-firm");
        });

        it("shows title for page", () => {
            cy.get(".govuk-grid-column-full > header").should(
                "contain",
                "Change firm"
            );
        });

        it("shows current firm name", () => {
            cy.get(".govuk-body").should("contain", "Current firm");
        });

        it("has a save button that can redirect to add-note page", () => {
            cy.get("#new-firm").click();
            cy.get(".govuk-button")
                .should("contain", "Save and continue")
                .click();
            cy.url().should(
                "contain",
                "/supervision/deputies/professional/deputy/1/add-firm"
            );
        });

        it("has a button that will show existing firm details", () => {
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown").should("be.visible");
        });

        it("has a cancel button that can redirect to deputy details page", () => {
            cy.get(".govuk-link").should("contain", "Cancel").click();
            cy.url().should(
                "contain",
                "/supervision/deputies/professional/deputy/1"
            );
        });
    });
    describe("Changing a firm to an existing firm", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/1/change-firm");
        });

        it("has a dropdown with the existing firm options", () => {
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("Firm");
            cy.get("#select-existing-firm__listbox")
                .find("li")
                .should("have.length", 2);
        });

        it("will redirect and show success banner when deputy allocated to firm", () => {
            cy.setCookie("success-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("Great");
            cy.contains(
                "#select-existing-firm__option--0",
                "Great Firm Corp - 1000002"
            ).click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".moj-banner").should("contain", "Firm changed to");
            cy.get("h1").should("contain", "Deputy details");
        });

        it("will allow searching based on firm id", () => {
            cy.setCookie("success-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm-dropdown > .govuk-label").should(
                "contain",
                "Enter a firm name or number"
            );
            cy.get("#select-existing-firm").click().type("1000002");
            cy.contains(
                "#select-existing-firm__option--0",
                "Great Firm Corp - 1000002"
            ).click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".moj-banner").should("contain", "Firm changed to");
            cy.get("h1").should("contain", "Deputy details");
        });

        it("will show a validation error if no options available", () => {
            cy.setCookie("fail-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#select-existing-firm")
                .click()
                .type("Unknown option for firm name");
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "Enter a firm name or number"
                );
            });
        });

        it("will show a validation error if form submitted when autocomplete empty", () => {
            cy.setCookie("fail-route", "allocateToFirm");
            cy.get("#f-existing-firm").click();
            cy.get("#existing-firm-or-new-firm-form").submit();
            cy.get(".govuk-error-summary__title").should(
                "contain",
                "There is a problem"
            );
            cy.get(".govuk-error-summary__list").within(() => {
                cy.get("li:first").should(
                    "contain",
                    "Enter a firm name or number"
                );
            });
        });
    });
});
