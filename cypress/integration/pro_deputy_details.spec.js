describe("Pro Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1");
    });

    describe("Deputy details", () => {
        it("shows all the deputy details", () => {
            cy.contains(".hook_deputy_name", "firstname surname");
            cy.contains(".hook_deputy_firm_name", "This is the Firm Name")
                .find('a')
                .should("have.attr", "href")
                .and('contain', "/supervision/deputies/firm/0");
            cy.contains(".hook_deputy_phone_number", "1111111");
            cy.contains(".hook_deputy_email", "email@something.com");
        });
    });
});
