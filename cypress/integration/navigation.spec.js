describe("Navigation bar", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1");
    });

    const expected = [
        ["Dashboard", "/supervision/deputies/professional/deputy/1"],
        ["Timeline", "/supervision/deputies/professional/deputy/1/timeline"],
    ];

    it("has titles and working nav links for all tabs in the correct order", () => {
        cy.get(".moj-sub-navigation__list")
            .children()
            .each(($el, index) => {
                cy.wrap($el).should("contain", expected[index][0]);
                cy.wrap($el).find('a').should("have.attr", "href").and("contain", expected[index][1]);
            });
    });
});