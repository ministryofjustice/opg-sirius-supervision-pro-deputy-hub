describe("Error handling", () => {
    it("renders the error page when the URL does not match a valid route", () => {
        cy.visit("/supervision/deputies/professional/client/1");
        cy.contains(".govuk-heading-l", "Page not found");
    });

    it("handles errors when an Internal Server Error is returned from Sirius", () => {
        cy.visit(
            "/supervision/deputies/professional/deputy/1/manage-deputy-contact-details"
        );

        cy.setCookie("fail-route", "internal-server-error");
        cy.setCookie("fail-code", "500");

        cy.get("form").submit();

        cy.contains(
            ".govuk-heading-l",
            "Sorry, there is a problem with the service"
        );
    });
});
