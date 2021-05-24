$(function() {
    // Set max dates
    var maxDOB = dateFns.subYears(new Date(), 16);
    var dobEl = $("input[name='dob']");
    dobEl.attr("max", dateFns.format(maxDOB, "YYYY-MM-DD"))


    // Check if other value is selected for own or rent and if so make the details text input required
    $("input[name='own_or_rent']").on('change', function() {
        var selected = $(this).val();
        var otherHomeDetailsEl = $("input[name='other_home_details']");

        otherHomeDetailsEl.removeAttr('required', true);
        if (selected === "other") {
            otherHomeDetailsEl.attr('required', true);
        }
    });

    // Check if can you drive is no, is so show can someone you know drive
    $("input[name='can_you_drive']").on('change', function() {
        var selected = $(this).val();
        var canYouDriveNoQuestionEL = $("#can-you-drive-no");

        if (selected === "No") {
            canYouDriveNoQuestionEL.removeClass("uk-hidden");
            canYouDriveNoQuestionEL.find("input").attr("required", true)
        } else {
            canYouDriveNoQuestionEL.addClass("uk-hidden");
            canYouDriveNoQuestionEL.find("input").attr("required", false)
        }
    });

    // Set date to today
    $("input[name='dated']").val(dateFns.format(new Date(), "YYYY-MM-DD"))
});