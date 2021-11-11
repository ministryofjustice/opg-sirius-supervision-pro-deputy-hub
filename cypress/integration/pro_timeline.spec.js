describe("Timeline", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1/timeline");
    });

    it("has a header called timeline", () => {
        cy.get(".main > header").should("contain", "Timeline");
    })

    it("contains appropriate test data for a timeline event", () => {
        cy.get(".moj-timeline__title").should("contain", "New client added to deputyship");
        cy.get(".moj-timeline__byline").should("contain", "by system admin (12345678)");
        cy.get("time").should("contain", "2021-09-09 14:01:59");
        cy.get(".govuk-list > :nth-child(1)").should("contain", "Order number: 03305972");
        cy.get(".govuk-list > :nth-child(2)").should("contain", "Sirius ID: 7000-0000-1995");
        cy.get(".govuk-list > :nth-child(3)").should("contain", "Order type: pfa");
        cy.get(".govuk-list > :nth-child(4)").should("contain", "Client: Duke John Fearless");
    })
});