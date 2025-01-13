async function submit(e) {
    e.preventDefault();

    let username = document.getElementById("input-username").value;
    let password = document.getElementById("input-password").value;

    // todo validate input

    await fetch("/auth/login", {
        method: 'POST',
        body: JSON.stringify({
            username: username,
            password: password
        }),
    })
        .then(r => r.json())
        .then(data => {
            window.location.replace(data['redirect_uri']);
        })
        .catch(err => document.getElementById("error-container").innerHTML = err); // todo
}

window.addEventListener('load', function () {
    document.getElementById("form-login").addEventListener("submit", submit);
})