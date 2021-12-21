describe("Manage Pii Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/1/manage-pii-details");
        });

        it("should navigate to the 'Manage pii contact details' page", () => {
            cy.get("[data-cy=manage-pii-details-btn]").click();
            cy.contains(".govuk-heading-l", "Manage professional indemnity insurance");
        });
    });
})