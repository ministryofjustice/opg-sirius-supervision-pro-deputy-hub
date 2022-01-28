describe("Pro Deputy Hub", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/professional/deputy/1");
    });

    describe("Header", () => {
        it("shows opg sirius within banner", () => {
            cy.contains(".moj-header__link", "OPG");
            cy.contains(".moj-header__link", "Sirius");
        });

        const expected = ["Supervision", "LPA", "Admin", "Logout"];

        it("has working nav links within header banner", () => {
            cy.get(".moj-header__navigation-list")
                .children()
                .each(($el, index) => {
                    cy.wrap($el).should("contain", expected[index]);
                    let $linkName = expected[index].toLowerCase();
                    cy.wrap($el)
                        .find("a")
                        .should("have.attr", "href")
                        .and("contain", `/${$linkName}`);
                });
        });
    });

    describe("Footer", () => {
        it("the footer should contain a link to the open government licence", () => {
            cy.get(
                ".govuk-footer__licence-description > .govuk-footer__link"
            ).should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/"
            );
        });

        it("the nav link should contain the crown copyright logo", () => {
            cy.get(".govuk-footer__copyright-logo").should(
                "have.attr",
                "href",
                "https://www.nationalarchives.gov.uk/information-management/re-using-public-sector-information/uk-government-licensing-framework/crown-copyright/"
            );
        });
    });

    describe("Pro Person deputy details", () => {
        it("the page should contain the deputy name", () => {
            cy.contains(".hook_header_deputy_name", "firstname surname");
        });

        it("the page should contain the deputy status", () => {
            cy.contains(".hook_header_deputy_status_person", "Active");
        });

        it("the page should contain the firm", () => {
            cy.contains(".hook_header_firm_name", "This is the Firm Name");
        });

        it("the page should contain the deputy number", () => {
            cy.contains(".hook_header_deputy_number", "Deputy Number: 1000");
        });

        it("the page should contain the executive case manager", () => {
            cy.contains(".hook_header_ecm", "Executive Case Manager: displayName");
        });

        it("the page should contain the address without the firm name", () => {
            const expected = [
                "addressLine1",
                "",
                "addressLine2",
                "",
                "addressLine3",
                "",
                "town",
                "",
                "county",
                "",
                "postcode",
            ];
            cy.get(".hook_deputy_address").each((el, i) => {
                // address rows are separated by <br/> tags
                if (i % 2 !== 0) {
                    cy.wrap(el).contains(expected[i]);
                }
            });
        });
    });

    describe("Pro Organisation details", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/professional/deputy/2");
        });

        it("the page should contain the deputy name", () => {
            cy.contains(".hook_header_organisation_name", "Organisation Ltd");
        });

        it("the page should contain the deputy status", () => {
            cy.contains(".hook_header_deputy_status_organisation", "Inactive");
        });
    });
});
