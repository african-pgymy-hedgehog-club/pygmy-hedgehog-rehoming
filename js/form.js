"use strict";

/**
 * Convert error to an object for json to stringify
 * @return {object}
 */
Error.prototype.toJSON = function() {
    var alt = {};

    Object.getOwnPropertyNames(this).forEach((key) => {
        alt[key] = this[key];
    });

    return alt;
};

/**
 * Send error to the back end to be logged
 * @param {object}
 * @return {Promise}
 */
const logError = async (err) => {
    let formData = new FormData();
    formData.append("error", JSON.stringify(err));

    return fetch("api/log", {
        credentials: "same-origin",
        method: "POST",
        body: formData
    }).catch(err => console.error(err));
};

/**
 * Try to parse response as json if it fails catch the error and parse as text
 * then display send the error to be logged
 * @throws {error}
 * @return {Promise}
 */
Response.prototype.jsonCatch = async function() {
    return this.clone().json().catch(err => {
        return this.text().then(async resErr => {
            UIkit.notify(resErr, { status: "danger", timeout: 3500 });

            await logError({ // Log error on the server
                err,
                resErr
            });

            throw err;
        });
    });
};

/**
 * Empty form inputs
 * @param {object} form
 */
const emptyForm = (form) => {
    form.find("input, textarea").val("");
};

$(document).ready(function () {
    $("form").submit(async function (e) { // On form submit send to api specified in the action attribute
        e.preventDefault();

        let route = $(this).attr("action");
        let [routeType] = route.split("/").slice(-1);
        let formData = new FormData( $(this)[0] );

        let buttons = $(this).find("button");

        let message;
        let messageTimeout = setTimeout(() => {
            message = UIkit.notify("Working to submit form... <i class='uk-icon-spinner uk-icon-spin'></i>", { status: "message", timeout: 0 });
        }, 1500);
        try {
            buttons.attr("disabled", true); // Disable form button

            let response = await fetch(route, { // Async post the data to the backend api
                credentials: "same-origin",
                method: "POST",
                body: formData
            });

            let {success, error} = await response.jsonCatch(); // Retrieve the response as json

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

                UIkit.notify(`Successfully submitted ${routeType} form`, { status: "success", timeout: 0 });
                emptyForm( $(this) );
            }
            else if(error) {
                UIkit.notify(error, { status: "danger", timeout: 3500 });
            }
            else {
                UIkit.notify("Sorry, there was an error", { status: "danger", timeout: 3500 });
            }

            buttons.attr("disabled", false); // Enable form buttons
        }
        catch(err) {
            if(message) {
                message.close();
            }
            clearTimeout(messageTimeout);

            UIkit.notify("Sorry, there was an error", { status: "danger", timeout: 3500 });

            await logError(err);
            buttons.attr("disabled", false); // Enable form buttons
            console.error(err);
        }
    });

    $("button#reset").click(function (e) { // On reset button click, empty form inputs
        e.preventDefault();

        emptyForm( $("form") );
    });
});
