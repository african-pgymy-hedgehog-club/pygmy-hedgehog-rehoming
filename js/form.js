"use strict";

/**
 * Convert error to an object for json to stringify
 * @return {object}
 */
Error.prototype.toJSON = function() {
    var alt = {};

    Object.getOwnPropertyNames(this).forEach(function(key) {
        alt[key] = this[key];
    });

    return alt;
};

/**
 * Send error to the back end to be logged
 * @param {object}
 * @return {Promise}
 */
function logError(err) {
    let formData = new FormData();
    formData.append("error", JSON.stringify(err));

    return fetch("api/log", {
        credentials: "same-origin",
        method: "POST",
        body: formData
    }).catch(function(err) {
        console.error(err);
    });
}

/**
 * Try to parse response as json if it fails catch the error and parse as text
 * then display send the error to be logged
 * @throws {error}
 * @return {Promise}
 */
Response.prototype.jsonCatch = function() {
    return this.clone().json().catch(function(err) {
        return this.text().then(function(resErr) {
            UIkit.notify(resErr, { status: "danger", timeout: 3500 });

            return logError({ // Log error on the server
                err: err,
                resErr: resErr
            }).then(function () {
                throw err;
            });
        });
    });
};

/**
 * Empty form inputs
 * @param {object} form
 */
function emptyForm(form) {
    form.find("input, textarea").val("");
    
    let firstSelectVal = form.find("select option").eq(0).val();
    form.find("select").val(firstSelectVal);
}

$(function () {
    $("form").on("submit", function (e) { // On form submit send to api specified in the action attribute
        e.preventDefault();

        let form = $(this);
        let route = $(this).attr("action");
        let routeType = route.split("/").slice(-1)[0];
        let formData = new FormData( $(this)[0] );

        let buttons = $(this).find("button");

        let message;
        let messageTimeout = setTimeout(function() {
            message = UIkit.notify("Working to submit form... <i class='uk-icon-spinner uk-icon-spin'></i>", { status: "message", timeout: 0 });
        }, 1500);
        buttons.attr("disabled", true); // Disable form button

        fetch(route, { // Async post the data to the backend api
            credentials: "same-origin",
            method: "POST",
            body: formData
        }).then(function (response) {
            return response.jsonCatch();
        }).then(function (data) {
            let success = data.success;
            let error = data.error || null;

            if(message) {
                message.close();
            }
            clearTimeout(messageTimeout);

            if(success) {
                if(routeType == "home-for-hog") {
                    routeType = "Find a Home For Your Hog";
                } else if(routeType == "foster-carer") {
                    routeType = "Foster Carers";
                }

                UIkit.notify("Successfully submitted " + routeType + " form", { status: "success", timeout: 0 });
                emptyForm(form);
            }
            else if(error) {
                UIkit.notify(error, { status: "danger", timeout: 3500 });
            }
            else {
                UIkit.notify("Sorry, there was an error", { status: "danger", timeout: 3500 });
            }

            buttons.attr("disabled", false); // Enable form buttons
        }).catch(err => {
            if(message) {
                message.close();
            }
            clearTimeout(messageTimeout);

            UIkit.notify("Sorry, there was an error", { status: "danger", timeout: 3500 });

            logError(err).then(function() {
                buttons.attr("disabled", false); // Enable form buttons
                console.error(err);
            });
        });
    });

    $("button#reset").click(function (e) { // On reset button click, empty form inputs
        e.preventDefault();

        emptyForm( $("form") );
    });
});
