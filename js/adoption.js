"use strict";

var hedgehogName = window.location.pathname.split("/")[2] || null;

function displayHedgehog(name) {
    $("div#" + name).removeClass("uk-hidden");
    $("div#adoption-details").addClass("uk-hidden");
}

function displayDetails() {
    $("div#adoption-details").removeClass("uk-hidden");
    $("div.hedgehog").addClass("uk-hidden");
}

function stateChange(e) {
    var state = e.state;

    if(state && state.hasOwnProperty("name")) {
        displayHedgehog(state.name);
    } else {
        displayDetails();
    }
}

window.onpopstate = stateChange;

$(document).ready(function () {
    if(hedgehogName) {
        displayHedgehog(hedgehogName);
    }

    $("a[name]").click(function (e) {
        e.preventDefault();

        var name = $(this).attr("name");

        history.pushState({name: name}, "Adoption", "adoption/" + name);
        displayHedgehog(name);
    });
});
