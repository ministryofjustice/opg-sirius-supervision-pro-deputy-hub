describe("Timeline", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1/timeline");
    });

    it("has a header called timeline", () => {
        cy.get(".main > header").should("contain", "Timeline");
    });

    it("lists timeline events in date ascending order", () => {
        const timelineItems = cy.get(".moj-timeline__item");

        timelineItems.first().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Deputy contact details changed"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by case manager (87654321)"
            );
            cy.wrap(item).contains("time", "14/12/2021 14:41:17");
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(1)",
                "First name: Bob"
            );
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(2)",
                "Surname: Builder"
            );
        });

        timelineItems.next().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "Executive Case Manager changed to ProTeam1 User1"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by case manager (12345678)"
            );
            cy.wrap(item).contains("time", "24/11/2021 11:01:59");
        });

        timelineItems.next().within((item) => {
            cy.wrap(item).contains(
                ".moj-timeline__title",
                "New client added to deputyship"
            );
            cy.wrap(item).contains(
                ".moj-timeline__byline",
                "by system admin (12345678)"
            );
            cy.wrap(item).contains("time", "09/09/2021 14:01:59");
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(1)",
                "Order number: 03305972"
            );
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(2)",
                "Sirius ID: 7000-0000-1995"
            );
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(3)",
                "Order type: pfa"
            );
            cy.wrap(item).contains(
                ".govuk-list > :nth-child(4)",
                "Client: Duke John Fearless"
            );
        });
    });
});
