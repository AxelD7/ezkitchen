// Initializes accordion sections and remembers their open/closed state using sessionStorage.
// Takes in the class name of the accordian and the CSS displayType
function setupAccordion(className, displayType) {
    let items = document.getElementsByClassName(className)

    for (let i = 0; i < items.length; i++) {
        const item = items[i]
        const panel = item.nextElementSibling
        if (!panel) continue

        const key = className + "-" + i
        const savedState = sessionStorage.getItem(key)

        // Restore previous accordion state
        if (savedState === "open") {
            item.classList.add("active")
            panel.style.display = displayType
        } else {
            panel.style.display = "none"
        }

        // Toggle accordion on click and save state
        item.addEventListener("click", function () {
            const isOpen = panel.style.display === "none"
            panel.style.display = isOpen ? displayType : "none"
            this.classList.toggle("active", isOpen)
            sessionStorage.setItem(key, isOpen ? "open" : "closed")
        })
    }
}

// Handles the "Add Product" button click to open the product modal and fetch items from the backend.
function setupAddProductBtn() {
    let addProductBtn = document.getElementsByClassName("add-product-btn")
    for (let i = 0; i < addProductBtn.length; i++) {
        addProductBtn[i].addEventListener("click", async function () {
            let category = this.id
            let subcat = this.closest(".subcategory-panel").id
            let color = ""

            // Skip subcategory filter if it matches category
            if (category == subcat) subcat = ""

            let url = `/product/get/?category=${category}&subcategory=${subcat}&color=${color}`

            try {
                const response = await csrfFetch(url)
                if (!response.ok)
                    throw new Error(`Response Status: ${response.status}`)

                const modalHTML = await response.text()
                document.querySelector(".product-modal").innerHTML = modalHTML
                document.querySelector(".product-modal").showModal()

                setupAddProductToEstimateBtn()
            } catch (error) {
                console.error(error.message)
            }
        })
    }

    // Close modal when clicking the close button
    document.addEventListener("click", (e) => {
        if (e.target.closest("#modal-close-btn")) {
            const modal = document.querySelector(".product-modal")
            modal.close()
        }
    })
}

// Handles adding a product from the modal to the current estimate.
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

            // Basic validation for quantity
            if (quantity < 1 || isNaN(quantity)) {
                showInlineError(this, "Quantity must be at least 1.")
                return
            }

            try {
                const response = await csrfFetch(
                    `/estimate/${estimateID}/items/`,
                    {
                        method: "POST",
                        body: JSON.stringify({
                            product_id: productID,
                            quantity: parseInt(quantity),
                        }),
                    },
                )

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

// Handles saving updates to existing estimate line items.
function setupUpdateItemBtn() {
    let saveItemBtn = document.getElementsByClassName("update-item-btn")

    for (let i = 0; i < saveItemBtn.length; i++) {
        saveItemBtn[i].addEventListener("click", async function () {
            const currentRow = this.closest("tr")
            const lineItemID = currentRow.dataset.lineItemId
            const newQuantity = currentRow.querySelector(
                ".estimate-item-quantity",
            ).value

            if (newQuantity < 1 || isNaN(newQuantity)) {
                alert("Quantity must be at least 1")
                return
            }

            try {
                const response = await csrfFetch(
                    `/estimate/items/${lineItemID}`,
                    {
                        method: "PUT",
                        body: JSON.stringify({
                            quantity: parseInt(newQuantity),
                        }),
                    },
                )
                if (!response.ok)
                    throw new Error(`Response status: ${response.status}`)
                location.reload()
            } catch (error) {
                console.error(error.message)
            }
        })
    }
}

// Handles deleting an estimate item when the "Delete" button is clicked.
function setupDeleteBtn() {
    let delItemBtn = document.getElementsByClassName("delete-item-btn")
    for (let i = 0; i < delItemBtn.length; i++) {
        delItemBtn[i].addEventListener("click", async function () {
            const lineItemID = this.closest("tr").dataset.lineItemId
            try {
                const response = await csrfFetch(
                    `/estimate/items/${lineItemID}`,
                    {
                        method: "DELETE",
                    },
                )
                if (!response.ok)
                    throw new Error(`Response status: ${response.status}`)
                location.reload()
            } catch (error) {
                console.error(error.message)
            }
        })
    }
}

// Displays error messages returned by the backend inline in the modal.
function handleProductErrors(button, errors) {
    document.querySelectorAll(".product-error").forEach((e) => e.remove())

    if (errors.quantity) showInlineError(button, errors.quantity)
    if (errors.product) showInlineError(button, errors.product)
}

// Creates a small inline error message below the input.
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

// Initialize ui features like accordians and buttons
setupAccordion("accordion", "block")
setupAccordion("subcategory-accordion", "flex")
setupAddProductBtn()
setupUpdateItemBtn()
setupDeleteBtn()
