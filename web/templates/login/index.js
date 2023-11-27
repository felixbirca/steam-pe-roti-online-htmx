const firebaseConfig = {
    apiKey: "AIzaSyDuZF8RF-xGHJuMd7dOjZ5W-erVFgxJFx8",
    authDomain: "stem-pe-roti-online.firebaseapp.com",
    projectId: "stem-pe-roti-online",
    storageBucket: "stem-pe-roti-online.appspot.com",
    messagingSenderId: "261534929781",
    appId: "1:261534929781:web:87bdfb838c958d9f6f793c",
    measurementId: "G-HVQ9MW0LB3"
};

firebase.initializeApp(firebaseConfig);

var provider = new firebase.auth.GoogleAuthProvider();


function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }


const loginBtn = document.getElementById("loginButton");
console.log(loginBtn);

firebase.auth().setPersistence(firebase.auth.Auth.Persistence.NONE);

async function postIdTokenToSessionLogin(url, idToken, csrfToken) {
    await fetch(`${url}`, {
        method: "POST",
        body: JSON.stringify({ idToken: idToken, csrfToken: csrfToken })
    }).then(result => {
        if (result.status == 200) {
            window.location.replace("/secured")
        }
    });
}


loginBtn.addEventListener("click", async () => {
    firebase.auth()
        .signInWithPopup(provider)
        .then((result) => {
            try {
                console.log("function?", postIdTokenToSessionLogin);
            }
            catch (e) {
                console.log("function?error", e);
            }
            console.log("result", result)
            /** @type {firebase.auth.OAuthCredential} */
            var credential = result.credential;

            // This gives you a Google Access Token. You can use it to access the Google API.
            var token = credential.accessToken;
            // The signed-in user info.
            var user = result.user;
            console.log("function?", postIdTokenToSessionLogin);

            // IdP data available in result.additionalUserInfo.profile.
            // ...
            return user.getIdToken().then(idToken => {
                console.log("idToken", idToken);
                console.log(getCookie);
                // Session login endpoint is queried and the session cookie is set.
                // CSRF protection should be taken into account.
                // ...
                const csrfToken = getCookie('csrfToken')
                console.log(postIdTokenToSessionLogin);
                return postIdTokenToSessionLogin('/sessionLogin', idToken, csrfToken);
            });
        }).catch((error) => {
            // Handle Errors here.
            var errorCode = error.code;
            var errorMessage = error.message;
            // The email of the user's account used.
            var email = error.email;
            // The firebase.auth.AuthCredential type that was used.
            var credential = error.credential;
            // ...
        });
})

