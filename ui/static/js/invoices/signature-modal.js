document.addEventListener("DOMContentLoaded", () => {
    const openBtn = document.querySelector(".open-signature-btn")
    const submitBtn = document.querySelector(".submit-agreement-btn")
    const agreementCheckbox = document.getElementById("agreement-checkbox")

    const modal = document.getElementById("signature-modal")
    const cancelBtn = document.getElementById("cancel-signature")
    const canvas = document.getElementById("signature-canvas")
    const clearBtn = document.getElementById("clear-signature")
    const confirmBtn = document.getElementById("confirm-signature")

    const fileInput = document.getElementById("signature-file")

    const previewContainer = document.querySelector(".signature-preview")
    const previewImg = document.getElementById("signature-preview-img")

    openBtn.disabled = true
    submitBtn.disabled = true

    function resizeCanvas() {
        const ratio = Math.max(window.devicePixelRatio || 1, 1)
        const ctx = canvas.getContext("2d")

        canvas.width = canvas.offsetWidth * ratio
        canvas.height = canvas.offsetHeight * ratio

        ctx.setTransform(ratio, 0, 0, ratio, 0, 0)
    }

    resizeCanvas()

    const signaturePad = new SignaturePad(canvas, {
        backgroundColor: "rgb(255,255,255)",
        penColor: "black",
    })

    agreementCheckbox.addEventListener("change", () => {
        openBtn.disabled = !agreementCheckbox.checked
    })

    openBtn.addEventListener("click", () => {
        modal.classList.remove("hidden")
        resizeCanvas()
    })

    cancelBtn.addEventListener("click", () => {
        modal.classList.add("hidden")
    })

    clearBtn.addEventListener("click", () => {
        signaturePad.clear()
    })

    confirmBtn.addEventListener("click", () => {
        if (signaturePad.isEmpty()) {
            alert("Please provide a signature first.")
            return
        }

        canvas.toBlob((blob) => {
            const file = new File([blob], "signature.png", {
                type: "image/png",
            })

            const dataTransfer = new DataTransfer()
            dataTransfer.items.add(file)
            fileInput.files = dataTransfer.files
            previewImg.src = URL.createObjectURL(blob)
            previewContainer.hidden = false

            submitBtn.disabled = false

            modal.classList.add("hidden")
        })
    })

    submitBtn.addEventListener("click", (e) => {
        if (!fileInput.files.length) {
            e.preventDefault()
            alert("You must add a signature before submitting.")
        }
    })
})
