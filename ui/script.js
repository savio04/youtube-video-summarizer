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

function isYouTubeUrl(url) {
  const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube|youtu|youtube-nocookie)\.(com|be)\/(watch\?v=[a-zA-Z0-9_-]+|(?:v|e(?:mbed)?)\/[a-zA-Z0-9_-]+)(?!.*\/shorts)/i;
  return youtubeRegex.test(url);
}

function getExtenalIdFromUrl(url) {
  const [, possibleExternalId] = url.split("v=")

  const [extenalId,] = possibleExternalId.split("&");

  return extenalId
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

async function getResume() {
  if (inputElement.value.trim() === "") {
    inputElement.classList.add("error")
    errorElement.innerText = "Url obrigatória"

    return
  }

  if (!isYouTubeUrl(inputElement.value)) {
    inputElement.classList.add("error")
    errorElement.innerText = "Url inválida"

    return
  }

  const extenalId = getExtenalIdFromUrl(inputElement.value)

  const alredyResult = resultElement.getAttribute("extenal_id")

  if (alredyResult && alredyResult === extenalId) return

  localStorage.setItem("extenalId", extenalId)

  inputElement.setAttribute("disabled", true)

  loading(true)

  const response = await fetch("http://localhost:8080/v1/videos", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ url: inputElement.value })
  })

  if (!response.ok) {
    return
  }

  const data = await response.json()

  await poolingGetVideo(data.payload.external_id)
}

async function poolingGetVideo(extenalId) {
  const response = await fetch(`http://localhost:8080/v1/videos/${extenalId}`)

  if (!response.ok) {
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
  } else {
    setTimeout(() => poolingGetVideo(extenalId), 1500)
  }
}

function formatText(text) {
  const items = text.split(".")

  const groupSize = 5
  let groupCount = 0

  let newText = ""
  let group = ""

  for (let item of items) {
    item += `${item}.`
    groupCount += 1

    if (groupSize === groupCount) {
      group = `<p>${group}</p>`
      newText += group

      groupCount = 0
      group = ""
    } else {
      group += item
    }
  }

  if (groupCount > 0) {
    newText += `<p>${group}</p>`
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
