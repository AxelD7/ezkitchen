function setupAccordion(className, displayType) {
    let items = document.getElementsByClassName(className)

    for (let i = 0; i < items.length; i++) {
        const item = items[i]
        const panel = item.nextElementSibling

        if (!panel) continue

        const key = className + "-" + i

        const savedState = sessionStorage.getItem(key)
        if (savedState === "open") {
            item.classList.add("active")
            panel.style.display = displayType
        } else {
            panel.style.display = "none"
        }

        item.addEventListener("click", function () {
            const isOpen = panel.style.display === "none"
            panel.style.display = isOpen ? displayType : "none"
            this.classList.toggle("active", isOpen)
            sessionStorage.setItem(key, isOpen ? "open" : "closed")
        })
    }
}

function setupAddProductBtn() {
    let addProductBtn = document.getElementsByClassName("add-product-btn")
    for (let i = 0; i < addProductBtn.length; i++) {
        addProductBtn[i].addEventListener("click", async function () {
            let category = this.id
            let subcat = this.closest(".subcategory-panel").id
            let color = ""
            if (category == subcat) {
                subcat = ""
            }
            let url = `/product/get/?category=${category}&subcategory=${subcat}&color=${color}`
            console.log(addProductBtn[i].id)
            try {
                const response = await fetch(url)
                if (!response.ok) {
                    throw new Error(`Response Status: ${response.status}`)
                }
                const modalHTML = await response.text()
                console.log(category)
                document.querySelector(".product-modal").innerHTML = modalHTML
                document.querySelector(".product-modal").showModal()

                setupAddProductToEstimateBtn()
            } catch (error) {
                console.error(error.message)
            }
        })
    }
    document.addEventListener("click", (e) => {
        if (e.target.closest("#modal-close-btn")) {
            const modal = document.querySelector(".product-modal")
            modal.close()
        }
    })
}

function setupAddProductToEstimateBtn() {
    let addItemBtn = document.getElementsByClassName("add-item-btn")
    for (let i = 0; i < addItemBtn.length; i++) {
        addItemBtn[i].addEventListener("click", async function () {
            const estimateID =
                document.querySelector(".estimate-ID").dataset.estimateId
            const quantity = this.closest(".card").querySelector(
                ".card-product-quantity",
            ).value
            const productID = parseInt(this.id)

            if (quantity < 1 || isNaN(quantity)) {
                showInlineError(this, "Quantity must be at least 1.")
                return
            }
            try {
                const response = await fetch(`/estimate/${estimateID}/items/`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        product_id: productID,
                        quantity: parseInt(quantity),
                    }),
                })

                if (!response.ok) {
                    const data = await response.json().catch(() => ({}))
                    if (data.errors) {
                        handleProductErrors(this, data.errors)
                        return
                    }
                    throw new Error(`Response Status: ${response.status}`)
                }

                location.reload()
            } catch (error) {
                alert("Something went wrong while adding the product.")
                console.error(error.message)
            }
        })
    }
}

function setupUpdateItemBtn() {
    let saveItemBtn = document.getElementsByClassName("update-item-btn")

    for (let i = 0; i < saveItemBtn.length; i++) {
        saveItemBtn[i].addEventListener("click", async function () {
            currentRow = this.closest("tr")
            lineItemID = currentRow.dataset.lineItemId
            newQuantity = currentRow.querySelector(
                ".estimate-item-quantity",
            ).value
            if (newQuantity < 1 || isNaN(newQuantity)) {
                alert("Quantity must be at least 1")
                return
            }
            try {
                const response = await fetch(`/estimate/items/${lineItemID}`, {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        quantity: parseInt(newQuantity),
                    }),
                })
                if (!response.ok) {
                    throw new Error(`Response status: ${response.status}`)
                }
                location.reload()
            } catch {
                console.Error(error.message)
            }
        })
    }
}

function setupDeleteBtn() {
    let delItemBtn = document.getElementsByClassName("delete-item-btn")
    for (let i = 0; i < delItemBtn.length; i++) {
        delItemBtn[i].addEventListener("click", async function () {
            lineItemID = this.closest("tr").dataset.lineItemId
            try {
                const response = await fetch(`/estimate/items/${lineItemID}`, {
                    method: "DELETE",
                })
                if (!response.ok) {
                    throw new Error(`Response status: ${response.status}`)
                }
                location.reload()
            } catch {
                console.Error(error.message)
            }
        })
    }
}
function handleProductErrors(button, errors) {
    document.querySelectorAll(".product-error").forEach((e) => e.remove())

    if (errors.quantity) {
        showInlineError(button, errors.quantity)
    }

    if (errors.product) {
        showInlineError(button, errors.product)
    }
}

function showInlineError(button, message) {
    const quantityInput = button
        .closest(".card")
        .querySelector(".card-product-quantity")
    const err = document.createElement("p")
    err.className = "product-error"
    err.style.color = "red"
    err.style.margin = "5px 0 0 0"
    err.textContent = message
    quantityInput.insertAdjacentElement("afterend", err)
}

setupAccordion("accordion", "block")
setupAccordion("subcategory-accordion", "flex")
setupAddProductBtn()
setupUpdateItemBtn()
setupDeleteBtn()
