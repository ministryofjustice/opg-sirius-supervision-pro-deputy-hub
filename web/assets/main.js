import "./main.scss";
import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import accessibleAutocomplete from "accessible-autocomplete";

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector(".app-main-class"));

if (document.querySelector("#select-existing-firm")) {
    accessibleAutocomplete.enhanceSelectElement({
        selectElement: document.querySelector("#select-existing-firm"),
        defaultValue: "",
    });
}

if (document.querySelector("#f-existing-firm")) {
    document.getElementById("f-existing-firm").onclick = function () {
        toggleChangeFirmAutoCompleteHideClass();
    };
}

function toggleChangeFirmAutoCompleteHideClass() {
    document
        .getElementById("change-firm-autocomplete")
        .classList.toggle("hide");
}
