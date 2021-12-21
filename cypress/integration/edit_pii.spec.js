describe("Manage Pii Details", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    describe("Navigation", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/1/manage-pii-details");
        });

        it("should navigate to the 'Manage pii contact details' page with relevant labels", () => {
            cy.get('.govuk-heading-l').should('contain', 'Manage professional indemnity insurance');
            cy.get(':nth-child(2) > .govuk-fieldset__legend > .govuk-fieldset__heading').should('contain', 'Current PII certificate')
            cy.get(':nth-child(2) > .govuk-label').should('contain', 'Date received')
            cy.get(':nth-child(3) > .govuk-label').should('contain', 'Expiry date')
            cy.get(':nth-child(4) > .govuk-label').should('contain', 'Amount')
            cy.get(':nth-child(3) > .govuk-fieldset__legend > .govuk-fieldset__heading').should('contain', 'New PII certificate')
            cy.get(':nth-child(3) > .govuk-form-group > .govuk-label').should('contain', 'Date requested')
        });

        it("should allow a form to be submitted and get a success message", () => {
            cy.setCookie("success-route", "editPii");
            cy.get('#pii-received').type('2020-01-01')
            cy.get('#pii-expiry').type('2020-01-02')
            cy.get('#pii-amount').type('1234')
            cy.get('#pii-requested').type('2020-01-03')
            cy.get('.govuk-button').click()
            cy.get('.moj-banner').should('contain', 'Pii details updated')
            cy.get('.govuk-heading-l').should('contain', 'Deputy details')
        });
    });
})