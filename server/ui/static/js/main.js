document.addEventListener('DOMContentLoaded', function () {
    const flashMessage = document.getElementById('flashMessage');
    if (flashMessage) {
        setTimeout(function () {
            flashMessage.classList.add('fade-out');
            setTimeout(function () {
                flashMessage.style.display = 'none';
            }, 300);
        }, 3000);
    }

    // const emailInput = document.getElementById('email');
    // if (emailInput) {
    //     emailInput.disabled = true;
    // }
});