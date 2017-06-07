"use strict";

var state = history.state;

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

$(document).ready(function () {
    window.onpopstate = stateChange;


    $("a[name]").click(function (e) {
        e.preventDefault();

        var name = $(this).attr("name");

        history.pushState({name: name}, "Adoption", "adoption/" + name);
        displayHedgehog(name);
    });
});
