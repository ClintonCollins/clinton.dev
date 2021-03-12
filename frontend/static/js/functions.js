const navLinks = document.querySelectorAll('.scroll-link');
const contactForm = document.getElementById('contactForm')
const contactInputs = document.querySelectorAll('.contact-input')
const middleScreen = screen.height / 4
const contactError = document.getElementById('contactError')
const contactSuccess = document.getElementById('contactSuccess')
const contactFormNameEl = document.getElementById('contactName')
const contactFormEmailEl = document.getElementById('contactEmail')
const contactFormMessageEl = document.getElementById('contactMessage')
const contactFormVerificationSumEl = document.getElementById('contactFormVerificationSum')

const homeLink = document.getElementById('home-link')
if (window.location.pathname === '/') {
    homeLink.classList.add('active-link')
}

document.addEventListener('DOMContentLoaded', function () {
    window.addEventListener('scroll', debounce(checkHeaderScrollPosition, 5))
    let homeLink = document.getElementById('home-link')
    homeLink.addEventListener('click', (event) => {
        if (window.location.pathname === '/') {
            event.preventDefault()
            window.scrollTo(0, 0)
            history.pushState(null, null, '/')
        }
    })
    if (window.fetch)  {
        contactForm.addEventListener('submit', submitForm, true)
    }
    for (let i = 0; i < contactInputs.length; i++) {
        let input = contactInputs[i]
        input.addEventListener('keyup', event => {
            debounce(validateInput(event.target), 750)
        })
        input.addEventListener('blur', event => {
            validateInput(event.target)
        })
    }
});

function checkHeaderScrollPosition() {
    let header = document.getElementById('header')
    let scrollTop = window.scrollY
    if (scrollTop >= 10) {
        header.classList.add('scroll-on')
    } else {
        header.classList.remove('scroll-on')
    }

    for (let i = 0; i < navLinks.length; i++) {
        let navLink = navLinks[i]
        if (navLink.hash !== "") {
            let section = document.querySelector(navLink.hash)
            let secOffset = section.offsetTop
            secOffset -= middleScreen
            if (secOffset <= scrollTop && secOffset + section.offsetHeight > scrollTop) {
                navLink.classList.add('active-link')
            } else {
                navLink.classList.remove('active-link')
            }
        } else {
            let section = document.querySelector('#home')
            let secOffset = section.offsetTop
            secOffset -= middleScreen
            if (secOffset <= scrollTop && secOffset + section.offsetHeight > scrollTop) {
                navLink.classList.add('active-link')
            } else {
                navLink.classList.remove('active-link')
            }
        }
    }

}

function debounce(func, waitMS) {
    let timeout;
    return function () {
        let context = this
        let args = arguments;
        let delayedFunction = function () {
            timeout = null;
            func.apply(context, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(delayedFunction, waitMS);
    };
}

function submitSuccess() {
    contactError.classList.add('hidden')
    contactSuccess.classList.remove('hidden')
}

function submitError(errors) {
    let formErrors = document.getElementById('formErrors')
    formErrors.innerHTML = ''
    for (let i = 0; i < errors.length; i++) {
        let errLi = document.createElement('li')
        errLi.innerHTML = errors[i]
        formErrors.appendChild(errLi)
    }
    contactError.classList.remove('hidden')
    contactSuccess.classList.add('hidden')
}

function validateInput(input) {
    let validity = true
    if (input.type === 'email') {
        // Test email inputs only.
        validity = /\S+@\S+\.\S+/.test(input.value);
    } else {
        // Test all other inputs to make sure they're not pure whitespace.
        validity = /\S/.test(input.value);
    }
    if (validity) {
        input.classList.remove('is-invalid')
        input.classList.add('is-valid')
    } else {
        input.classList.remove('is-valid')
        input.classList.add('is-invalid')
    }
    return validity
}

function validateInputs(contactName, contactEmail, contactMessage) {
    let validity = true
    if (!validateInput(contactName)) {
        validity = false
    }
    if (!validateInput(contactEmail)) {
        validity = false
    }
    if (!validateInput(contactMessage)) {
        validity = false
    }
    return validity
}

function disableContactFormInputs() {
    let contactInputs = document.getElementsByClassName('contact-input')
    let submitButton = document.getElementById('contactSubmit')
    for (let i = 0; i < contactInputs.length; i++) {
        let element = contactInputs[i]
        element.disabled = true
    }
    submitButton.classList.add('disabled-button')
}

function enableContactFormInputs() {
    let contactInputs = document.getElementsByClassName('contact-input')
    let submitButton = document.getElementById('contactSubmit')
    for (let i = 0; i < contactInputs.length; i++) {
        let element = contactInputs[i]
        element.disabled = false
    }
    submitButton.classList.remove('disabled-button')
}

function clearContactFormInputs() {
    contactFormNameEl.classList.remove('is-valid')
    contactFormEmailEl.classList.remove('is-valid')
    contactFormMessageEl.classList.remove('is-valid')
    contactFormVerificationSumEl.classList.remove('is-valid')
    contactForm.reset()
}

function submitForm(event) {
    event.preventDefault()
    let formData = new FormData(contactForm)
    let contactToken = document.getElementById('contactFormToken')
    let contactVerificationSumLabel = document.getElementById('contactFormVerificationSumLabel')
    if (!validateInputs(contactFormNameEl, contactFormEmailEl, contactFormMessageEl)) {
        return
    }
    disableContactFormInputs()
    let data = new URLSearchParams()
    for (const keyValue of formData) {
        data.append(keyValue[0], keyValue[1].toString());
        // console.log(`${keyValue[0]} = ${keyValue[1].toString()}`)
    }
    fetch('/api/contact', {
        method: 'POST',
        body: data,
    }).then(response => {
        if (response.status !== 200) {
            submitError([response.statusText])
            throw Error(response.statusText)
        }
        enableContactFormInputs()
        return response.json()
    }).then(data => {
        if (data.success) {
            submitSuccess()
            let formSuccessMessageEl = document.getElementById('successMessage')
            formSuccessMessageEl.innerHTML = `Thank you ${contactFormNameEl.value},
             I'll be in contact with your shortly!`
            clearContactFormInputs()
        } else {
            // Got an error.
            submitError(data.errors)
        }
        // Update with a new token and verification question.
        contactToken.value = data.new_token
        contactVerificationSumLabel.innerHTML = `What is ${data.new_form_num_1} + ${data.new_form_num_2} ?`
        enableContactFormInputs()
    }).catch(err => {
        submitError([err])
        console.error(err)
        enableContactFormInputs()
    })
}
