document.addEventListener('DOMContentLoaded', function() {
    setTimeout(function() {
        var messageContainer = document.getElementById('message-container');
        var captchaErrorContainer = document.getElementById('captcha-error-container');
        if (messageContainer) {
            messageContainer.style.display='none';
        }
        if (captchaErrorContainer) {
            captchaErrorContainer.style.visibility='hidden';
        }
    }, 4000);
});
