window.API_URL = "http://localhost:8080"

const bodyElement = document.body
const toggleSlider = document.querySelector('.toggle-slider i')

if (localStorage.getItem("theme") === "dark") {
  bodyElement.classList.remove("light-theme");
  bodyElement.classList.add("dark-theme");

  toggleSlider.className = "fas fa-moon"
}

const inputElement = document.getElementById("url")
const errorElement = document.getElementById("error")
const resultElement = document.getElementById("result")
const videoElement = document.getElementById("videoyt")
const loadingElement = document.getElementById("loading")
const footerElement = document.getElementById("footer")

function verifyLocalStorage() {
  const data = localStorage.getItem("extenalId")
  resultElement.removeAttribute("extenalId")

  if (data) {
    inputElement.value = `https://www.youtube.com/watch?v=${data}`
    inputElement.setAttribute("disabled", true)

    loading(true)
    poolingGetVideo(data)
  }
}

verifyLocalStorage()

function extractYouTubeVideoId(url) {
  const regex =
    /(?:youtube\.com\/(?:[^/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?/ ]{11})/;

  const match = url.match(regex);

  return match ? match[1] : null;
}

function loading(value) {
  if (value) {
    videoElement.style.display = "none"
    resultElement.style.display = "none"
    footerElement.style.display = "none"
    loadingElement.style.display = "flex"
  } else {
    loadingElement.style.display = "none"
  }
}

function setCookie(name, value) {
  const date = new Date()

  date.setTime(date.getTime() + (2 * 60 * 1000))

  const expires = `expires=${date.toUTCString()}`

  document.cookie = `${name}=${value}; ${expires}; path=/`
}

function getCookie(name) {
  const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'))

  if (match) return match[2]

  return null
}

function clearCookie(name) {
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/`
}

async function getToken() {
  return new Promise((resolve, reject) => {
    grecaptcha.ready(function() {
      grecaptcha.execute('6LfifNcqAAAAAJ0Ztlfa0TUgp5rqLYh9y8LAt_ab', { action: 'submit' }).then(function(token) {
        resolve(token)
      })
        .catch(error => reject(error))
    })
  })
}

async function verifyToken() {
  let token = getCookie("token")

  if (!token) {
    try {
      token = await getToken()

      console.log({ token })
    } catch (error) {
      window.alert(`Erro ao executar o reCAPTCHA: ${error}`)
      throw error
    }
  }

  //Verify in backend
  try {
    await fetch(`${window.API_URL}/v1/reptcha`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ token })
    })

    setCookie("token", token)
  } catch (error) {
    window.alert(`Token invalido`)
    clearCookie("token")
    throw error
  }

}

async function getResume() {
  await verifyToken()

  if (inputElement.value.trim() === "") {
    inputElement.classList.add("error")
    errorElement.innerText = "Url obrigatória"

    return
  }

  const externalId = extractYouTubeVideoId(inputElement.value)

  if (!externalId) {
    inputElement.classList.add("error")
    errorElement.innerText = "Url inválida"

    return
  }

  const alredyResult = resultElement.getAttribute("extenal_id")

  if (alredyResult && alredyResult === externalId) return

  localStorage.setItem("extenalId", externalId)

  inputElement.setAttribute("disabled", true)

  loading(true)

  try {
    const response = await fetch(`${window.API_URL}/v1/videos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ url: inputElement.value, externalId })
    })

    if (!response.ok) {
      localStorage.removeItem("extenalId")
      loading(false)
      window.alert("Falha ao enviar dados")
      return
    }

    const data = await response.json()

    await poolingGetVideo(data.payload.external_id)
  } catch (error) {
    localStorage.removeItem("extenalId")
    loading(false)
    window.alert("Falha ao enviar dados")
  }
}

async function poolingGetVideo(extenalId) {
  try {
    const response = await fetch(`${window.API_URL}/v1/videos/${extenalId}`)

    if (!response.ok) {
      localStorage.removeItem("extenalId")
      loading(false)
      window.alert("Falha ao enviar dados")

      return
    }

    const data = await response.json()

    if (data.payload.status === "COMPLETED") {
      localStorage.removeItem("extenalId")
      inputElement.removeAttribute("disabled")

      const text = formatText(data.payload.summary)

      resultElement.innerHTML = `<h2>Resumo</h2> ${text}`
      resultElement.setAttribute("extenal_id", data.payload.external_id)
      resultElement.style.display = "flex"
      videoElement.src = `https://www.youtube.com/embed/${data.payload.external_id}`
      videoElement.style.display = "block"
      footerElement.style.display = "block"

      loading(false)
    } else if (data.payload.status === "FAILED") {
      loading(false)
      localStorage.removeItem("extenalId")
      inputElement.removeAttribute("disabled")

      window.alert("Falha ao fazer download do video, tente novamente! (Ou tente outra url)")

    } else {
      setTimeout(() => poolingGetVideo(extenalId), 1500)
    }
  } catch (error) {
    localStorage.removeItem("extenalId")
    loading(false)
    window.alert("Falha ao enviar dados")
  }
}

function formatText(text) {
  const items = text.split("\n\n")

  let newText = ""

  for (const item of items) {
    newText += `<p>${item}</p>`
  }

  return newText
}

function toggleTheme() {
  if (bodyElement.classList.contains("dark-theme")) {
    bodyElement.classList.remove("dark-theme")
    bodyElement.classList.add("light-theme")

    toggleSlider.className = "fas fa-sun"

    localStorage.setItem("theme", "light")
  } else {
    bodyElement.classList.remove("light-theme")
    bodyElement.classList.add("dark-theme")

    toggleSlider.className = "fas fa-moon"


    localStorage.setItem("theme", "dark")
  }
}

function clearInputError() {
  inputElement.classList.remove("error")
  errorElement.innerText = ""
}


//RECAPCHA
// function onClick(e) {
//   e.preventDefault();
//   grecaptcha.ready(function() {
//     grecaptcha.execute('6LfifNcqAAAAAJ0Ztlfa0TUgp5rqLYh9y8LAt_ab', { action: 'submit' }).then(function(token) {
//       // Add your logic to submit to your backend server here.
//       console.log({ token })
//     });
//   });
// }
