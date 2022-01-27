describe("Change ECM", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1/change-ecm");
    });

    it("has headers for different sections", () => {
        cy.get("h1").should("contain", "Change Executive Case Manager");
    });

    it("leaves current ecm blank if none is set", () => {
        cy.get(".govuk-body").should("contain", "Current ECM:");
        cy.get(".govuk-label").should(
            "contain",
            "Enter an Executive Case Manager name"
        );
    });

    it("shows ecm if is set", () => {
        cy.visit("/supervision/deputies/professional/deputy/2/change-ecm");
        cy.get(".govuk-body").should("contain", "Current ECM:");
        cy.get(".govuk-body").should("contain", "ProTeam1 User1");
    });

    it("has a drop down populated with members of the Pro Deputy Teams", () => {
        cy.get("#select-ecm").type("S");
        cy.get("#select-ecm__listbox").find("li").should("have.length", 3);
        cy.get("#select-ecm").type("now");
        cy.get("#select-ecm__listbox").find("li").should("have.length", 1);
    });

    it("directs me back to deputy details page if I press cancel", () => {
        cy.get(".data-emc-cancel").should("contain", "Cancel").click();
        cy.url().should("not.include", "/change-ecm");
        cy.get("h1").should("contain", "Deputy details");
    });

    it("allows me to fill in and submit the ecm form", () => {
        cy.visit("/supervision/deputies/professional/deputy/1/change-ecm");
        cy.setCookie("success-route", "ecm");
        cy.get("#select-ecm").type("S");
        cy.contains("#select-ecm__listbox", "Jon Snow").click();
        cy.get("form").submit();
        cy.get("h1").should("contain", "Deputy details");
        cy.get(".moj-banner--success").should("contain", "Ecm changed to");
    });

    it("displays warning when no ecm chosen and form submitted", () => {
        cy.setCookie("fail-route", "ecm");
        cy.get("#select-ecm").type("S");
        cy.get("form").submit();
        cy.get(".govuk-error-summary").should("contain", "There is a problem");
        cy.get(".govuk-list > li > a").should(
            "contain",
            "Select an executive case manager"
        );
    });
});
