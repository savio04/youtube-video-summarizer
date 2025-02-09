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

async function getResume() {
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
    const response = await fetch("https://yt-api.savioaraujogomes.com/v1/videos", {
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
    const response = await fetch(`https://yt-api.savioaraujogomes.com/v1/videos/${extenalId}`)

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
