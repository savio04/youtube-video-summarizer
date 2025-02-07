const bodyElement = document.body

if (localStorage.getItem("theme") === "dark") {
  bodyElement.classList.remove("light-theme");
  bodyElement.classList.add("dark-theme");
}

const inputElement = document.getElementById("url")
const errorElement = document.getElementById("error")
const resultElement = document.getElementById("result")
const videoElement = document.getElementById("videoyt")
const loadingElement = document.getElementById("loading")

function verifyLocalStorage() {
  resultElement.removeAttribute("extenalId")
  const data = localStorage.getItem("extenalId")

  if (data) {
    inputElement.value = `https://www.youtube.com/watch?v=${data}`
    inputElement.setAttribute("disabled", true)

    loading(true)
    poolingGetVideo(data)

    console.log({ data })
  }
}

verifyLocalStorage()

function getExtenalIdFromUrl(url) {
  const [, possibleExternalId] = url.split("v=")

  const [extenalId,] = possibleExternalId.split("&");

  return extenalId
}

function loading(value) {
  if (value) {
    loadingElement.style.display = "flex"
  } else {
    loadingElement.style.display = "none"
  }
}

function getResume() {
  if (inputElement.value.trim() === "") {
    errorElement.innerText = "Url invalida"

    return
  }

  const extenalId = getExtenalIdFromUrl(inputElement.value)

  const alredyResult = resultElement.getAttribute("externalId")

  if (alredyResult && alredyResult === extenalId) return

  localStorage.setItem("extenalId", extenalId)

  inputElement.setAttribute("disabled", true)

  loading(true)

  fetch("http://localhost:8080/v1/videos", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ url: inputElement.value })
  })
    .then(response => {
      if (response.ok) {
        return response.json()
      }
    })
    .then(data => {
      poolingGetVideo(data.payload.external_id)
    })
    .catch(err => console.log({ err }))
}

function poolingGetVideo(extenalId) {
  fetch(`http://localhost:8080/v1/videos/${extenalId}`)
    .then(response => {
      if (response.ok) {
        return response.json()
      }
    })
    .then(data => {
      if (data.payload.status !== "COMPLETED") {
        return setTimeout(() => poolingGetVideo(extenalId), 1000)
      }
      localStorage.removeItem("extenalId")
      inputElement.removeAttribute("disabled")

      const text = formatText(data.payload.summary)
      resultElement.innerHTML = text
      videoElement.src = `https://www.youtube.com/embed/${data.payload.external_id}`
      videoElement.style.display = "block"
      loading(false)
    })
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
    newText += `<p>${group}</p>`;
  }

  return newText
}

function toggleTheme() {
  if (bodyElement.classList.contains("dark-theme")) {
    bodyElement.classList.remove("dark-theme");
    bodyElement.classList.add("light-theme");

    localStorage.setItem("theme", "light")
  } else {
    bodyElement.classList.remove("light-theme");
    bodyElement.classList.add("dark-theme");

    localStorage.setItem("theme", "dark")
  }
}
