function getResume() {
  const inputElement = document.getElementById("url")

  if (inputElement.value.trim() === "") {
    const errorElement = document.getElementById("error")
    errorElement.innerText = "Url invalida"

    return
  }

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
      if (data.payload.status !== "COMPLETED") {
        return poolingGetVideo(data.payload.external_id)
      }

      const resultElement = document.getElementById("result")
      const videoElement = document.getElementById("videoyt")

      const text = formatText(data.payload.summary)
      resultElement.innerHTML = `<div>${text}</div>`
      videoElement.src = `https://www.youtube.com/embed/${data.payload.external_id}`
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

      const resultElement = document.getElementById("result")
      const videoElement = document.getElementById("videoyt")

      resultElement.innerHTML = `<p>${data.payload.summary}</p>`
      videoElement.src = `https://www.youtube.com/embed/${data.payload.external_id}`
    })
}

function formatText(text) {
  const items = text.split(".")

  const groupSize = 6
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
