describe("Pro Deputy Hub", () => {
  beforeEach(() => {
      cy.setCookie("Other", "other");
      cy.setCookie("XSRF-TOKEN", "abcde");
      cy.visit("/supervision/deputies/professional/");
  });

    it('finds the content "Hello world!"', () => {
        cy.contains('Hello world!')
    })
});
