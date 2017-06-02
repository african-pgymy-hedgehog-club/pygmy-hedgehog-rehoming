"use strict";

/**
 * FormValidator class to validate a form
 */
class FormValidator {

    /**
     * @param jqueryObject form
     * @param array inputs Format: ['name', 'email']
     */
    constructor(form, inputs) {
        this.form = form;
        this.inputs = {};
        inputs.forEach(input => {
            this.inputs[input] = {
                required: false
            };
        });
    }

    /**
     * Set inputs to required
     * @param array inputs
     */
    required(inputs = null) {
        if(inputs !== null) {
            inputs.forEach(input => {
                if(!Object.hasOwnProperty(input)) {
                    this.inputs[input] = {}
                }
                
                this.inputs[input].required = true;
            });
        }
        else { // If no inputs are passed, set all inputs passed in constructor as required
            Object.keys(this.inputs).forEach(input => {
                this.inputs[input].required = true;
            });
        }
    }

    /**
     * Set input to be matched against postcode regex
     * @param array inputs
     */
    postcode(inputs) {
        const regex = /^([a-zA-Z]){1,2}([0-9]){1,2}(\s?){1}([0-9]){1}([a-zA-Z]){2}$/
        inputs.forEach(input => {
            if(!this.inputs.hasOwnProperty(input)) {
                this.inputs[input] = {required: false};
            }

            this.inputs[input].postcode = regex;
        });
    }

    /**
     * Highlight form element to show it it's value invalid
     * @param jqueryObject element
     * @param options
     */
     invalid(element, options) {
         element.addClass('fv-invalid');
         let nodeName = element[0].nodeName.toLowerCase();

         if(nodeName === 'input' && element.attr('type') !== 'file') {
             element.bind('keyup change mouseup', function () {
                 let valid = true;
                 for (var option in options) { // For each option set on element check if they match the requirements
                     if(option == 'required' && options[option]) {
                         if(element.val() === '' || element.val() === null || Number(element.val()) === 0) {
                             valid = false;

                             break;
                         }
                     }
                     else {
                         let regex = options[option];
                         if(regex && !element.val().match(regex) && element.val() !== '' && Number(element.val()) === 0) {
                             valid = false;

                             break;
                         }
                     }
                 }

                 if(valid) {
                     $(this).removeClass('fv-invalid');
                 }
                 else {
                     $(this).addClass('fv-invalid');
                 }
             });
         }
         else {
             element.bind('keyup change mouseup', function () {
                 let valid = true;
                 for (var option in options) { // For each option set on element check if they match the requirements
                     if(option == 'required' && options[option]) {
                         if(element.val() === '' || element.val() === null || Number(element.val()) === 0) {
                             valid = false;

                             break;
                         }
                     }
                     else {
                         let regex = options[option];
                         if(regex && !element.val().match(regex)) {
                             valid = false;

                             break;
                         }
                     }
                 }

                if(valid) {
                     $(this).removeClass('fv-invalid');
                }
                else {
                     $(this).addClass('fv-invalid');
                }

                if(nodeName === 'select') {
                    let nextElement = $(this).next();
                    let nextClass = nextElement.attr('class') || '';
                    if(nextClass.includes('CaptionCont SelectBox')) {
                        if(valid) {
                            nextElement.removeClass('fv-invalid');
                        }
                        else {
                            nextElement.addClass('fv-invalid');
                        }
                    }
                }
            });

            if(nodeName === 'select') {
                let nextElement = element.next();
                let nextClass = nextElement.attr('class') || '';
                if(nextClass.includes('CaptionCont SelectBox')) {
                    nextElement.addClass('fv-invalid');
                }
            }
        }

        element.addClass('fv-invalid');
    }

    /**
     * Validate form
     * @return bool
     */
    validate() {
        let valid = true;

        Object.keys(this.inputs).forEach(input => { // For each input object key check if the input element meets the validation
            let options = this.inputs[input];
            let element = null;
            let elements = null;
            let inputs = null;

            // If the input name has a two names with a | character splitting then treat as an or for validation
            let pattern = /([\w\d]+)\|([\w\d]+)/;
            if(!input.match(pattern)) {
                element = this.form.find(`select[name="${input}"]`);

                if(element.length === 0) {
                    element = this.form.find(`[name="${input}"]`);
                }
            }
            else {
                inputs = input.split('|');
                elements = inputs.map( input => this.form.find(`[name="${input}"]`) );
            }

            for (var option in options) { // For each option set on element check if they match the requirements
                if(option == 'required' && options[option]) {
                    if(elements && Array.isArray(elements)) {
                        elements.forEach(element => element.removeClass('fv-invalid')); // Set or elements to valid

                        elements = elements.filter(element => ( element.val() === '' || element.val() === null || Number(element.val()) === 0 ));

                        if(elements.length === inputs.length) {
                            elements.forEach(element => this.invalid(element, option));

                            valid = false;
                            break;
                        }
                    }
                    else {
                        if(element.val() === '' || element.val() === null || Number(element.val()) === 0) {
                            this.invalid(element, options);

                            valid = false;
                            break;
                        }
                    }
                }
                else {
                    let regex = options[option];
                    if(regex && !element.val().match(regex) && element.val() !== '') {
                        this.invalid(element, options);
                        valid = false;

                        break;
                    }
                }
            }

            if(valid && element) {
                element.removeClass('fv-invalid');
            }
        });

        return valid;
    }
}
